package model

import (
	"errors"
	"fmt"
	"runmate_api/internal/entity"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDateRequired = errors.New("date is required")
	ErrPastDate     = errors.New("date must be in the future")
)

type Event struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Date         time.Time `json:"date"`
	Finished     bool      `json:"finished"`
	Participants []*User   `json:"participants,omitempty"`
}

func NewEventFromEntity(c *entity.Event) *Event {
	participants := make([]*User, 0, len(c.Users))
	for _, item := range c.Users {
		participants = append(participants, NewUserFromEntity(item))
	}

	return &Event{
		ID:           c.ID.String(),
		Title:        c.Title,
		Date:         c.Date,
		Finished:     c.Date.Before(time.Now()),
		Participants: participants,
	}
}

type CreateEventInput struct {
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
	UserID string    `json:"created_by"`
}

func (c *CreateEventInput) Validate() error {
	if c.Date.IsZero() {
		return ErrDateRequired
	}

	if c.Date.Before(time.Now()) {
		return ErrPastDate
	}

	return nil
}

func (c *CreateEventInput) ToEntity() (*entity.Event, error) {
	userID, err := uuid.Parse(c.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user id: %v", err)
	}

	return &entity.Event{
		Title:     c.Title,
		Date:      c.Date,
		CreatedBy: userID,
	}, nil
}

type JoinQuitEventInput struct {
	UserID  string `json:"user_id"`
	EventID string `json:"event_id"`
}
