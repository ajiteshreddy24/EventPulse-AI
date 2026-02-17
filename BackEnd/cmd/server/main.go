package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/db"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/handlers"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/queries"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/service"
)

func main() {
	db.Connect()

	repo := &queries.EventRepository{DB: db.DB}
	service := &service.EventService{Repo: repo}
	handler := &handlers.EventHandler{Service: service}

	r := mux.NewRouter()

	r.HandleFunc("/events", handler.CreateEvent).Methods("POST")
	r.HandleFunc("/events", handler.GetEvents).Methods("GET")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
