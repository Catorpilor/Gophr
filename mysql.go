package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var globalMySQLDB *sql.DB

/*
func init() {
	db, err := NewMySQLDB("zhimaa:zhimaa@tcp(127.0.0.1)/gopher?parseTime=true")
	if err != nil {
		panic(err)
	}
	globalMySQLDB = db
}
*/

func NewMySQLDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
