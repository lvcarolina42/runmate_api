package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"runmate_api/internal/entity"
)

type Activity struct {
	db *gorm.DB
}

func NewActivity(db *gorm.DB) *Activity {
	return &Activity{db: db}
}

func (a *Activity) Create(ctx context.Context, activity *entity.Activity) error {
	result := a.db.WithContext(ctx).Create(activity)
	if result.Error != nil {
		return fmt.Errorf("failed to create activity: %v", result.Error)
	}

	return nil
}

func (a *Activity) GetAll(ctx context.Context) ([]*entity.Activity, error) {
	var activities []*entity.Activity
	result := a.db.
		WithContext(ctx).
		Preload("Coordinates", func(db *gorm.DB) *gorm.DB {
			return db.Order("coordinates.order ASC")
		}).
		Preload("User").
		Find(&activities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get activities: %v", result.Error)
	}

	return activities, nil
}

func (a *Activity) GetByUserID(ctx context.Context, userID string) ([]*entity.Activity, error) {
	var activities []*entity.Activity
	result := a.db.
		WithContext(ctx).
		Preload("Coordinates", func(db *gorm.DB) *gorm.DB {
			return db.Order("coordinates.order ASC")
		}).
		Preload("User").
		Where("user_id = ?", userID).
		Find(&activities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get activities for user %s: %v", userID, result.Error)
	}

	return activities, nil
}

func (a *Activity) GetByUserIDAndDateRange(ctx context.Context, userID string, start, end time.Time) ([]*entity.Activity, error) {
	var activities []*entity.Activity
	result := a.db.
		WithContext(ctx).
		Preload("Coordinates", func(db *gorm.DB) *gorm.DB {
			return db.Order("coordinates.order ASC")
		}).
		Preload("User").
		Where("user_id = ? AND date >= ? AND date <= ?", userID, start, end).
		Find(&activities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get activities for user %s in date range %s - %s: %v", userID, start, end, result.Error)
	}

	return activities, nil
}

func (a *Activity) Delete(ctx context.Context, id string) error {
	result := a.db.WithContext(ctx).Select(clause.Associations).Where("id = ?", id).Delete(&entity.Activity{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete activity %s: %v", id, result.Error)
	}

	return nil
}
