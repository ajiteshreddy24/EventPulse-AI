# EventPulseAI – Sprint 2 Completion Report

## Sprint 2 Overview

Sprint 2 focused on:

- Full CRUD functionality
- Authentication implementation
- Frontend–Backend integration
- Backend layered unit testing
- End-to-End and Unit testing for frontend

---

#  Sprint 2 Demo Videos

🎬 **Frontend Demo:**  
https://youtu.be/kHt29qN58cE  

🎬 **Backend Demo (Updated):**  
https://youtu.be/Slh6_E3W0nI  

---

# System Architecture

## Frontend
- React (Vite)
- React Router
- Cypress (E2E Testing)
- Vitest (Unit Testing)

## Backend
- Go (Golang)
- Gorilla Mux Router
- PostgreSQL (Production)
- SQLite (:memory:) for testing
- JWT Authentication
- Layered Architecture:
  - Handler Layer
  - Service Layer
  - Repository Layer

---

# 💻 Frontend – Sprint 2

## Completed Work

- Integrated frontend with backend APIs
- Implemented Event Update Feature
- Implemented Event Delete Feature
- Wrote Cypress End-to-End Tests
- Added Unit Tests using Vitest

---

## 🔹 CRUD API Integration

| Method | Endpoint | Purpose |
|--------|----------|----------|
| GET | /events | Retrieve all events |
| POST | /events | Create event |
| PUT | /events/{id} | Update event |
| DELETE | /events/{id} | Delete event |

All operations update the database and reflect immediately in the UI.

---

# Frontend Testing

## Cypress End-to-End Testing

### Test Files:

- `basic.cy.js`
- `createEvent.cy.js`
- `updateEvent.cy.js`
- `deleteEvent.cy.js`
- `navigation.cy.js`

### What They Validate:

- Event creation workflow
- Event update workflow
- Event deletion from UI
- Navigation between pages
- Complete CRUD flow integration

### How to Run Cypress:

```bash
cd FrontEnd
npm install cypress --save-dev
npx cypress open
```

Select and run any test file from the Cypress UI.

---

# Backend – Sprint 2

## Implemented Features

- DELETE `/events/{id}` endpoint
- Proper error handling
- Authentication system
- Layered unit testing

---

# Authentication Implementation

### Endpoints:

| Method | Endpoint | Purpose |
|--------|----------|----------|
| POST | /auth/register | Register user |
| POST | /auth/login | Login user |
| GET | /auth/me | Get authenticated user details |

### Security Features:

- Password hashing
- JWT token generation
- Middleware token validation
- Protected routes

---

# Backend Unit Testing

Testing is divided into:

- Service Layer Tests
- Handler Layer Tests

All tests use:

- Go `httptest`
- In-memory SQLite database (`:memory:`)
- Isolated test environment

---

##  Backend Unit Test Cases List

---

## 🔹 Service Layer

### TestCreateEvent

Verifies that a new event is successfully created.

Checks:
- No error is returned
- Event ID is generated
- Database insert succeeds

---

### TestGetEvents

Verifies retrieval of events from the database.

Checks:
- No error is returned
- Correct number of events is fetched (after inserting one)

---

## 🔹 Handler Layer

### TestCreateEventHandler

Verifies HTTP POST `/events`.

Checks:
- Request is properly handled
- Response status is **201 Created**

---

### TestGetEventsHandler

Verifies HTTP GET `/events`.

Checks:
- Request is properly handled
- Response status is **200 OK**
- Data returned successfully

---

### TestUpdateEventHandler

Verifies HTTP PUT `/events/{id}`.

Checks:
- Response status is **200 OK**
- Event data is updated correctly

---

### TestDeleteEventHandler

Verifies HTTP DELETE `/events/{id}`.

Checks:
- Response status is **200 OK**
- Event is removed from database
- Event cannot be retrieved afterward

---

# Full System Flow

1. User interacts with React UI
2. API call sent to Go backend
3. Handler validates HTTP request
4. Service layer applies business logic
5. Repository interacts with database
6. Response returned to frontend
7. UI updates dynamically

---

# Sprint 2 Achievements

✔ Full CRUD implementation  
✔ JWT Authentication  
✔ Layered backend testing  
✔ Frontend–Backend integration  
✔ Cypress E2E coverage  
✔ Vitest component coverage  
✔ Clean modular architecture  

---

# Future Improvements

- Role-based access control
- Event filtering & search
- Pagination
- AI-powered recommendations
- Deployment automation
- CI/CD integration

---

# Academic Summary

GatorHive Sprint 2 demonstrates:

- Full-stack integration
- Secure authentication
- REST API design
- Test-driven backend architecture
- End-to-end frontend validation
- Clean separation of concerns

The system is now scalable, secure, and fully tested across both frontend and backend layers.