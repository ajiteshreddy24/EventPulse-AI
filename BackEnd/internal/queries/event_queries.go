package queries

import (
	"database/sql"

	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/models"
)

type EventRepository struct {
	DB *sql.DB
}

func (r *EventRepository) Create(event *models.Event) error {
	query := `
	INSERT INTO events (title, description, location, event_date)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at`

	return r.DB.QueryRow(
		query,
		event.Title,
		event.Description,
		event.Location,
		event.EventDate,
	).Scan(&event.ID, &event.CreatedAt)
}

func (r *EventRepository) Update(event *models.Event) error {
	query := `
	UPDATE events
	SET title = $1,
	    description = $2,
	    location = $3,
	    event_date = $4
	WHERE id = $5
	RETURNING created_at;
	`

	return r.DB.QueryRow(
		query,
		event.Title,
		event.Description,
		event.Location,
		event.EventDate,
		event.ID,
	).Scan(&event.CreatedAt)
}

func (r *EventRepository) GetAll() ([]models.Event, error) {
	rows, err := r.DB.Query(`SELECT id, title, description, location, event_date, created_at FROM events ORDER BY event_date`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.EventDate, &e.CreatedAt)
		events = append(events, e)
	}
	return events, nil
}
