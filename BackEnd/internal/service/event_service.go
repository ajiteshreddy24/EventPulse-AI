package service

import (
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/models"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/queries"
)

type EventService struct {
	Repo *queries.EventRepository
}

func (s *EventService) CreateEvent(e *models.Event) error {
	return s.Repo.Create(e)
}

func (s *EventService) GetEvents() ([]models.Event, error) {
	return s.Repo.GetAll()
}

func (s *EventService) GetEventByID(id int) (*models.Event, error) {
	return s.Repo.GetEventByID(id)
}

func (s *EventService) UpdateEvent(e *models.Event) error {
	return s.Repo.Update(e)
}

func (s *EventService) DeleteEvent(id int) error {
	return s.Repo.Delete(id)
}

func (s *EventService) RSVP(userID, eventID int) error {
	return s.Repo.AddRSVP(userID, eventID)
}

func (s *EventService) CancelRSVP(userID, eventID int) error {
	return s.Repo.RemoveRSVP(userID, eventID)
}