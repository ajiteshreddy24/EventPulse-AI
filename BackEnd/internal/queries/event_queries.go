package queries

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/models"
)

type EventRepository struct {
	DB *sql.DB
}

var ErrEventNotFound = errors.New("event not found")
var ErrRSVPAlreadyExists = errors.New("rsvp already exists")
var ErrEventFull = errors.New("event is at full capacity")
var ErrAlreadyOnWaitlist = errors.New("user already on waitlist")
var ErrWaitlistNotFound = errors.New("waitlist entry not found")

func (r *EventRepository) Create(e *models.Event) error {
	if e.Capacity <= 0 {
		e.Capacity = 50
	}

	return r.DB.QueryRow(`
		INSERT INTO events (title, description, location, event_date, capacity)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`,
		e.Title,
		e.Description,
		e.Location,
		e.EventDate,
		e.Capacity,
	).Scan(&e.ID, &e.CreatedAt)
}

func (r *EventRepository) GetAll(userID *int) ([]models.Event, error) {
	rows, err := r.DB.Query(`
		SELECT
			e.id,
			e.title,
			e.description,
			e.location,
			e.event_date,
			e.capacity,
			e.created_at,
			COALESCE((SELECT COUNT(*) FROM rsvps r WHERE r.event_id = e.id), 0) AS rsvp_count,
			CASE
				WHEN $1::int IS NULL THEN FALSE
				ELSE EXISTS(
					SELECT 1 FROM rsvps r
					WHERE r.event_id = e.id AND r.user_id = $1
				)
			END AS user_has_rsvp,
			COALESCE((SELECT COUNT(*) FROM waitlist_entries w WHERE w.event_id = e.id), 0) AS waitlist_count,
			CASE
				WHEN $1::int IS NULL THEN FALSE
				ELSE EXISTS(
					SELECT 1 FROM waitlist_entries w
					WHERE w.event_id = e.id AND w.user_id = $1
				)
			END AS user_on_waitlist
		FROM events e
		ORDER BY event_date ASC
	`, nullableUserID(userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var e models.Event
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.Location,
			&e.EventDate,
			&e.Capacity,
			&e.CreatedAt,
			&e.RSVPCount,
			&e.UserHasRSVP,
			&e.WaitlistCount,
			&e.UserOnWaitlist,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *EventRepository) GetEventByID(id int, userID *int) (*models.Event, error) {
	var e models.Event

	err := r.DB.QueryRow(`
		SELECT
			e.id,
			e.title,
			e.description,
			e.location,
			e.event_date,
			e.capacity,
			e.created_at,
			COALESCE((SELECT COUNT(*) FROM rsvps r WHERE r.event_id = e.id), 0) AS rsvp_count,
			CASE
				WHEN $2::int IS NULL THEN FALSE
				ELSE EXISTS(
					SELECT 1 FROM rsvps r
					WHERE r.event_id = e.id AND r.user_id = $2
				)
			END AS user_has_rsvp,
			COALESCE((SELECT COUNT(*) FROM waitlist_entries w WHERE w.event_id = e.id), 0) AS waitlist_count,
			CASE
				WHEN $2::int IS NULL THEN FALSE
				ELSE EXISTS(
					SELECT 1 FROM waitlist_entries w
					WHERE w.event_id = e.id AND w.user_id = $2
				)
			END AS user_on_waitlist
		FROM events e
		WHERE e.id = $1
	`, id, nullableUserID(userID)).Scan(
		&e.ID,
		&e.Title,
		&e.Description,
		&e.Location,
		&e.EventDate,
		&e.Capacity,
		&e.CreatedAt,
		&e.RSVPCount,
		&e.UserHasRSVP,
		&e.WaitlistCount,
		&e.UserOnWaitlist,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEventNotFound
	}
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *EventRepository) Update(e *models.Event) error {
	if e.Capacity <= 0 {
		e.Capacity = 50
	}

	result, err := r.DB.Exec(`
		UPDATE events
		SET title = $1,
		    description = $2,
		    location = $3,
		    event_date = $4,
		    capacity = $5
		WHERE id = $6
	`,
		e.Title,
		e.Description,
		e.Location,
		e.EventDate,
		e.Capacity,
		e.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (r *EventRepository) Delete(id int) error {
	result, err := r.DB.Exec(`DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (r *EventRepository) AddRSVP(userID, eventID int) error {
	_, err := r.DB.Exec(`
		INSERT INTO rsvps (user_id, event_id)
		VALUES ($1, $2)
	`, userID, eventID)
	if err != nil {
		lowerErr := strings.ToLower(err.Error())
		if strings.Contains(lowerErr, "duplicate") || strings.Contains(lowerErr, "unique") {
			return ErrRSVPAlreadyExists
		}
		if strings.Contains(lowerErr, "foreign key") {
			return ErrEventNotFound
		}
	}
	return err
}

func (r *EventRepository) RemoveRSVP(userID, eventID int) error {
	result, err := r.DB.Exec(`
		DELETE FROM rsvps
		WHERE user_id = $1 AND event_id = $2
	`, userID, eventID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (r *EventRepository) AddToWaitlist(userID, eventID int) error {
	_, err := r.DB.Exec(`
		INSERT INTO waitlist_entries (user_id, event_id)
		VALUES ($1, $2)
	`, userID, eventID)
	if err != nil {
		lowerErr := strings.ToLower(err.Error())
		if strings.Contains(lowerErr, "duplicate") || strings.Contains(lowerErr, "unique") {
			return ErrAlreadyOnWaitlist
		}
		if strings.Contains(lowerErr, "foreign key") {
			return ErrEventNotFound
		}
	}
	return err
}

func (r *EventRepository) RemoveFromWaitlist(userID, eventID int) error {
	_, err := r.DB.Exec(`
		DELETE FROM waitlist_entries
		WHERE user_id = $1 AND event_id = $2
	`, userID, eventID)
	return err
}

func (r *EventRepository) PromoteNextWaitlistedUser(eventID int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int
	err = tx.QueryRow(`
		SELECT user_id
		FROM waitlist_entries
		WHERE event_id = $1
		ORDER BY created_at ASC, id ASC
		LIMIT 1
	`, eventID).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO rsvps (user_id, event_id)
		VALUES ($1, $2)
	`, userID, eventID); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		DELETE FROM waitlist_entries
		WHERE user_id = $1 AND event_id = $2
	`, userID, eventID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *EventRepository) PromoteNextWaitlistedUserWithResult(eventID int) (*int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var userID int
	err = tx.QueryRow(`
		SELECT user_id
		FROM waitlist_entries
		WHERE event_id = $1
		ORDER BY created_at ASC, id ASC
		LIMIT 1
	`, eventID).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if _, err := tx.Exec(`
		INSERT INTO rsvps (user_id, event_id)
		VALUES ($1, $2)
	`, userID, eventID); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(`
		DELETE FROM waitlist_entries
		WHERE user_id = $1 AND event_id = $2
	`, userID, eventID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &userID, nil
}

func (r *EventRepository) CreateComment(eventID, userID int, content string) (*models.Comment, error) {
	comment := &models.Comment{}

	err := r.DB.QueryRow(`
		INSERT INTO comments (event_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, event_id, user_id, content, created_at
	`, eventID, userID, content).Scan(
		&comment.ID,
		&comment.EventID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := r.DB.QueryRow(`
		SELECT id, name, email
		FROM users
		WHERE id = $1
	`, userID).Scan(
		&comment.User.ID,
		&comment.User.Name,
		&comment.User.Email,
	); err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *EventRepository) GetCommentsByEventID(eventID int) ([]models.Comment, error) {
	rows, err := r.DB.Query(`
		SELECT
			c.id,
			c.event_id,
			c.user_id,
			c.content,
			c.created_at,
			u.id,
			u.name,
			u.email
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.event_id = $1
		ORDER BY c.created_at ASC, c.id ASC
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.EventID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.User.ID,
			&comment.User.Name,
			&comment.User.Email,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

func (r *EventRepository) DeleteComment(commentID int) error {
	result, err := r.DB.Exec(`DELETE FROM comments WHERE id = $1`, commentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (r *EventRepository) GetRSVPedEventsByUser(userID int) ([]models.Event, error) {
	rows, err := r.DB.Query(`
		SELECT
			e.id,
			e.title,
			e.description,
			e.location,
			e.event_date,
			e.capacity,
			e.created_at,
			COALESCE((SELECT COUNT(*) FROM rsvps r WHERE r.event_id = e.id), 0) AS rsvp_count
		FROM events e
		JOIN rsvps rsvp ON rsvp.event_id = e.id
		WHERE rsvp.user_id = $1
		ORDER BY e.event_date ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Location,
			&event.EventDate,
			&event.Capacity,
			&event.CreatedAt,
			&event.RSVPCount,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

func nullableUserID(userID *int) any {
	if userID == nil {
		return nil
	}
	return *userID
}

func (r *EventRepository) GetUserContactByID(userID int) (models.CommentUser, error) {
	var user models.CommentUser

	err := r.DB.QueryRow(`
		SELECT id, name, email
		FROM users
		WHERE id = $1
	`, userID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return models.CommentUser{}, err
	}

	return user, nil
}

func (r *EventRepository) DebugString(eventID int) string {
	return fmt.Sprintf("event-%d", eventID)
}
