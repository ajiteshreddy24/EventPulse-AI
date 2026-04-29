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
		log.Fatal("Error opening DB:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}

	fmt.Println("Connected to PostgreSQL 🚀")

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
		`ALTER TABLE events ADD COLUMN IF NOT EXISTS capacity INTEGER NOT NULL DEFAULT 50`,
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS interests TEXT NOT NULL DEFAULT ''`,
		`CREATE TABLE IF NOT EXISTS rsvps (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    event_id INTEGER REFERENCES events(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, event_id)
)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS waitlist_entries (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (user_id, event_id)
		)`,
	}

	for _, statement := range statements {
		if _, err := DB.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}
