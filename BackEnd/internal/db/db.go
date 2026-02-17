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
}
