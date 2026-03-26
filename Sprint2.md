# GatorHive – Sprint 2 Comprehensive Report

## Project Overview

GatorHive is a centralized event discovery and management platform designed to streamline event organization and participation for students and communities.

Sprint 2 focused on:

- Completing full CRUD functionality
- Implementing secure authentication
- Integrating frontend and backend systems
- Adding automated testing (E2E + Unit)
- Strengthening backend architecture and validation

---

# Sprint 2 Demo Videos

 Frontend Demo:  
https://youtu.be/kHt29qN58cE  

 Backend Demo:  
https://youtu.be/ZZrMQJcO2_8  

---

#  System Architecture

## Frontend
- React (Vite)
- React Router
- API abstraction layer
- Cypress (E2E testing)
- Vitest (Unit testing)

## Backend
- Go (Golang)
- Gorilla Mux Router
- PostgreSQL (Primary Database)
- SQLite (In-memory testing)
- JWT Authentication
- Layered architecture (Handler → Service → Repository)

---

#  Frontend – Sprint 2 Work

## 1️⃣ Frontend-Backend Integration

The frontend was fully integrated with backend APIs to enable complete event lifecycle management.

### Integrated APIs

| Method | Endpoint | Purpose |
|--------|----------|----------|
| GET | /events | Fetch all events |
| POST | /events | Create event |
| PUT | /events/{id} | Update event |
| DELETE | /events/{id} | Delete event |

All CRUD operations now persist changes in the database and reflect immediately in the UI.

---

## 2️⃣ Update Event Feature

Users can now:

- Navigate to Edit page
- Modify event details
- Submit changes
- Automatically redirect back to event list
- View updated event information

This ensures complete modification capability from UI to database.

---

## 3️⃣ Delete Event Feature

Users can:

- Delete any event directly from the event listing page
- Immediately see UI update
- Confirm removal from database

Backend ensures proper validation and safe deletion.

---

# Frontend Testing

## 🔹 Cypress End-to-End Testing

Implemented automated E2E tests simulating real user workflows:

- createEvent.cy.js – Tests event creation flow
- updateEvent.cy.js – Tests editing functionality
- deleteEvent.cy.js – Tests deletion behavior
- navigation.cy.js – Tests routing and page navigation

These tests validate:

- Form interactions
- API calls
- Page transitions
- UI updates after backend operations

### ▶ Run Cypress

```bash
cd FrontEnd
npm install cypress --save-dev
npx cypress open