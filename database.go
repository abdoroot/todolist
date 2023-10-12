package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "todo"
)

func PgDbConnect() (*sql.DB, error) {
	var db *sql.DB
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return db, err
	}

	if err = db.Ping(); err != nil {
		return db, err
	}

	return db, err
}

func MySqlDbConnect() (*sql.DB, error) {
	var db *sql.DB
	mysqllInfo := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, password, host, port, dbname)
	db, err := sql.Open("mysql", mysqllInfo)
	if err != nil {
		return db, err
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		return db, err
	}

	return db, nil
}
