package service

import (
	"context"

	"runmate_api/internal/entity"
	"runmate_api/internal/repository"
)

type Activity struct {
	repo *repository.Activity
}

func NewActivity(repo *repository.Activity) *Activity {
	return &Activity{repo: repo}
}

func (a *Activity) Create(ctx context.Context, activity *entity.Activity) error {
	for i, coordinate := range activity.Coordinates {
		coordinate.Order = int(i)
	}

	return a.repo.Create(ctx, activity)
}

func (a *Activity) ListAll(ctx context.Context) ([]*entity.Activity, error) {
	return a.repo.GetAll(ctx)
}

func (a *Activity) ListByUser(ctx context.Context, userID string) ([]*entity.Activity, error) {
	return a.repo.GetByUserID(ctx, userID)
}

func (a *Activity) Delete(ctx context.Context, id string) error {
	return a.repo.Delete(ctx, id)
}
