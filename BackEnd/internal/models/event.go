package models

import "time"

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
