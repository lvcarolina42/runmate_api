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

func (u *User) Authenticate(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := u.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	if user.Password != password {
		return nil, nil
	}

	return user, nil
}

func (u *User) AddFriend(ctx context.Context, userID, friendID string) error {
	user, err := u.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	friend, err := u.GetByID(ctx, friendID)
	if err != nil {
		return err
	}

	if friend == nil {
		return ErrUserNotFound
	}

	if user.ID == friend.ID {
		return errors.New("user and friend are the same")
	}

	return u.repo.CreateFriend(ctx, user, friend)
}

func (u *User) ListFriends(ctx context.Context, userID string) ([]*entity.User, error) {
	user, err := u.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return u.repo.ListFriends(ctx, user)
}

func (u *User) RemoveFriend(ctx context.Context, userID, friendID string) error {
	user, err := u.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	friend, err := u.GetByID(ctx, friendID)
	if err != nil {
		return err
	}

	if friend == nil {
		return ErrUserNotFound
	}

	return u.repo.DeleteFriend(ctx, user, friend)
}
