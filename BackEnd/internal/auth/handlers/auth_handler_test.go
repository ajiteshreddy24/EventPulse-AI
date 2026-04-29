package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

	authMiddleware "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/middleware"
	authModels "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/models"
	authQueries "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/queries"
	authService "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/service"
)

func setupAuthHandler(t *testing.T) (*AuthHandler, *authMiddleware.AuthMiddleware, *authService.AuthService) {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	_, err = db.Exec(`
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		interests TEXT NOT NULL DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	repo := &authQueries.UserRepository{DB: db}
	svc := &authService.AuthService{Repo: repo}
	return &AuthHandler{Service: svc}, &authMiddleware.AuthMiddleware{Service: svc}, svc
}

func TestRegisterHandlerSuccess(t *testing.T) {
	handler, _, _ := setupAuthHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{
		"name":"Casey",
		"email":"casey@example.com",
		"password":"secret123",
		"interests":["music"]
	}`))
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var response authModels.AuthResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Token == "" || response.User.Email != "casey@example.com" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestLoginHandlerInvalidCredentials(t *testing.T) {
	handler, _, svc := setupAuthHandler(t)

	_, err := svc.Register(authModels.RegisterRequest{
		Name:     "Casey",
		Email:    "casey@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{
		"email":"casey@example.com",
		"password":"wrong-pass"
	}`))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestMeHandlerUnauthorized(t *testing.T) {
	handler, _, _ := setupAuthHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rec := httptest.NewRecorder()

	handler.Me(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestUpdateInterestsHandlerSuccess(t *testing.T) {
	handler, authMW, svc := setupAuthHandler(t)

	registered, err := svc.Register(authModels.RegisterRequest{
		Name:     "Riley",
		Email:    "riley@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	router := mux.NewRouter()
	router.Handle("/auth/me/interests",
		authMW.RequireAuth(http.HandlerFunc(handler.UpdateInterests)),
	).Methods(http.MethodPut)

	req := httptest.NewRequest(http.MethodPut, "/auth/me/interests", bytes.NewBufferString(`{
		"interests":["Art"," art ","Career"]
	}`))
	req.Header.Set("Authorization", "Bearer "+registered.Token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var user authModels.User
	if err := json.NewDecoder(rec.Body).Decode(&user); err != nil {
		t.Fatalf("failed to decode user: %v", err)
	}

	want := []string{"art", "career"}
	if !reflect.DeepEqual(user.Interests, want) {
		t.Fatalf("expected interests %v, got %v", want, user.Interests)
	}
}
