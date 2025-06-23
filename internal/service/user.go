package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"runmate_api/internal/entity"
	"runmate_api/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	activityRepo *repository.Activity
	userRepo     *repository.User
}

func NewUser(activityRepo *repository.Activity, userRepo *repository.User) *User {
	return &User{activityRepo: activityRepo, userRepo: userRepo}
}

func (u *User) enrichUserWithWeekActivities(ctx context.Context, user *entity.User) error {
	today := time.Now()
	todayWithoutHour := today.Truncate(24 * time.Hour)
	offset := (int(time.Sunday) - int(todayWithoutHour.Weekday()) - 7) % 7
	start := todayWithoutHour.Add(time.Duration(offset*24) * time.Hour)

	activities, err := u.activityRepo.GetByUserIDAndDateRange(ctx, user.ID.String(), start, today)
	if err != nil {
		return err
	}

	fmt.Println("activities", activities)
	weekActivities := make(map[string]*entity.UserDayActitivy, (-offset + 1))
	for _, activity := range activities {
		dateKey := activity.Date.Format("2006-01-02")
		if dayActivity, ok := weekActivities[dateKey]; ok {
			dayActivity.Distance += activity.Distance
		} else {
			weekActivities[dateKey] = &entity.UserDayActitivy{
				Date:     activity.Date,
				Distance: activity.Distance,
			}
		}
	}

	fmt.Println("weekActivities", weekActivities)
	user.WeekActivities = weekActivities
	return nil
}

func (u *User) Create(ctx context.Context, user *entity.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	return u.userRepo.Create(ctx, user)
}

func (u *User) ListAll(ctx context.Context) ([]*entity.User, error) {
	return u.userRepo.GetAll(ctx)
}

func (u *User) ListAllNonFriends(ctx context.Context, userID string) ([]*entity.User, error) {
	user, err := u.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return u.userRepo.GetAllNonFriends(ctx, user)
}

func (u *User) GetByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.enrichUserWithWeekActivities(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return u.userRepo.GetByUsername(ctx, username)
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

	return u.userRepo.Update(ctx, user)
}

func (u *User) UpdateFCMToken(ctx context.Context, userID string, fcmToken string) error {
	user, err := u.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	user.FCMToken = fcmToken
	return u.userRepo.Update(ctx, user)
}

func (u *User) UpdateGoal(ctx context.Context, userID string, days, distance *int) error {
	user, err := u.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	user.GoalDays = days
	user.GoalDailyDistance = distance
	return u.userRepo.Update(ctx, user)
}

func (u *User) Delete(ctx context.Context, id string) error {
	return u.userRepo.Delete(ctx, id)
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

	return u.userRepo.CreateFriend(ctx, user, friend)
}

func (u *User) ListFriends(ctx context.Context, userID string) ([]*entity.User, error) {
	user, err := u.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return u.userRepo.ListFriends(ctx, user)
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

	return u.userRepo.DeleteFriend(ctx, user, friend)
}
