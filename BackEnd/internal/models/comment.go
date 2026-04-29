package models

import "time"

type Comment struct {
	ID        int         `json:"id"`
	EventID   int         `json:"event_id"`
	UserID    int         `json:"user_id"`
	Content   string      `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
	User      CommentUser `json:"user"`
}

type CommentUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
