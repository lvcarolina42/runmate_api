package model

import (
	"time"

	"runmate_api/internal/entity"
)

type GoalDayActivity struct {
	Date     string `json:"date"`
	Distance int    `json:"distance"`
}

func newGoalDayActivityFromEntity(activity *entity.UserDayActitivy) *GoalDayActivity {
	return &GoalDayActivity{
		Date:     activity.Date.Format("2006-01-02"),
		Distance: activity.Distance,
	}
}

type Goal struct {
	Days           *int               `json:"days,omitempty"`
	DailyDistance  *int               `json:"daily_distance,omitempty"`
	WeekActivities []*GoalDayActivity `json:"week_activities,omitempty"`
}

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Birthdate   time.Time `json:"birthdate"`
	Role        int8      `json:"role"`
	XP          int       `json:"xp"`
	Level       int       `json:"level"`
	NextLevelXP int       `json:"next_level_xp"`
	Goal        *Goal     `json:"goal,omitempty"`
}

func NewUserFromEntity(user *entity.User) *User {
	weekActivities := make([]*GoalDayActivity, 0, 7)
	for _, activity := range user.WeekActivities {
		weekActivities = append(weekActivities, newGoalDayActivityFromEntity(activity))
	}

	return &User{
		ID:          user.ID.String(),
		Username:    user.Username,
		Name:        user.Name,
		Email:       user.Email,
		Birthdate:   user.Birthdate,
		Role:        user.Role,
		XP:          user.XP,
		Level:       user.CurrentLevel(),
		NextLevelXP: user.NextLevelXP(),
		Goal: &Goal{
			Days:           user.GoalDays,
			DailyDistance:  user.GoalDailyDistance,
			WeekActivities: weekActivities,
		},
	}
}

type CreateUserInput struct {
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Birthdate time.Time `json:"birthdate"`
}

func (c *CreateUserInput) ToEntity() *entity.User {
	return &entity.User{
		Username:  c.Username,
		Name:      c.Name,
		Email:     c.Email,
		Password:  c.Password,
		Birthdate: c.Birthdate,
		Role:      0,
	}
}

type FriendInput struct {
	UserID   string `json:"user_id"`
	FriendID string `json:"friend_id"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserFCMTokenInput struct {
	Token string `json:"token"`
}

type UpdateUserGoalInput struct {
	Days          int `json:"days"`
	DailyDistance int `json:"distance"`
}
