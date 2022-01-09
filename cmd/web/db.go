/**
package db

Main package for opening up a database with mysql driver using dsn
**/
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func OpenDB(table, user, pass string, params ...string) (*sql.DB, error) {
	p := ""
	for _, param := range params {
		p += fmt.Sprintf("?%v", param)
	}
	dsn := fmt.Sprintf("%s:%s@/%s%s", user, pass, table, p)
	db, err := sql.Open("mysql", dsn)
	fmt.Println("p is", dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}
