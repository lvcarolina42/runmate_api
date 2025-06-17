package repository

import (
	"context"
	"fmt"

	"runmate_api/internal/entity"

	"gorm.io/gorm"
)

type Event struct {
	db *gorm.DB
}

func NewEvent(db *gorm.DB) *Event {
	return &Event{db: db}
}

func (e *Event) Create(ctx context.Context, event *entity.Event) error {
	result := e.db.WithContext(ctx).Create(event)
	if result.Error != nil {
		return fmt.Errorf("failed to create event: %v", result.Error)
	}

	return nil
}

func (e *Event) GetByID(ctx context.Context, id string) (*entity.Event, error) {
	var event entity.Event
	result := e.db.WithContext(ctx).Preload("Users").Where("id = ?", id).First(&event)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get event %s: %v", id, result.Error)
	}

	return &event, nil
}

func (e *Event) GetAllActive(ctx context.Context) ([]*entity.Event, error) {
	var events []*entity.Event
	result := e.db.WithContext(ctx).Preload("Users").Where("date >= NOW()").Find(&events)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get active events: %v", result.Error)
	}

	return events, nil
}

func (e *Event) Update(ctx context.Context, event *entity.Event) error {
	result := e.db.WithContext(ctx).Save(event)
	if result.Error != nil {
		return fmt.Errorf("failed to update event: %v", result.Error)
	}

	return nil
}

func (e *Event) GetAllActiveWithoutUser(ctx context.Context, user *entity.User) ([]*entity.Event, error) {
	var events []*entity.Event
	err := e.db.WithContext(ctx).
		Preload("Users").
		Table("events").
		Select("events.*").
		Joins("LEFT JOIN user_events ON events.id = user_events.event_id AND user_events.user_id = ?", user.ID).
		Where("user_events.event_id IS NULL AND date >= NOW()").
		Find(&events).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %v", err)
	}

	return events, nil
}

func (e *Event) GetAllActiveByUser(ctx context.Context, user *entity.User) ([]*entity.Event, error) {
	var events []*entity.Event
	err := e.db.WithContext(ctx).Model(&user).Where("date >= NOW()").Association("Events").Find(&events)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %s active events: %v", user.ID.String(), err)
	}

	return events, nil
}

func (e *Event) GetAllByUser(ctx context.Context, user *entity.User) ([]*entity.Event, error) {
	var events []*entity.Event
	err := e.db.WithContext(ctx).Model(&user).Association("Events").Find(&events)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %s events: %v", user.ID.String(), err)
	}

	return events, nil
}

func (e *Event) AddUser(ctx context.Context, event *entity.Event, user *entity.User) error {
	err := e.db.WithContext(ctx).Model(&event).Association("Users").Append(user)
	if err != nil {
		return fmt.Errorf("failed to add user to event: %v", err)
	}

	return nil
}

func (e *Event) RemoveUser(ctx context.Context, event *entity.Event, user *entity.User) error {
	err := e.db.WithContext(ctx).Model(&event).Association("Users").Delete(user)
	if err != nil {
		return fmt.Errorf("failed to remove user from event: %v", err)
	}

	return nil
}
