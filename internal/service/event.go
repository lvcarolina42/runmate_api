package service

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"math/rand"
	"slices"

	"runmate_api/internal/entity"
	"runmate_api/internal/firebase"
	"runmate_api/internal/repository"
)

var (
	ErrEventNotFound = errors.New("event not found")

	newEventNotificationTitles = []string{"Novo evento na área!", "Novo evento para você!"}
)

func newEventNotification(userName, eventTitle string) *firebase.Notification {
	return &firebase.Notification{
		Title: newEventNotificationTitles[rand.Intn(len(newEventNotificationTitles))],
		Body:  fmt.Sprintf("%s criou um novo evento: %s", userName, eventTitle),
	}
}

type Event struct {
	eventRepo *repository.Event
	userRepo  *repository.User

	firebaseClient *firebase.Client
}

func NewEvent(eventRepo *repository.Event, userRepo *repository.User, firebaseClient *firebase.Client) *Event {
	return &Event{
		eventRepo: eventRepo,
		userRepo:  userRepo,

		firebaseClient: firebaseClient,
	}
}

func (c *Event) Create(ctx context.Context, event *entity.Event) error {
	owner, err := c.userRepo.GetByID(ctx, event.CreatedBy.String())
	if err != nil {
		return err
	}

	if owner == nil {
		return ErrUserNotFound
	}

	event.Users = []*entity.User{owner}
	err = c.eventRepo.Create(ctx, event)
	if err != nil {
		return err
	}

	users, err := c.userRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	tokens := make(map[string]any, len(users)-1)
	for _, user := range users {
		if user.ID == owner.ID || user.FCMToken == "" {
			continue
		}

		tokens[user.FCMToken] = struct{}{}
	}

	notification := newEventNotification(owner.Name, event.Title)
	return c.firebaseClient.SendNotification(ctx, notification, slices.Collect(maps.Keys(tokens)))
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
