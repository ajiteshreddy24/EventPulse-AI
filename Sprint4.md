# EventPulse – Sprint 4 Completion Report

## Project Overview
EventPulse is a full-stack event discovery and event management platform that enables users to browse, create, and engage with events in a centralized system. 

Sprint 4 focused on improving **user engagement, interactivity, and personalization** by introducing search, discussions, waitlist support, and AI-based recommendations.

---
### Frontend Demo
**https://youtu.be/DQPLaBe3d38**

---

###  Backend Demo
**https://youtu.be/kdV4xjiC4DQ**

---


## Frontend (Sprint 4)

During Sprint 4, we enhanced the user experience by implementing new interactive features including **Search, Comments/Discussions, Waitlist, and AI Recommendations**, along with improved testing coverage. :contentReference[oaicite:0]{index=0}

### Features Implemented

### 1.  Search Functionality
- Added search input to filter events by:
  - Title  
  - Description  
  - Location  
- Enables faster discovery of relevant events.

---

### 2. Comments / Discussions
- Users can:
  - View comments
  - Toggle comments (Show/Hide)
  - Add new comments
  - Delete comments
- Uses API:
  - `GET /events/{id}/comments`
  - `POST /events/{id}/comments`

---

### 3. Waitlist Functionality
- If event reaches capacity:
  - Users can join waitlist
  - Waitlist count is displayed
  - Button is disabled if already waitlisted

---

### 4. AI Recommendations
- Personalized recommendations based on user interests
- Matches:
  - Event title
  - Description
  - Location
- Displays most relevant events for the user

---

### 5. Profile Interest Management
- Users can:
  - Add/remove interests
  - Save preferences
- Used for improving recommendation accuracy

---

## Frontend Testing

### Unit Tests (Vitest + React Testing Library)
- `EventList.test.jsx` → RSVP logic, rendering, capacity
- `SearchEventList.test.jsx` → search filtering
- `EventComments.test.jsx` → comments toggle behavior
- `WaitList.test.jsx` → waitlist button behavior

---

### Cypress E2E Tests
- `waitlist.cy.js` → full-capacity waitlist flow
- `comments.cy.js` → comments interaction
- `search.cy.js` → event filtering

**Validates:**
- Search functionality  
- Comment interactions  
- Waitlist behavior  

---


## ⚙️ Backend (Sprint 4)

Sprint 4 backend focused on implementing **comments, waitlist, email notifications, and personalization logic**.

---

### Features Implemented

### 1. 💬 Comment System
- **POST /events/{id}/comments**
  - Requires JWT authentication
  - Creates comment linked to user and event

- **GET /events/{id}/comments**
  - Fetches all comments for an event

- **DELETE /comments/{id}**
  - Deletes a comment

---

### 2. Waitlist System
- **POST /events/{id}/waitlist**
- Adds user when event is full
- Prevents duplicate entries
- Automatically promotes users when slots open

---

### 3. Email Notifications
- RSVP Confirmation Email
- RSVP Cancellation Email
- Implemented using:
  - Go routines (async processing)
  - Gmail SMTP (`net/smtp`)

---

### 4. RSVP Enhancements
- RSVP triggers email
- Cancel RSVP:
  - Sends cancellation email
  - Promotes waitlisted users
- Includes event capacity handling

---

### 5. Recommendation Engine
- Stores user interests
- Matches interests with event metadata
- Scores and ranks events
- Returns personalized recommendations

---

## Backend Unit Tests

### 1. TestCreateComment
- Validates comment creation with JWT
- Checks user_id, event_id, and content

---

### 2. TestGetComments
- Returns ordered comments list
- Ensures correct response structure

---

### 3. TestDeleteComment
- Deletes comment successfully
- Confirms removal from DB

---

### 4. TestJoinWaitlist_Success
- Adds user to waitlist
- Ensures DB entry is created

---

### 5. TestRSVP_EmailSent
- Verifies email trigger on RSVP
- Checks:
  - Recipient
  - Subject
  - Invocation

---



##  Database Changes

- Added `capacity` field to events
- Created `comments` table
- Created `waitlist_entries` table
- Added `interests` field to users

---

## How to Run the Project

### 1. Start Database
```bash
docker compose up -d
