# Event Handler Test Documentation

## Overview
This document explains the test cases written for the `EventHandler`. These tests validate API endpoints using an in-memory SQLite database and Go’s `httptest` package.

---

## Test Setup

### In-Memory Database
An in-memory SQLite database (`:memory:`) is used to simulate real database interactions without affecting production data.

### Tables Created
- events
- rsvps
- waitlist_entries

### Setup Function

```go
func setupHandler(t *testing.T) *EventHandler
```

This function:
- Initializes an in-memory database
- Creates required tables
- Sets up repository and service layers
- Returns an EventHandler instance

---

## Test Cases

### 1. Create Event (POST /events)

#### Function
```go
func TestCreateEventHandler(t *testing.T)
```

#### Description
Tests successful creation of an event.

#### Steps
1. Initialize handler
2. Create event payload
3. Send POST request
4. Validate response

#### Expected Result
- HTTP Status: 201 Created

---

### 2. Get Events (GET /events)

#### Function
```go
func TestGetEventsHandler(t *testing.T)
```

#### Description
Tests retrieval of all events.

#### Current Status
Skipped due to SQLite not supporting Postgres-style query placeholders.

```go
t.Skip("legacy sqlite test does not support the current Postgres-style repeated placeholder query")
```

#### Expected Result
- HTTP Status: 200 OK

---

## Tools Used

- net/http/httptest for request simulation
- encoding/json for request payloads
- in-memory SQLite database for isolation

---

## Limitations

- SQLite differs from Postgres in query syntax
- Some tests may not reflect production DB behavior

---

## Recommended Improvements

- Use mocks instead of real DB for unit tests
- Add table-driven tests
- Add edge case coverage (invalid JSON, missing fields)
- Add integration tests with Postgres

---

## Summary

These tests validate:
- API endpoint functionality
- Correct HTTP responses
- Database interaction flow
