package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/db"
)

func main() {
	db.Connect()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "EventPulse-AI running 🚀")
	})

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
