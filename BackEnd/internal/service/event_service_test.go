package service

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	authQueries "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/queries"
	_ "github.com/mattn/go-sqlite3"

	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/models"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/queries"
)

type fakeResponse struct {
	columns []string
	rows    [][]driver.Value
	err     error
}

type fakeDBState struct {
	mu        sync.Mutex
	responses map[string]fakeResponse
}

func (s *fakeDBState) set(key string, response fakeResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.responses[key] = response
}

func (s *fakeDBState) get(query string) fakeResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, response := range s.responses {
		if strings.Contains(query, key) {
			return response
		}
	}

	return fakeResponse{err: errors.New("unexpected query: " + query)}
}

type fakeDriver struct {
	state *fakeDBState
}

func (d *fakeDriver) Open(name string) (driver.Conn, error) {
	return &fakeConn{state: d.state}, nil
}

type fakeConn struct {
	state *fakeDBState
}

func (c *fakeConn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("prepare not supported")
}

func (c *fakeConn) Close() error {
	return nil
}

func (c *fakeConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions not supported")
}

func (c *fakeConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	response := c.state.get(query)
	if response.err != nil {
		return nil, response.err
	}

	return &fakeRows{
		columns: response.columns,
		rows:    response.rows,
	}, nil
}

func (c *fakeConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	response := c.state.get(query)
	if response.err != nil {
		return nil, response.err
	}

	return driver.RowsAffected(1), nil
}

type fakeRows struct {
	columns []string
	rows    [][]driver.Value
	index   int
}

func (r *fakeRows) Columns() []string {
	return r.columns
}

func (r *fakeRows) Close() error {
	return nil
}

func (r *fakeRows) Next(dest []driver.Value) error {
	if r.index >= len(r.rows) {
		return io.EOF
	}

	for i, value := range r.rows[r.index] {
		dest[i] = value
	}

	r.index++
	return nil
}

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

func setupFakeService(t *testing.T, state *fakeDBState) *EventService {
	t.Helper()

	driverName := "event_service_fake_driver_" + strings.ReplaceAll(t.Name(), "/", "_")
	sql.Register(driverName, &fakeDriver{state: state})

	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("failed to open fake db: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	repo := &queries.EventRepository{DB: db}
	userRepo := &authQueries.UserRepository{DB: db}
	return &EventService{Repo: repo, UserRepo: userRepo}
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

func TestGetRSVPedEvents(t *testing.T) {
	svc := setupService(t)

	event := &models.Event{
		Title:       "AI Mixer",
		Description: "Meet builders on campus",
		Location:    "Innovation Lab",
		EventDate:   time.Now().Add(24 * time.Hour),
	}

	if err := svc.CreateEvent(event); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if _, err := svc.Repo.DB.Exec(`
		INSERT INTO rsvps (user_id, event_id)
		VALUES ($1, $2)
	`, 7, event.ID); err != nil {
		t.Fatalf("failed to insert rsvp: %v", err)
	}

	events, err := svc.GetRSVPedEvents(7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 attending event, got %d", len(events))
	}

	if events[0].ID != event.ID {
		t.Fatalf("expected event ID %d, got %d", event.ID, events[0].ID)
	}
}

func TestCreateCommentRejectsBlankComment(t *testing.T) {
	svc := &EventService{}

	comment, err := svc.CreateComment(1, 2, "   ")
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if comment != nil {
		t.Fatalf("expected nil comment, got %+v", comment)
	}
}

func TestCreateCommentSuccess(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	svc := setupFakeService(t, state)
	now := time.Now()

	state.set("WHERE e.id = $1", fakeResponse{
		columns: []string{
			"id", "title", "description", "location", "event_date", "capacity", "created_at",
			"rsvp_count", "user_has_rsvp", "waitlist_count", "user_on_waitlist",
		},
		rows: [][]driver.Value{{int64(2), "Tech Talk", "Go tips", "Room 101", now, int64(50), now, int64(5), false, int64(0), false}},
	})
	state.set("INSERT INTO comments", fakeResponse{
		columns: []string{"id", "event_id", "user_id", "content", "created_at"},
		rows:    [][]driver.Value{{int64(11), int64(2), int64(7), "Great session", now}},
	})
	state.set("FROM users", fakeResponse{
		columns: []string{"id", "name", "email"},
		rows:    [][]driver.Value{{int64(7), "Alex", "alex@example.com"}},
	})

	comment, err := svc.CreateComment(7, 2, "  Great session  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if comment.Content != "Great session" {
		t.Fatalf("expected trimmed comment content, got %q", comment.Content)
	}

	if comment.User.Name != "Alex" {
		t.Fatalf("expected comment author Alex, got %+v", comment.User)
	}
}

func TestJoinWaitlistSuccess(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	svc := setupFakeService(t, state)
	now := time.Now()

	state.set("WHERE e.id = $1", fakeResponse{
		columns: []string{
			"id", "title", "description", "location", "event_date", "capacity", "created_at",
			"rsvp_count", "user_has_rsvp", "waitlist_count", "user_on_waitlist",
		},
		rows: [][]driver.Value{{int64(9), "Career Fair", "Meet employers", "Ballroom", now, int64(10), now, int64(10), false, int64(2), false}},
	})
	state.set("INSERT INTO waitlist_entries", fakeResponse{})

	if err := svc.JoinWaitlist(4, 9); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJoinWaitlistOpenSeatsReturnsEventFull(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	svc := setupFakeService(t, state)
	now := time.Now()

	state.set("WHERE e.id = $1", fakeResponse{
		columns: []string{
			"id", "title", "description", "location", "event_date", "capacity", "created_at",
			"rsvp_count", "user_has_rsvp", "waitlist_count", "user_on_waitlist",
		},
		rows: [][]driver.Value{{int64(9), "Career Fair", "Meet employers", "Ballroom", now, int64(10), now, int64(9), false, int64(0), false}},
	})

	err := svc.JoinWaitlist(4, 9)
	if err != queries.ErrEventFull {
		t.Fatalf("expected ErrEventFull, got %v", err)
	}
}

func TestGetRecommendationsRanksAndFilters(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	svc := setupFakeService(t, state)
	now := time.Now()

	state.set("WHERE id = $1", fakeResponse{
		columns: []string{"id", "name", "email", "password_hash", "interests", "created_at"},
		rows:    [][]driver.Value{{int64(15), "Morgan", "morgan@example.com", "hash", "technology,career", now}},
	})
	state.set("ORDER BY event_date ASC", fakeResponse{
		columns: []string{
			"id", "title", "description", "location", "event_date", "capacity", "created_at",
			"rsvp_count", "user_has_rsvp", "waitlist_count", "user_on_waitlist",
		},
		rows: [][]driver.Value{
			{int64(2), "Technology Meetup", "Software builders welcome", "Innovation Lab", now.Add(24 * time.Hour), int64(0), now, int64(20), false, int64(0), false},
			{int64(1), "Coding Career Workshop", "Resume and interview prep", "Career Center", now.Add(48 * time.Hour), int64(30), now, int64(5), false, int64(0), false},
			{int64(3), "Gaming Night", "Board games and pizza", "Student Union", now.Add(72 * time.Hour), int64(40), now, int64(12), false, int64(0), false},
			{int64(4), "Career Expo", "Professional networking", "Arena", now.Add(24 * time.Hour), int64(100), now, int64(8), true, int64(0), false},
		},
	})

	events, err := svc.GetRecommendations(15)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 recommendations, got %d", len(events))
	}

	if events[0].ID != 2 || events[1].ID != 1 {
		t.Fatalf("unexpected recommendation order: %+v", events)
	}

	if events[0].RecommendationReason == "" || events[1].RecommendationReason == "" {
		t.Fatalf("expected recommendation reasons, got %+v", events)
	}
}

func TestGetRecommendationsWithNoInterestsReturnsEmpty(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	svc := setupFakeService(t, state)
	now := time.Now()

	state.set("WHERE id = $1", fakeResponse{
		columns: []string{"id", "name", "email", "password_hash", "interests", "created_at"},
		rows:    [][]driver.Value{{int64(15), "Morgan", "morgan@example.com", "hash", "", now}},
	})

	events, err := svc.GetRecommendations(15)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 0 {
		t.Fatalf("expected no recommendations, got %+v", events)
	}
}

func TestInterestKeywordsFallback(t *testing.T) {
	got := interestKeywords("climate")
	want := []string{"climate"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
