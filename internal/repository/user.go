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

func (u *User) GetAll(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User
	result := u.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get users: %v", result.Error)
	}

	return users, nil
}

func (u *User) GetAllNonFriends(ctx context.Context, user *entity.User) ([]*entity.User, error) {
	var users []*entity.User
	result := u.db.WithContext(ctx).
		Table("users").
		Select("users.*").
		Joins("LEFT JOIN user_friends ON users.id = user_friends.friend_id AND user_friends.user_id = ?", user.ID).
		Where("user_friends.friend_id IS NULL AND users.id != ?", user.ID).
		Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get users: %v", result.Error)
	}

	return users, nil
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

func (u *User) CreateFriend(ctx context.Context, user, friend *entity.User) error {
	err := u.db.WithContext(ctx).Model(&user).Association("Friends").Append(friend)
	if err != nil {
		return fmt.Errorf("failed to create friend: %v", err)
	}

	return nil
}

func (u *User) ListFriends(ctx context.Context, user *entity.User) ([]*entity.User, error) {
	var friends []*entity.User
	err := u.db.WithContext(ctx).Model(&user).Association("Friends").Find(&friends)
	if err != nil {
		return nil, fmt.Errorf("failed to list friends: %v", err)
	}

	return friends, nil
}

func (u *User) DeleteFriend(ctx context.Context, user, friend *entity.User) error {
	err := u.db.WithContext(ctx).Model(&user).Association("Friends").Delete(friend)
	if err != nil {
		return fmt.Errorf("failed to delete friend: %v", err)
	}

	return nil
}
