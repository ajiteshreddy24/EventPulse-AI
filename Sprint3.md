# GatorHive - Sprint 3 Completion Report

## Overview

Sprint 3 focused on completing the authenticated event-management flow for GatorHive. The project now supports user login and signup, protected event creation/edit/delete actions, RSVP and cancel-RSVP actions, a monthly calendar view, and a refreshed University of Florida-inspired interface with a gator-themed home page.

The application is split into:

- Frontend: React, Vite, React Router, FullCalendar, Cypress.
- Backend: Go, Gorilla Mux, PostgreSQL, JWT authentication, repository/service/handler layers.

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

## Video Demo Script

## Suggested Video Structure

Target length: 6 to 8 minutes.

Speaker plan:

- Frontend Speaker 1: UI, homepage, navigation, auth pages.
- Frontend Speaker 2: event management, calendar, RSVP frontend flow.
- Backend Speaker 1: backend architecture and authentication APIs.
- Backend Speaker 2: event APIs, RSVP APIs, database, tests.

## Opening - 20 to 30 seconds

Speaker: Frontend Speaker 1

Script:

```text
Hi everyone, this is our Sprint 3 demo for GatorHive, our event discovery and management platform. In this sprint, we focused on completing the authenticated user flow, improving event management, adding RSVP support, integrating a calendar view, and refreshing the user interface with a cleaner gator-themed design. We will walk through both the frontend and backend accomplishments, then show how the features connect end to end.
```

Suggested screen recording:

- Show the project running in the browser.
- Start on the home page.
- Slowly move through the navigation bar.

## Frontend Speaker 1 - UI, Navigation, Authentication

Duration: 1.5 to 2 minutes.

Script:

```text
I will start with the frontend interface. For Sprint 3, we refreshed the GatorHive home page with a stronger visual identity and a gator-inspired wallpaper. The design uses a blue, orange, white, and green palette, with a hero section that explains the purpose of the app and gives users quick actions to explore events or create a new event.

The navigation bar now provides direct access to Events, Calendar, Create Event, Login, and Signup. Authentication is handled through a React context provider. When users log in or sign up, the JWT token is stored in localStorage, and the app uses that token for protected requests.

The create and edit routes are protected. If a user tries to access protected event actions without being logged in, the frontend redirects them to the login page. This keeps the user experience clear and keeps protected operations behind authentication.
```

Suggested screen recording:

- Show home page hero.
- Click Events, Calendar, Create Event.
- If logged out, show Create Event redirecting to Login.
- Show Signup and Login pages.
- Login with a test user account.
- Show the navbar updating after login.

## Frontend Speaker 2 - Events, Calendar, RSVP Flow

Duration: 1.5 to 2 minutes.

Script:

```text
Next, I will cover the event features on the frontend. The Events page fetches events from the backend and displays them as cards with the title, description, location, and date. When a user is logged in, they can edit or delete events directly from the event card.

We also added RSVP behavior. If the user has not RSVPed to an event, the button shows RSVP. If the user has already RSVPed, the button changes to Cancel RSVP. After a user RSVPs or cancels, the frontend refreshes the event list so the UI stays synchronized with the backend.

We also integrated FullCalendar for the Calendar page. The calendar pulls the same event data from the backend and displays events in a monthly view, giving users an easier way to browse events by date.
```

Suggested screen recording:

- Open Events page.
- Show an event card.
- Click RSVP, then show the button/state update.
- Click Cancel RSVP if available.
- Click Edit and update an event.
- Return to Events and show the updated event.
- Click Calendar and show events in the monthly grid.

## Backend Speaker 1 - Architecture and Authentication

Duration: 1.5 to 2 minutes.

Script:

```text
Now I will explain the backend work for Sprint 3. The backend is written in Go and uses Gorilla Mux for routing. The code is organized into handlers, services, repositories, models, database setup, and authentication middleware.

For authentication, we implemented register, login, and current-user endpoints under the /api/auth route. Registration stores user records with hashed passwords using bcrypt. Login verifies credentials and returns a JWT token. The /api/auth/me endpoint uses the JWT middleware to identify the current user.

Protected routes use the RequireAuth middleware. This middleware reads the Authorization header, validates the bearer token, extracts the user ID from the token subject, and places it into the request context. Event create, update, delete, RSVP, and cancel RSVP all require this authentication layer.
```

Suggested screen recording:

- Show `BackEnd/cmd/server/main.go`.
- Point to `/api/auth/register`, `/api/auth/login`, and `/api/auth/me`.
- Show `internal/auth/service/auth_service.go`.
- Show JWT generation and bcrypt password hashing.
- Show `internal/auth/middleware/auth_middleware.go`.

## Backend Speaker 2 - Events, RSVP, Database, Tests

Duration: 1.5 to 2 minutes.

Script:

```text
For the event and RSVP backend work, we maintain event routes under /api/events. Public users can fetch all events or fetch a single event by ID. Authenticated users can create, update, and delete events.

Sprint 3 also includes RSVP support. The backend exposes POST /api/events/{id}/rsvp to RSVP and DELETE /api/events/{id}/rsvp to cancel an RSVP. The database has an rsvps table that links users and events, with a primary key across user_id and event_id. This prevents duplicate RSVPs for the same user and event.

We also improved error handling. For example, missing events return not found responses, and duplicate RSVP attempts are treated as conflicts. For validation, we added and ran Go tests around event handlers and services, and we verified the backend with go test ./....
```

Suggested screen recording:

- Show `BackEnd/internal/db/db.go` and the `events`, `users`, and `rsvps` migrations.
- Show `BackEnd/internal/handlers/event_handler.go`.
- Show `BackEnd/internal/queries/event_queries.go`.
- Show backend test files:
  - `BackEnd/internal/handlers/event_api_test.go`
  - `BackEnd/internal/handlers/event_handler_test.go`
  - `BackEnd/internal/service/event_service_test.go`
- Run or show:

```bash
cd BackEnd
go test ./...
```

## Closing - 20 to 30 seconds

Speaker: Frontend Speaker 1 or Backend Speaker 1

Script:

```text
To summarize, Sprint 3 completed the main authenticated event workflow for GatorHive. Users can register, log in, create and manage events, RSVP or cancel RSVP, and view events on a calendar. The backend now supports JWT-protected APIs and RSVP persistence, while the frontend provides a cleaner gator-themed interface and connected user flow. These changes make the application feel much more complete and ready for future recommendation or AI concierge features.
```

## Demo Checklist

- Start PostgreSQL with Docker.
- Start backend with `go run ./cmd/server`.
- Start frontend with `npm run dev`.
- Confirm login/signup works.
- Confirm create event works.
- Confirm edit event works.
- Confirm delete event works.
- Confirm RSVP and cancel RSVP work.
- Confirm calendar view loads events.
- Show frontend lint/build success.
- Show backend test success.

