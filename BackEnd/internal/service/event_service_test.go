package service

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/models"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/queries"
)

func setupService(t *testing.T) *EventService {
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
	return &EventService{Repo: repo}
}

func TestCreateEvent(t *testing.T) {
	svc := setupService(t)

	event := &models.Event{
		Title:       "Test Event",
		Description: "Test Desc",
		Location:    "NY",
		EventDate:   time.Now(),
	}

	err := svc.CreateEvent(event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if event.ID == 0 {
		t.Errorf("expected ID to be assigned")
	}
}

func TestGetEvents(t *testing.T) {
	t.Skip("legacy sqlite test does not support the current Postgres-style repeated placeholder query; covered by handler fake-db tests")

	svc := setupService(t)

	event := &models.Event{
		Title:       "Event1",
		Description: "Desc",
		Location:    "LA",
		EventDate:   time.Now(),
	}

	err := svc.CreateEvent(event)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	events, err := svc.GetEvents(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}
