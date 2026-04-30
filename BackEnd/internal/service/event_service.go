package service

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"sort"
	"strings"
	"time"

	authQueries "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/queries"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/models"
	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/queries"
)

type EventService struct {
	Repo     *queries.EventRepository
	UserRepo *authQueries.UserRepository
}

var sendEmail = defaultSendEmail

func (s *EventService) CreateEvent(e *models.Event) error {
	return s.Repo.Create(e)
}

func (s *EventService) GetEvents(userID *int) ([]models.Event, error) {
	return s.Repo.GetAll(userID)
}

func (s *EventService) GetEventByID(id int, userID *int) (*models.Event, error) {
	return s.Repo.GetEventByID(id, userID)
}

func (s *EventService) UpdateEvent(e *models.Event) error {
	return s.Repo.Update(e)
}

func (s *EventService) DeleteEvent(id int) error {
	return s.Repo.Delete(id)
}

func (s *EventService) RSVP(userID, eventID int) error {
	event, err := s.Repo.GetEventByID(eventID, &userID)
	if err != nil {
		return err
	}
	if event.UserHasRSVP {
		return queries.ErrRSVPAlreadyExists
	}
	if event.Capacity > 0 && event.RSVPCount >= event.Capacity {
		return queries.ErrEventFull
	}

	err = s.Repo.AddRSVP(userID, eventID)
	if err != nil {
		return err
	}

	_ = s.Repo.RemoveFromWaitlist(userID, eventID)
	go s.sendRSVPNotification(userID, eventID, "RSVP confirmed", "You have successfully RSVPed to")
	return nil
}

func (s *EventService) JoinWaitlist(userID, eventID int) error {
	event, err := s.Repo.GetEventByID(eventID, &userID)
	if err != nil {
		return err
	}
	if event.UserHasRSVP {
		return queries.ErrRSVPAlreadyExists
	}
	if event.UserOnWaitlist {
		return queries.ErrAlreadyOnWaitlist
	}
	if event.Capacity <= 0 || event.RSVPCount < event.Capacity {
		return queries.ErrEventFull
	}

	return s.Repo.AddToWaitlist(userID, eventID)
}

func (s *EventService) CancelRSVP(userID, eventID int) error {
	event, err := s.Repo.GetEventByID(eventID, &userID)
	if err != nil {
		return err
	}

	err = s.Repo.RemoveRSVP(userID, eventID)
	if err != nil {
		return err
	}

	go s.sendRSVPNotification(userID, eventID, "RSVP cancelled", "Your RSVP has been cancelled for")

	promotedUserID, err := s.Repo.PromoteNextWaitlistedUserWithResult(eventID)
	if err != nil {
		return err
	}

	if promotedUserID != nil {
		go s.sendPromotionNotification(*promotedUserID, *event)
	}

	return nil
}

func (s *EventService) CreateComment(userID, eventID int, content string) (*models.Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("comment is required")
	}

	if _, err := s.Repo.GetEventByID(eventID, nil); err != nil {
		return nil, err
	}

	return s.Repo.CreateComment(eventID, userID, content)
}

func (s *EventService) GetComments(eventID int) ([]models.Comment, error) {
	return s.Repo.GetCommentsByEventID(eventID)
}

func (s *EventService) DeleteComment(commentID int) error {
	return s.Repo.DeleteComment(commentID)
}

func (s *EventService) GetRecommendations(userID int) ([]models.Event, error) {
	if s.UserRepo == nil {
		return []models.Event{}, nil
	}

	user, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if len(user.Interests) == 0 {
		return []models.Event{}, nil
	}

	events, err := s.Repo.GetAll(&userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var candidates []models.Event
	for _, event := range events {
		if event.UserHasRSVP {
			continue
		}

		matches := matchedInterests(event, user.Interests)
		if len(matches) == 0 {
			continue
		}

		score := float64(len(matches))*5 + float64(event.RSVPCount)*0.35
		reason := fmt.Sprintf("matches your interest in %s", matches[0])

		if event.EventDate.After(now) {
			hoursAway := event.EventDate.Sub(now).Hours()
			if hoursAway < 72 {
				score += 2
			}
		}

		if event.Capacity == 0 || event.RSVPCount < event.Capacity {
			score += 1
		}

		event.RecommendationScore = score
		event.RecommendationReason = reason

		candidates = append(candidates, event)
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].RecommendationScore == candidates[j].RecommendationScore {
			return candidates[i].EventDate.Before(candidates[j].EventDate)
		}
		return candidates[i].RecommendationScore > candidates[j].RecommendationScore
	})

	if len(candidates) > 5 {
		candidates = candidates[:5]
	}

	return candidates, nil
}

func (s *EventService) GetRSVPedEvents(userID int) ([]models.Event, error) {
	return s.Repo.GetRSVPedEventsByUser(userID)
}

func recommendationTokens(event models.Event) []string {
	fields := []string{event.Title, event.Description, event.Location}
	var tokens []string
	for _, field := range fields {
		for _, token := range strings.Fields(strings.ToLower(field)) {
			cleaned := strings.Trim(token, ".,!?;:-_/()[]{}\"'")
			if len(cleaned) >= 4 {
				tokens = append(tokens, cleaned)
			}
		}
	}

	return tokens
}

func matchedInterests(event models.Event, interests []string) []string {
	tokens := recommendationTokens(event)
	tokenSet := map[string]struct{}{}
	for _, token := range tokens {
		tokenSet[token] = struct{}{}
	}

	var matches []string
	for _, interest := range interests {
		if eventMatchesInterest(tokenSet, interest) {
			matches = append(matches, interest)
		}
	}

	return matches
}

func eventMatchesInterest(tokenSet map[string]struct{}, interest string) bool {
	for _, keyword := range interestKeywords(interest) {
		if _, exists := tokenSet[keyword]; exists {
			return true
		}
	}

	return false
}

func interestKeywords(interest string) []string {
	keywords := map[string][]string{
		"technology":       {"technology", "tech", "ai", "coding", "code", "hackathon", "software", "engineering", "robotics", "startup"},
		"sports":           {"sports", "sport", "game", "match", "soccer", "football", "basketball", "fitness", "athletics", "tournament"},
		"music":            {"music", "concert", "band", "choir", "dj", "jam", "performance", "festival"},
		"art":              {"art", "design", "gallery", "painting", "creative", "photography", "theater", "film"},
		"networking":       {"networking", "mixer", "social", "meetup", "community", "connect", "alumni"},
		"career":           {"career", "resume", "interview", "recruiting", "professional", "leadership", "internship", "jobs"},
		"entrepreneurship": {"entrepreneurship", "startup", "pitch", "venture", "founder", "business", "innovation"},
		"gaming":           {"gaming", "game", "esports", "nintendo", "xbox", "playstation", "tournament", "board"},
		"community":        {"community", "service", "volunteer", "charity", "outreach", "cultural", "club"},
		"wellness":         {"wellness", "yoga", "mindfulness", "health", "meditation", "selfcare", "fitness"},
	}

	if words, ok := keywords[interest]; ok {
		return words
	}

	return []string{interest}
}

func (s *EventService) sendRSVPNotification(userID, eventID int, subject, prefix string) {
	user, err := s.Repo.GetUserContactByID(userID)
	if err != nil {
		return
	}

	event, err := s.Repo.GetEventByID(eventID, nil)
	if err != nil {
		return
	}

	body := fmt.Sprintf(
		"Hello %s,\n\n%s %s at %s on %s.\n\nThanks,\nGatorHive",
		user.Name,
		prefix,
		event.Title,
		event.Location,
		event.EventDate.Format(time.RFC1123),
	)

	_ = sendEmail(user.Email, subject, body)
}

func (s *EventService) sendPromotionNotification(userID int, event models.Event) {
	user, err := s.Repo.GetUserContactByID(userID)
	if err != nil {
		return
	}

	body := fmt.Sprintf(
		"Hello %s,\n\nA seat opened up and you have been moved from the waitlist to RSVP for %s at %s on %s.\n\nThanks,\nGatorHive",
		user.Name,
		event.Title,
		event.Location,
		event.EventDate.Format(time.RFC1123),
	)

	_ = sendEmail(user.Email, "Moved from waitlist", body)
}

func defaultSendEmail(to, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	fromEmail := os.Getenv("SMTP_FROM")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" || fromEmail == "" {
		return nil
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	message := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{to}, message)
}
