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

	// CORS
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	api := r.PathPrefix("/api").Subrouter()

	// EVENTS (Protected where needed)
	api.HandleFunc("/events", handler.GetEvents).Methods("GET")
	api.HandleFunc("/events/{id}", handler.GetEventByID).Methods("GET")

	api.Handle("/events",
		authMW.RequireAuth(http.HandlerFunc(handler.CreateEvent)),
	).Methods("POST")

	api.Handle("/events/{id}",
		authMW.RequireAuth(http.HandlerFunc(handler.UpdateEvent)),
	).Methods("PUT")

	api.Handle("/events/{id}",
		authMW.RequireAuth(http.HandlerFunc(handler.DeleteEvent)),
	).Methods("DELETE")

	api.Handle("/events/{id}/rsvp",
		authMW.RequireAuth(http.HandlerFunc(handler.RSVP)),
	).Methods("POST")

	api.Handle("/events/{id}/rsvp",
		authMW.RequireAuth(http.HandlerFunc(handler.CancelRSVP)),
	).Methods("DELETE")

	// AUTH
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.Handle("/auth/me",
		authMW.RequireAuth(http.HandlerFunc(authHandler.Me)),
	).Methods("GET")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}