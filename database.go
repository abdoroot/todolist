package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "127.0.0.1"
	port     = 3306
	user     = "root"
	password = ""
	dbname   = "todo"
)

func DbConnect() (*sql.DB, error) {
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
