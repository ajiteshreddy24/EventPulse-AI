package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	authHandlers "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/handlers"
	authMiddleware "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/middleware"
	authQueries "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/queries"
	authService "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/service"
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
	authRepo := &authQueries.UserRepository{DB: db.DB}
	authSvc := &authService.AuthService{Repo: authRepo}
	authHandler := &authHandlers.AuthHandler{Service: authSvc}
	authMW := &authMiddleware.AuthMiddleware{Service: authSvc}

	r := mux.NewRouter()

	r.HandleFunc("/events", handler.CreateEvent).Methods("POST")
	r.HandleFunc("/events", handler.GetEvents).Methods("GET")
	r.HandleFunc("/events/{id}", handler.UpdateEvent).Methods("PUT")
	r.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	r.Handle("/auth/me", authMW.RequireAuth(http.HandlerFunc(authHandler.Me))).Methods("GET")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
