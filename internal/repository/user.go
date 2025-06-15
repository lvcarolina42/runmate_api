package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"runmate_api/internal/entity"
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{db: db}
}

func (u *User) Create(ctx context.Context, user *entity.User) error {
	result := u.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %v", result.Error)
	}

	return nil
}

func (u *User) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	result := u.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user %s: %v", id, result.Error)
	}

	return &user, nil
}

func (u *User) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	result := u.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user by username %s: %v", username, result.Error)
	}

	return &user, nil
}

func (u *User) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	result := u.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user by email %s: %v", email, result.Error)
	}

	return &user, nil
}

func (u *User) Update(ctx context.Context, user *entity.User) error {
	result := u.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %v", result.Error)
	}

	return nil
}

func (u *User) Delete(ctx context.Context, id string) error {
	result := u.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user %s: %v", id, result.Error)
	}

	return nil
}
