# Sprint 1 Report – EventPulse-AI

## Project Overview
Students and communities often miss relevant events due to fragmented announcements and information overload. *EventPulse-AI* is a centralized platform that enables users to create, discover, and manage events. Sprint 1 focused on establishing a solid full-stack foundation using *React, **Go, **PostgreSQL, and **Docker, along with API testing using **Postman*.

---

## User Stories

•⁠  ⁠As a user, I want to create an event so that I can organize and share details with others.
•⁠  ⁠As a user, I want to view a list of upcoming events so that I can decide which ones to attend.
•⁠  ⁠As a user, I want an intuitive navigation system so that I can easily move between pages.
•⁠  ⁠As a user, I want a confirmation message after creating an event so that I know my submission was successful.
•⁠  ⁠As a user, I want a visually appealing homepage so that I feel engaged when visiting the site.

---

## Planned Issues for Sprint 1

### Frontend
•⁠  ⁠Set up React frontend environment
•⁠  ⁠Implement event creation form with controlled inputs
•⁠  ⁠Display a list of upcoming events
•⁠  ⁠Configure page navigation using React Router
•⁠  ⁠Show confirmation message upon successful event creation
•⁠  ⁠Add background image and basic styling for homepage UI

### Backend
•⁠  ⁠Set up backend using Go (Golang)
•⁠  ⁠Configure PostgreSQL database using Docker
•⁠  ⁠Establish database connection and schema
•⁠  ⁠Implement RESTful APIs:
  - Create Event API
  - Get All Events API
  - Update Event API
  - Delete Event API
•⁠  ⁠Test backend APIs using Postman

---

## Successfully Completed Issues

### Frontend
•⁠  ⁠Successfully set up React frontend using Vite
•⁠  ⁠Implemented event creation form with validation
•⁠  ⁠Implemented event listing page displaying data dynamically
•⁠  ⁠Configured navigation between pages using React Router
•⁠  ⁠Displayed confirmation message after successful event creation
•⁠  ⁠Added background image and improved homepage UI styling

### Backend
•⁠  ⁠Successfully set up backend using Go with modular project structure
•⁠  ⁠Configured PostgreSQL database using Docker and Docker Compose
•⁠  ⁠Created database schema for events
•⁠  ⁠Implemented RESTful APIs:
  - *POST /events* – Create a new event
  - *GET /events* – Retrieve all events
  - *PUT /events/{id}* – Update an existing event
•⁠  ⁠Verified API functionality using *Postman* with proper request/response validation
•⁠  ⁠Ensured backend services run independently of frontend

---

## Issues Not Completed and Reasons

•⁠  ⁠*Event deletion API*  
  The Delete Event API was planned but not completed due to time constraints. This feature will be prioritized in the next sprint.

•⁠  ⁠*Frontend event persistence on refresh*  
  Currently, events displayed on the frontend rely on in-memory React state and are lost upon page refresh. Full persistence using backend integration and database storage will be implemented in the next sprint.

---

## Upcoming Sprints

In upcoming sprints, the team plans to:
•⁠  ⁠Complete the Delete Event API
•⁠  ⁠Fully integrate frontend with backend APIs for persistent data storage
•⁠  ⁠Enhance error handling and validation on both frontend and backend
•⁠  ⁠Improve UI/UX and scalability of the application
•⁠  ⁠Begin integrating AI-powered features for personalized event recommendations

---

## Demo Links

•⁠  ⁠*Frontend Video Demo:*  
  https://www.youtube.com/watch?v=8cGDy8C7_6s

•⁠  ⁠*Backend Video Demo:*  
  https://youtu.be/5ZAnOgZp5-s