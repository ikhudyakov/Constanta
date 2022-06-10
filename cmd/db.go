package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func dbConnect() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlconn)
	CheckError(err)
	err = db.Ping()
	CheckError(err)
}

func dbInit() {
	go func() {
		dbConnect()
		defer db.Close()
		_, err := db.Exec(`
		DROP TABLE IF EXISTS public.transactions;
		DROP SEQUENCE IF EXISTS public.id_sequence;
		`)
		CheckError(err)

		_, err = db.Exec(`
		SET statement_timeout = 0;
		SET lock_timeout = 0;
		SET idle_in_transaction_session_timeout = 0;
		SET client_encoding = 'UTF8';
		SET standard_conforming_strings = on;
		SELECT pg_catalog.set_config('search_path', '', false);
		SET check_function_bodies = false;
		SET xmloption = content;
		SET client_min_messages = warning;
		SET row_security = off;

		CREATE SEQUENCE public.id_sequence
			START WITH 1
			INCREMENT BY 1
			NO MINVALUE
			NO MAXVALUE
			CACHE 1;

		ALTER TABLE public.id_sequence OWNER TO constanta;
		
		SET default_tablespace = '';
		
		CREATE TABLE public.transactions (
			id bigint DEFAULT nextval('public.id_sequence'::regclass) NOT NULL,
			userid bigint,
			useremail character varying(50),
			amount real,
			currency character(3),
			initdate character varying,
			moddate character varying,
			status character varying(7)
		);
		
		ALTER TABLE public.transactions OWNER TO constanta;
		
		ALTER TABLE ONLY public.transactions
			ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);`)
		CheckError(err)
	}()
}

// Сохранение транзакции в базу данных
func SaveToDB(transaction *Transaction) int64 {
	c1 := make(chan int64)
	go func() {
		dbConnect()
		defer db.Close()

		var lastID int64
		err = db.QueryRow(
			"INSERT INTO transactions (userid, useremail, amount, currency, initdate, moddate, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
			transaction.UserId,
			transaction.UserEmail,
			transaction.Amount,
			transaction.Currency,
			transaction.InitDate,
			transaction.ModDate,
			transaction.Status).Scan(&lastID)
		c1 <- lastID

		CheckError(err)
	}()
	return <-c1
}

// Получение транзакций из базы данных по UserID
func GetFromDBByUserID(userId int64) []Transaction {
	c1 := make(chan []Transaction)

	go func() {
		dbConnect()
		defer db.Close()
		rows, err := db.Query("select * from Transactions where userid = $1", userId)
		CheckError(err)
		defer rows.Close()
		transactions := []Transaction{}

		for rows.Next() {
			t := Transaction{}
			err := rows.Scan(&t.ID, &t.UserId, &t.UserEmail, &t.Amount, &t.Currency, &t.InitDate, &t.ModDate, &t.Status)
			CheckError(err)
			transactions = append(transactions, t)
		}
		c1 <- transactions
	}()
	return <-c1
}

// Получение транзакций из базы данных по UserEmail
func GetFromDBByUserEmail(userEmail string) []Transaction {
	c1 := make(chan []Transaction)

	go func() {
		dbConnect()
		defer db.Close()
		rows, err := db.Query("select * from Transactions where useremail = $1", userEmail)
		CheckError(err)
		defer rows.Close()
		transactions := []Transaction{}

		for rows.Next() {
			t := Transaction{}
			err := rows.Scan(&t.ID, &t.UserId, &t.UserEmail, &t.Amount, &t.Currency, &t.InitDate, &t.ModDate, &t.Status)
			CheckError(err)
			transactions = append(transactions, t)
		}
		c1 <- transactions
	}()
	return <-c1
}

// Получение транзакции из базы данных по ее ID
func GetFromDBByID(id int64) Transaction {
	c1 := make(chan Transaction)

	go func() {
		dbConnect()
		defer db.Close()
		rows, err := db.Query("select * from Transactions where id = $1", id)
		CheckError(err)
		defer rows.Close()
		transactions := Transaction{}

		for rows.Next() {
			t := Transaction{}
			err := rows.Scan(&t.ID, &t.UserId, &t.UserEmail, &t.Amount, &t.Currency, &t.InitDate, &t.ModDate, &t.Status)
			CheckError(err)
			transactions = t
		}
		c1 <- transactions
	}()
	return <-c1
}

// Обновление статуса транзакции в базе данных по ее ID
func UpdateStatusDBByID(id int64, status string) string {
	c1 := make(chan string)

	go func() {
		dbConnect()
		defer db.Close()
		var statusFromBd string
		rows, err := db.Query("select status from Transactions where id = $1", id)
		CheckError(err)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&statusFromBd)
			CheckError(err)
		}

		if statusFromBd != FAIL && statusFromBd != SUCCESS && statusFromBd != "" {
			_, err := db.Exec("update Transactions set status = $1, moddate = $2 where id = $3", status, time.Now().String(), id)
			CheckError(err)
			c1 <- "status updated successfully"
			return
		}
		c1 <- "status update error"
	}()
	return <-c1
}

// Удаление транзакции из базе данных по ее ID
func DeleteFromDBByID(id int64) string {
	c1 := make(chan string)

	go func() {
		dbConnect()
		defer db.Close()
		var statusFromBd string
		rows, err := db.Query("select status from Transactions where id = $1", id)
		CheckError(err)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&statusFromBd)
			CheckError(err)
		}

		if statusFromBd == "" {
			c1 <- "transaction not found"
			return
		}

		if statusFromBd != FAIL && statusFromBd != SUCCESS {
			result, err := db.Exec("delete from Transactions where id = $1", id)
			CheckError(err)
			d, _ := result.RowsAffected()
			if d > 0 {
				c1 <- "transaction deleted successfully"
				return
			}
		}
		c1 <- "transaction not deleted"
	}()
	return <-c1
}

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
