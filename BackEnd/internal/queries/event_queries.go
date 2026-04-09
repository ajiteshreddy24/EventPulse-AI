package queries

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/models"
)

type EventRepository struct {
	DB *sql.DB
}

var ErrEventNotFound = errors.New("event not found")
var ErrRSVPAlreadyExists = errors.New("rsvp already exists")

func (r *EventRepository) Create(e *models.Event) error {
	return r.DB.QueryRow(`
		INSERT INTO events (title, description, location, event_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`,
		e.Title,
		e.Description,
		e.Location,
		e.EventDate,
	).Scan(&e.ID, &e.CreatedAt)
}

func (r *EventRepository) GetAll() ([]models.Event, error) {
	rows, err := r.DB.Query(`
		SELECT id, title, description, location, event_date, created_at
		FROM events
		ORDER BY event_date ASC
	`)
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
			&e.CreatedAt,
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

func (r *EventRepository) GetEventByID(id int) (*models.Event, error) {
	var e models.Event

	err := r.DB.QueryRow(`
		SELECT id, title, description, location, event_date, created_at
		FROM events
		WHERE id = $1
	`, id).Scan(
		&e.ID,
		&e.Title,
		&e.Description,
		&e.Location,
		&e.EventDate,
		&e.CreatedAt,
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
	result, err := r.DB.Exec(`
		UPDATE events
		SET title = $1,
		    description = $2,
		    location = $3,
		    event_date = $4
		WHERE id = $5
	`,
		e.Title,
		e.Description,
		e.Location,
		e.EventDate,
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
