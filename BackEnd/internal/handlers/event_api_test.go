package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/models"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/queries"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/service"
	"github.com/gorilla/mux"
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

func newTestHandler(t *testing.T, state *fakeDBState) *EventHandler {
	t.Helper()

	driverName := "event_api_fake_driver"
	sql.Register(driverName, &fakeDriver{state: state})

	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("failed to open fake db: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	repo := &queries.EventRepository{DB: db}
	svc := &service.EventService{Repo: repo}
	return &EventHandler{Service: svc}
}

func newTestHandlerWithUniqueDriver(t *testing.T, state *fakeDBState) *EventHandler {
	t.Helper()

	driverName := "event_api_fake_driver_" + strings.ReplaceAll(t.Name(), "/", "_")
	sql.Register(driverName, &fakeDriver{state: state})

	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("failed to open fake db: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	repo := &queries.EventRepository{DB: db}
	svc := &service.EventService{Repo: repo}
	return &EventHandler{Service: svc}
}

func TestCreateEventSuccess(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	now := time.Now()
	state.set("INSERT INTO events", fakeResponse{
		columns: []string{"id", "created_at"},
		rows:    [][]driver.Value{{int64(1), now}},
	})

	handler := newTestHandlerWithUniqueDriver(t, state)
	body := `{"title":"Demo","description":"Launch","location":"NYC","event_date":"2026-03-25T12:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.CreateEvent(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var event models.Event
	if err := json.NewDecoder(rec.Body).Decode(&event); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if event.ID != 1 {
		t.Fatalf("expected event ID 1, got %d", event.ID)
	}
}

func TestCreateEventDatabaseError(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	state.set("INSERT INTO events", fakeResponse{err: errors.New("insert failed")})

	handler := newTestHandlerWithUniqueDriver(t, state)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(`{"title":"Demo"}`))
	rec := httptest.NewRecorder()

	handler.CreateEvent(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestGetEventsSuccess(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	now := time.Now()
	state.set("FROM events", fakeResponse{
		columns: []string{"id", "title", "description", "location", "event_date", "created_at"},
		rows: [][]driver.Value{
			{int64(1), "Demo", "Launch", "NYC", now, now},
			{int64(2), "Meetup", "Community", "Boston", now, now},
		},
	})

	handler := newTestHandlerWithUniqueDriver(t, state)
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rec := httptest.NewRecorder()

	handler.GetEvents(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var events []models.Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("failed to decode events: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestGetEventsDatabaseError(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	state.set("FROM events", fakeResponse{
		err: errors.New("select failed"),
	})

	handler := newTestHandlerWithUniqueDriver(t, state)
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rec := httptest.NewRecorder()

	handler.GetEvents(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestUpdateEventSuccess(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	now := time.Now()
	state.set("UPDATE events", fakeResponse{
		columns: []string{"created_at"},
		rows:    [][]driver.Value{{now}},
	})

	handler := newTestHandlerWithUniqueDriver(t, state)
	body := `{"title":"Updated","description":"Updated desc","location":"NYC","event_date":"2026-03-26T12:00:00Z"}`
	req := httptest.NewRequest(http.MethodPut, "/events/5", bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	rec := httptest.NewRecorder()

	handler.UpdateEvent(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestUpdateEventInvalidID(t *testing.T) {
	handler := &EventHandler{}
	req := httptest.NewRequest(http.MethodPut, "/events/abc", bytes.NewBufferString(`{}`))
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rec := httptest.NewRecorder()

	handler.UpdateEvent(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestUpdateEventInvalidBody(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	handler := newTestHandlerWithUniqueDriver(t, state)
	req := httptest.NewRequest(http.MethodPut, "/events/5", bytes.NewBufferString(`{"title":`))
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	rec := httptest.NewRecorder()

	handler.UpdateEvent(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestUpdateEventNotFound(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	state.set("UPDATE events", fakeResponse{err: sql.ErrNoRows})

	handler := newTestHandlerWithUniqueDriver(t, state)
	req := httptest.NewRequest(http.MethodPut, "/events/5", bytes.NewBufferString(`{"title":"Updated"}`))
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	rec := httptest.NewRecorder()

	handler.UpdateEvent(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 with current implementation, got %d", rec.Code)
	}
}

func TestDeleteEventInvalidID(t *testing.T) {
	handler := &EventHandler{}
	req := httptest.NewRequest(http.MethodDelete, "/events/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rec := httptest.NewRecorder()

	handler.DeleteEvent(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestDeleteEventSuccess(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	state.set("DELETE FROM events", fakeResponse{})

	handler := newTestHandlerWithUniqueDriver(t, state)
	req := httptest.NewRequest(http.MethodDelete, "/events/5", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	rec := httptest.NewRecorder()

	handler.DeleteEvent(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}
