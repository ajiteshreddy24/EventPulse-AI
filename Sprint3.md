# GatorHive - Sprint 3 Completion Report

## Overview

Sprint 3 focused on completing the authenticated event-management flow for GatorHive. The project now supports user login and signup, protected event creation/edit/delete actions, RSVP and cancel-RSVP actions, a monthly calendar view, and a refreshed University of Florida-inspired interface with a gator-themed home page.

The application is split into:

- Frontend: React, Vite, React Router, FullCalendar, Cypress.
- Backend: Go, Gorilla Mux, PostgreSQL, JWT authentication, repository/service/handler layers.

#  Sprint 3 Demo Videos

🎬 **Frontend Demo:**  
[https://youtu.be/kHt29qN58cE](https://www.youtube.com/watch?v=z2KrQTPA_ZM)  

🎬 **Backend Demo (Updated):**  
[https://youtu.be/Slh6_E3W0nI](https://www.youtube.com/watch?v=aBJRcEeNJoI)  

## Frontend Sprint 3

## Features Implemented

### 1. Authentication UI and Session Handling

- Added login and signup pages for user authentication.
- Added an `AuthProvider` context to manage the logged-in user, loading state, login, and logout.
- Stores JWT tokens in `localStorage` after successful login/signup.
- Sends the JWT token in protected API requests with the `Authorization: Bearer <token>` header.
- Protects create and edit event routes using `ProtectedRoute`.
- Shows authenticated navigation actions only when a user is logged in.

### 2. Event Management UI

- Users can view all events on the events page.
- Authenticated users can create events through the protected `/create` route.
- Authenticated users can edit existing events through `/edit/:id`.
- Authenticated users can delete events from the events page.
- API failures are surfaced in the UI through error messages instead of failing silently.
- Event dates are converted to ISO format before being sent to the backend.

### 3. RSVP and Cancel RSVP

- Added RSVP and cancel-RSVP buttons on event cards.
- Logged-out users are redirected to login when attempting protected actions.
- RSVP requests call `POST /api/events/{id}/rsvp`.
- Cancel RSVP requests call `DELETE /api/events/{id}/rsvp`.
- The frontend reloads events after RSVP/delete actions so the displayed state stays synchronized with the backend.

### 4. Full Calendar View

- Integrated `@fullcalendar/react` and `@fullcalendar/daygrid`.
- Added a `/calendar` route.
- Events are fetched from the backend and rendered in a monthly calendar grid.
- The calendar gives users another way to scan upcoming events by date.

### 5. UI Refresh

- Reworked the home page into a cleaner landing page with hero copy, action buttons, and event preview cards.
- Added a gator-themed wallpaper using blue, orange, white, and green visual elements inspired by the provided reference.
- Removed emoji text from the frontend for a cleaner, more professional UI.
- Improved card layering so event buttons are clickable and not blocked by decorative overlays.
- Updated navbar, buttons, cards, form fields, and responsive layouts.

## Frontend Testing and Validation

The frontend currently includes Cypress E2E test files for:

- Home page load.
- Navigation between pages.
- Event creation flow.
- Event update flow.
- Event delete flow.

Validation commands used:

```bash
cd FrontEnd
npm run lint
npm run build
```

Note: the current frontend `package.json` does not define an `npm test` script. Cypress tests can be run with Cypress tooling after the frontend/backend stack is running.

## Backend Sprint 3

## Features Implemented

### 1. Authentication APIs

- Added user registration through:

```text
POST /api/auth/register
```

- Added user login through:

```text
POST /api/auth/login
```

- Added current-user lookup through:

```text
GET /api/auth/me
```

- Passwords are hashed using bcrypt.
- JWT tokens are generated for authenticated users.
- JWT middleware validates protected requests and extracts the logged-in user's ID from the token.

### 2. Protected Event APIs

- Public event routes:

```text
GET /api/events
GET /api/events/{id}
```

- Protected event routes:

```text
POST /api/events
PUT /api/events/{id}
DELETE /api/events/{id}
```

- Protected routes require a valid JWT bearer token.
- Event create/update/delete logic is organized through handler, service, and repository layers.
- Event not-found behavior now returns clearer API responses instead of generic failures.

### 3. RSVP APIs

- RSVP endpoint:

```text
POST /api/events/{id}/rsvp
```

- Cancel RSVP endpoint:

```text
DELETE /api/events/{id}/rsvp
```

- RSVP records are stored in the `rsvps` table with a `(user_id, event_id)` primary key.
- Duplicate RSVP attempts are handled as conflicts.
- RSVP deletion is supported for authenticated users.

### 4. Database Updates

- Added/maintained database tables for:

```text
events
users
rsvps
```

- `users` includes unique emails and password hashes.
- `rsvps` links users to events and cascades deletes when a user or event is removed.

## Backend Testing and Validation

Backend test coverage includes handler and service tests for event operations, including:

- Creating an event successfully.
- Handling create-event database failures.
- Getting events successfully.
- Handling get-events database failures.
- Updating events.
- Testing service-level event creation and retrieval.

Validation command:

```bash
cd BackEnd
go test ./...
```

## How to Run the Project

### 1. Start PostgreSQL

From the repo root:

```bash
docker compose up -d
```

### 2. Start the Backend

```bash
cd BackEnd
go run ./cmd/server
```

The backend runs on:

```text
http://localhost:8080
```

### 3. Start the Frontend

In a second terminal:

```bash
cd FrontEnd
npm install
npm run dev
```

The frontend runs on:

```text
http://localhost:5173
```


