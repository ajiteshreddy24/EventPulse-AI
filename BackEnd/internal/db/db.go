package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := "host=localhost port=5432 user=eventuser password=password dbname=eventpulse sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening Database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error connecting to Database:", err)
	}

	fmt.Println("Connected to PostgreSQL ")

	if err := migrate(); err != nil {
		log.Fatal("Error migrating DB:", err)
	}
}

func migrate() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			event_date TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, statement := range statements {
		if _, err := DB.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}
