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
