package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func initDBConnection() *sql.DB {
	connStr := os.Getenv("DB_DSN")
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Panic("couldn't connect to database", err)
	}

	return db
}
