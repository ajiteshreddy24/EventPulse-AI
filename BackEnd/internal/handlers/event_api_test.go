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
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	authMiddleware "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/middleware"
	authService "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/service"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/models"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/queries"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/service"
	"github.com/golang-jwt/jwt/v5"
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

func serveAuthenticatedRequest(t *testing.T, handler http.HandlerFunc, method, routePattern, requestPath string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()

	authMW := &authMiddleware.AuthMiddleware{Service: &authService.AuthService{}}
	router := mux.NewRouter()
	router.Handle(routePattern, authMW.RequireAuth(handler)).Methods(method)

	req := httptest.NewRequest(method, requestPath, body)
	req.Header.Set("Authorization", "Bearer "+testBearerToken(t, 42))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
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
		columns: []string{"id", "title", "description", "location", "event_date", "capacity", "created_at", "rsvp_count", "user_has_rsvp", "waitlist_count", "user_on_waitlist"},
		rows: [][]driver.Value{
			{int64(1), "Demo", "Launch", "NYC", now, int64(50), now, int64(3), false, int64(1), false},
			{int64(2), "Meetup", "Community", "Boston", now, int64(100), now, int64(8), true, int64(0), false},
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

func TestGetEventByIDSuccess(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	now := time.Now()
	state.set("WHERE e.id = $1", fakeResponse{
		columns: []string{"id", "title", "description", "location", "event_date", "capacity", "created_at", "rsvp_count", "user_has_rsvp", "waitlist_count", "user_on_waitlist"},
		rows: [][]driver.Value{
			{int64(7), "Demo Event", "Launch night", "NYC", now, int64(50), now, int64(12), true, int64(2), false},
		},
	})

	handler := newTestHandlerWithUniqueDriver(t, state)
	handler.AuthService = &authService.AuthService{}

	req := httptest.NewRequest(http.MethodGet, "/events/7", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "7"})
	req.Header.Set("Authorization", "Bearer "+testBearerToken(t, 99))
	rec := httptest.NewRecorder()

	handler.GetEventByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var event models.Event
	if err := json.NewDecoder(rec.Body).Decode(&event); err != nil {
		t.Fatalf("failed to decode event: %v", err)
	}

	if event.ID != 7 || event.Title != "Demo Event" {
		t.Fatalf("unexpected event payload: %+v", event)
	}
}

func TestGetEventByIDInvalidID(t *testing.T) {
	handler := &EventHandler{}
	req := httptest.NewRequest(http.MethodGet, "/events/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rec := httptest.NewRecorder()

	handler.GetEventByID(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestGetEventByIDNotFound(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	state.set("WHERE e.id = $1", fakeResponse{err: sql.ErrNoRows})

	handler := newTestHandlerWithUniqueDriver(t, state)
	req := httptest.NewRequest(http.MethodGet, "/events/17", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "17"})
	rec := httptest.NewRecorder()

	handler.GetEventByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
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

func TestRSVPAlreadyExists(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	now := time.Now()
	state.set("WHERE e.id = $1", fakeResponse{
		columns: []string{"id", "title", "description", "location", "event_date", "capacity", "created_at", "rsvp_count", "user_has_rsvp", "waitlist_count", "user_on_waitlist"},
		rows: [][]driver.Value{
			{int64(3), "Hack Night", "Build night", "Lab", now, int64(20), now, int64(4), true, int64(0), false},
		},
	})

	handler := newTestHandlerWithUniqueDriver(t, state)
	rec := serveAuthenticatedRequest(t, http.HandlerFunc(handler.RSVP), http.MethodPost, "/events/{id}/rsvp", "/events/3/rsvp", nil)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", rec.Code)
	}
}

func TestJoinWaitlistSuccess(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	now := time.Now()
	state.set("WHERE e.id = $1", fakeResponse{
		columns: []string{"id", "title", "description", "location", "event_date", "capacity", "created_at", "rsvp_count", "user_has_rsvp", "waitlist_count", "user_on_waitlist"},
		rows: [][]driver.Value{
			{int64(5), "Career Expo", "Meet recruiters", "Hall A", now, int64(10), now, int64(10), false, int64(3), false},
		},
	})
	state.set("INSERT INTO waitlist_entries", fakeResponse{})

	handler := newTestHandlerWithUniqueDriver(t, state)
	rec := serveAuthenticatedRequest(t, http.HandlerFunc(handler.JoinWaitlist), http.MethodPost, "/events/{id}/waitlist", "/events/5/waitlist", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var payload map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload["message"] != "Added to waitlist" {
		t.Fatalf("unexpected response payload: %+v", payload)
	}
}

func TestCreateCommentValidationError(t *testing.T) {
	handler := &EventHandler{Service: &service.EventService{}}
	rec := serveAuthenticatedRequest(
		t,
		http.HandlerFunc(handler.CreateComment),
		http.MethodPost,
		"/events/{id}/comments",
		"/events/8/comments",
		bytes.NewBufferString(`{"comment":"   "}`),
	)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCreateCommentSuccess(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	now := time.Now()
	state.set("WHERE e.id = $1", fakeResponse{
		columns: []string{"id", "title", "description", "location", "event_date", "capacity", "created_at", "rsvp_count", "user_has_rsvp", "waitlist_count", "user_on_waitlist"},
		rows: [][]driver.Value{
			{int64(8), "Music Fest", "Outdoor concert", "Amphitheater", now, int64(80), now, int64(15), false, int64(0), false},
		},
	})
	state.set("INSERT INTO comments", fakeResponse{
		columns: []string{"id", "event_id", "user_id", "content", "created_at"},
		rows:    [][]driver.Value{{int64(11), int64(8), int64(42), "Can't wait!", now}},
	})
	state.set("FROM users", fakeResponse{
		columns: []string{"id", "name", "email"},
		rows:    [][]driver.Value{{int64(42), "Taylor", "taylor@example.com"}},
	})

	handler := newTestHandlerWithUniqueDriver(t, state)
	rec := serveAuthenticatedRequest(
		t,
		http.HandlerFunc(handler.CreateComment),
		http.MethodPost,
		"/events/{id}/comments",
		"/events/8/comments",
		bytes.NewBufferString(`{"comment":"Can't wait!"}`),
	)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var comment models.Comment
	if err := json.NewDecoder(rec.Body).Decode(&comment); err != nil {
		t.Fatalf("failed to decode comment: %v", err)
	}

	if comment.ID != 11 || comment.User.Name != "Taylor" {
		t.Fatalf("unexpected comment payload: %+v", comment)
	}
}

func TestGetCommentsSuccess(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	now := time.Now()
	state.set("FROM comments c", fakeResponse{
		columns: []string{"id", "event_id", "user_id", "content", "created_at", "user_id_join", "name", "email"},
		rows: [][]driver.Value{
			{int64(1), int64(9), int64(4), "First!", now, int64(4), "Sam", "sam@example.com"},
			{int64(2), int64(9), int64(5), "Looking forward to it", now, int64(5), "Jordan", "jordan@example.com"},
		},
	})

	handler := newTestHandlerWithUniqueDriver(t, state)
	req := httptest.NewRequest(http.MethodGet, "/events/9/comments", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "9"})
	rec := httptest.NewRecorder()

	handler.GetComments(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var comments []models.Comment
	if err := json.NewDecoder(rec.Body).Decode(&comments); err != nil {
		t.Fatalf("failed to decode comments: %v", err)
	}

	if len(comments) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(comments))
	}
}

func TestDeleteCommentNotFound(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	state.set("DELETE FROM comments", fakeResponse{err: queries.ErrEventNotFound})

	handler := newTestHandlerWithUniqueDriver(t, state)
	req := httptest.NewRequest(http.MethodDelete, "/comments/7", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "7"})
	rec := httptest.NewRecorder()

	handler.DeleteComment(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetRecommendationsUnauthorized(t *testing.T) {
	handler := &EventHandler{Service: &service.EventService{}}
	req := httptest.NewRequest(http.MethodGet, "/events/recommendations", nil)
	rec := httptest.NewRecorder()

	handler.GetRecommendations(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestGetMyRSVPedEventsSuccess(t *testing.T) {
	state := &fakeDBState{responses: map[string]fakeResponse{}}
	now := time.Now()
	state.set("JOIN rsvps rsvp", fakeResponse{
		columns: []string{"id", "title", "description", "location", "event_date", "capacity", "created_at", "rsvp_count"},
		rows: [][]driver.Value{
			{int64(4), "AI Meetup", "Campus networking", "Innovation Hub", now, int64(80), now, int64(22)},
			{int64(9), "Career Fair", "Internship recruiting", "Student Center", now.Add(2 * time.Hour), int64(150), now, int64(89)},
		},
	})

	handler := newTestHandlerWithUniqueDriver(t, state)
	authMW := &authMiddleware.AuthMiddleware{Service: &authService.AuthService{}}
	router := mux.NewRouter()
	router.Handle("/events/attending",
		authMW.RequireAuth(http.HandlerFunc(handler.GetMyRSVPedEvents)),
	).Methods(http.MethodGet)

	req := httptest.NewRequest(http.MethodGet, "/events/attending", nil)
	req.Header.Set("Authorization", "Bearer "+testBearerToken(t, 42))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var events []models.Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("failed to decode attending events: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 attending events, got %d", len(events))
	}

	if events[0].ID != 4 || events[1].ID != 9 {
		t.Fatalf("unexpected attending events returned: %+v", events)
	}
}

func TestGetMyRSVPedEventsUnauthorized(t *testing.T) {
	handler := &EventHandler{}
	authMW := &authMiddleware.AuthMiddleware{Service: &authService.AuthService{}}
	router := mux.NewRouter()
	router.Handle("/events/attending",
		authMW.RequireAuth(http.HandlerFunc(handler.GetMyRSVPedEvents)),
	).Methods(http.MethodGet)

	req := httptest.NewRequest(http.MethodGet, "/events/attending", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func testBearerToken(t *testing.T, userID int) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(userID),
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	signedToken, err := token.SignedString([]byte("eventpulse-dev-secret"))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	return signedToken
}
