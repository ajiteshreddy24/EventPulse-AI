# Event Model Documentation

## Overview
The `Event` struct represents an event in the system, including its details, user interaction status, and recommendation-related information.

---

## Struct Definition

```go
type Event struct {
	ID                   int       `json:"id"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	Location             string    `json:"location"`
	EventDate            time.Time `json:"event_date"`
	Capacity             int       `json:"capacity"`
	CreatedAt            time.Time `json:"created_at"`
	RSVPCount            int       `json:"rsvpCount"`
	UserHasRSVP          bool      `json:"userHasRSVP"`
	WaitlistCount        int       `json:"waitlistCount"`
	UserOnWaitlist       bool      `json:"userOnWaitlist"`
	RecommendationReason string    `json:"recommendationReason,omitempty"`
	RecommendationScore  float64   `json:"recommendationScore,omitempty"`
}
```

---

## Field Descriptions

### Basic Information
- **ID**: Unique identifier for the event  
- **Title**: Name of the event  
- **Description**: Details about the event  
- **Location**: Where the event takes place  
- **EventDate**: Date and time of the event  
- **Capacity**: Maximum number of attendees  
- **CreatedAt**: Timestamp when the event was created  

---

### RSVP & Participation
- **RSVPCount**: Number of users who RSVP’d  
- **UserHasRSVP**: Whether the current user has RSVP’d  

---

### Waitlist
- **WaitlistCount**: Number of users on the waitlist  
- **UserOnWaitlist**: Whether the current user is on the waitlist  

---

### Recommendation System
- **RecommendationReason**: Why the event is recommended (optional)  
- **RecommendationScore**: Strength of recommendation (optional)  

---

## Example JSON

```json
{
  "id": 1,
  "title": "Tech Meetup",
  "description": "A meetup for developers",
  "location": "New York",
  "event_date": "2026-05-01T18:00:00Z",
  "capacity": 100,
  "created_at": "2026-04-01T12:00:00Z",
  "rsvpCount": 50,
  "userHasRSVP": true,
  "waitlistCount": 10,
  "userOnWaitlist": false
}
```

---

## Notes
- Fields with `omitempty` will not appear in responses if empty  
- `time.Time` fields use ISO 8601 format  
- User-specific fields are dynamically set per request  

---

## Usage
Used in:
- Event creation APIs  
- Event listing APIs  
- RSVP and waitlist handling  
- Recommendation systems  
