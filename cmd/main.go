package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	c "paymentService/internal/config"
	"paymentService/internal/database"
	"paymentService/internal/repo"
	pg "paymentService/internal/repo/postgresql"
	"paymentService/internal/server"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	var dbmanager repo.DBmanager
	var db *sql.DB
	var err error
	var conf *c.Config

	conf, err = c.GetConfig()
	if err != nil {
		log.Println(err)
		return
	}

	if err = migration(conf); err != nil {
		log.Println(err)
	}

	if db, err = database.Connect(conf); err != nil {
		log.Println((err.Error()))
		return
	}

	dbmanager = &pg.PGmanager{
		DB: db,
	}

	r := server.HandlersInit(dbmanager, conf)
	server := &http.Server{
		Addr:    conf.Port,
		Handler: r,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}

func migration(conf *c.Config) error {
	m, err := migrate.New(
		"file://./migrations",
		fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", conf.ConnectionType, conf.User, conf.Password, conf.Host, conf.Postgresqlport, conf.DBname))
	if err != nil {
		return err
	}
	// if err := m.Down(); err != nil {
	// 	return err
	// }
	if err := m.Up(); err != nil {
		return err
	}
	return err
}
