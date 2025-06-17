package service

import (
	"context"
	"errors"

	"runmate_api/internal/entity"
	"runmate_api/internal/repository"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type Event struct {
	eventRepo *repository.Event
	userRepo  *repository.User
}

func NewEvent(eventRepo *repository.Event, userRepo *repository.User) *Event {
	return &Event{eventRepo: eventRepo, userRepo: userRepo}
}

func (c *Event) Create(ctx context.Context, event *entity.Event) error {
	user, err := c.userRepo.GetByID(ctx, event.CreatedBy.String())
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	event.Users = []*entity.User{user}
	return c.eventRepo.Create(ctx, event)
}

func (c *Event) GetByID(ctx context.Context, id string) (*entity.Event, error) {
	return c.eventRepo.GetByID(ctx, id)
}

func (c *Event) ListAllActive(ctx context.Context) ([]*entity.Event, error) {
	return c.eventRepo.GetAllActive(ctx)
}

func (c *Event) ListAllActiveWithoutUserID(ctx context.Context, userID string) ([]*entity.Event, error) {
	user, err := c.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return c.eventRepo.GetAllActiveWithoutUser(ctx, user)
}

func (c *Event) ListAllActiveByUserID(ctx context.Context, userID string) ([]*entity.Event, error) {
	user, err := c.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return c.eventRepo.GetAllActiveByUser(ctx, user)
}

func (c *Event) ListAllByUserID(ctx context.Context, userID string) ([]*entity.Event, error) {
	user, err := c.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return c.eventRepo.GetAllByUser(ctx, user)
}

func (c *Event) Join(ctx context.Context, eventID, userID string) error {
	user, err := c.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	event, err := c.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event == nil {
		return ErrEventNotFound
	}

	return c.eventRepo.AddUser(ctx, event, user)
}

func (c *Event) Quit(ctx context.Context, eventID, userID string) error {
	user, err := c.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	event, err := c.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event == nil {
		return ErrEventNotFound
	}

	return c.eventRepo.RemoveUser(ctx, event, user)
}
