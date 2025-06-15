package service

import (
	"context"
	"errors"

	"runmate_api/internal/entity"
	"runmate_api/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	repo *repository.User
}

func NewUser(repo *repository.User) *User {
	return &User{repo: repo}
}

func (u *User) Create(ctx context.Context, user *entity.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	return u.repo.Create(ctx, user)
}

func (u *User) ListAll(ctx context.Context) ([]*entity.User, error) {
	return u.repo.GetAll(ctx)
}

func (u *User) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *User) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return u.repo.GetByUsername(ctx, username)
}

func (u *User) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.repo.GetByEmail(ctx, email)
}

func (u *User) Update(ctx context.Context, user *entity.User) error {
	currentUser, err := u.GetByID(ctx, user.ID.String())
	if err != nil {
		return err
	}

	if currentUser == nil {
		return ErrUserNotFound
	}

	if err := user.Validate(); err != nil {
		return err
	}

	return u.repo.Update(ctx, user)
}

func (u *User) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
