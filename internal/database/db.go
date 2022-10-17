package database

import (
	"database/sql"
	"fmt"
	c "paymentService/internal/config"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func Connect(conf *c.Config) (*sql.DB, error) {
	var err error
	var conn string
	var db *sql.DB

	switch conf.ConnectionType {
	case "postgres":
		conn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Host, conf.Postgresqlport, conf.User, conf.Password, conf.DBname)
		if db, err = sql.Open(conf.ConnectionType, conn); err != nil {
			return nil, err
		}
	case "sqlite3":
		conn = conf.DBname
		if db, err = sql.Open(conf.ConnectionType, conn); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid base type")
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
