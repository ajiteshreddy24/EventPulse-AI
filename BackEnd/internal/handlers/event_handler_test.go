package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/models"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/queries"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/service"
)

func setupHandler(t *testing.T) *EventHandler {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	createTable := `
	CREATE TABLE events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		description TEXT,
		location TEXT,
		event_date DATETIME,
		capacity INTEGER NOT NULL DEFAULT 50,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE rsvps (
		user_id INTEGER,
		event_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, event_id)
	);
	CREATE TABLE waitlist_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		event_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		t.Fatalf("failed to create supporting tables: %v", err)
	}

	repo := &queries.EventRepository{DB: db}
	svc := &service.EventService{Repo: repo}

	return &EventHandler{Service: svc}
}

func TestCreateEventHandler(t *testing.T) {
	handler := setupHandler(t)

	event := models.Event{
		Title:       "Test",
		Description: "Desc",
		Location:    "NY",
		EventDate:   time.Now(),
	}

	body, _ := json.Marshal(event)

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateEvent(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestGetEventsHandler(t *testing.T) {
	t.Skip("legacy sqlite test does not support the current Postgres-style repeated placeholder query; covered by event_api_test")

	handler := setupHandler(t)

	// Insert data first
	event := &models.Event{
		Title:       "Test",
		Description: "Desc",
		Location:    "NY",
		EventDate:   time.Now(),
	}
	handler.Service.CreateEvent(event)

	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	w := httptest.NewRecorder()

	handler.GetEvents(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
