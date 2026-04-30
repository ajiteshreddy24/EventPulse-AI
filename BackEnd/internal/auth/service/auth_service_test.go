package service

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	authModels "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/models"
	authQueries "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/queries"
)

func setupAuthService(t *testing.T) (*AuthService, *authQueries.UserRepository) {
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
	return &AuthService{Repo: repo}, repo
}

func TestRegisterSuccessNormalizesInput(t *testing.T) {
	svc, repo := setupAuthService(t)

	response, err := svc.Register(authModels.RegisterRequest{
		Name:      "  Alice  ",
		Email:     "  ALICE@EXAMPLE.COM  ",
		Password:  "secret123",
		Interests: []string{"Technology", " technology ", "", "Career"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Token == "" {
		t.Fatal("expected token to be returned")
	}

	if response.User.Email != "alice@example.com" {
		t.Fatalf("expected normalized email, got %q", response.User.Email)
	}

	wantInterests := []string{"technology", "career"}
	if !reflect.DeepEqual(response.User.Interests, wantInterests) {
		t.Fatalf("expected interests %v, got %v", wantInterests, response.User.Interests)
	}

	if response.User.PasswordHash != "" {
		t.Fatalf("expected sanitized user without password hash, got %q", response.User.PasswordHash)
	}

	storedUser, err := repo.GetByEmail("alice@example.com")
	if err != nil {
		t.Fatalf("failed to fetch stored user: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte("secret123")); err != nil {
		t.Fatalf("expected stored password hash to match original password: %v", err)
	}
}

func TestRegisterInvalidPayload(t *testing.T) {
	svc, _ := setupAuthService(t)

	_, err := svc.Register(authModels.RegisterRequest{
		Name:     "Alex",
		Email:    "",
		Password: "secret123",
	})
	if err != ErrInvalidAuthPayload {
		t.Fatalf("expected ErrInvalidAuthPayload, got %v", err)
	}
}

func TestLoginSuccess(t *testing.T) {
	svc, _ := setupAuthService(t)

	registered, err := svc.Register(authModels.RegisterRequest{
		Name:     "Jordan",
		Email:    "jordan@example.com",
		Password: "pass1234",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	response, err := svc.Login(authModels.LoginRequest{
		Email:    "JORDAN@example.com",
		Password: "pass1234",
	})
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}

	if response.Token == "" {
		t.Fatal("expected login token")
	}

	if response.User.ID != registered.User.ID {
		t.Fatalf("expected user ID %d, got %d", registered.User.ID, response.User.ID)
	}
}

func TestUpdateInterestsNormalizesAndPersists(t *testing.T) {
	svc, repo := setupAuthService(t)

	registered, err := svc.Register(authModels.RegisterRequest{
		Name:     "Taylor",
		Email:    "taylor@example.com",
		Password: "pass1234",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	user, err := svc.UpdateInterests(registered.User.ID, []string{" Sports ", "", "sports", "AI"})
	if err != nil {
		t.Fatalf("unexpected error updating interests: %v", err)
	}

	want := []string{"sports", "ai"}
	if !reflect.DeepEqual(user.Interests, want) {
		t.Fatalf("expected interests %v, got %v", want, user.Interests)
	}

	storedUser, err := repo.GetByID(registered.User.ID)
	if err != nil {
		t.Fatalf("failed to fetch stored user: %v", err)
	}

	if !reflect.DeepEqual(storedUser.Interests, want) {
		t.Fatalf("expected persisted interests %v, got %v", want, storedUser.Interests)
	}
}
