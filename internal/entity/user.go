package entity

import (
	"errors"
	"math"
	"regexp"
	"time"

	"github.com/google/uuid"
)

const (
	UserRoleUser  = 0
	UserRoleAdmin = 1
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`)
)

type User struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username          string    `gorm:"unique"`
	Email             string    `gorm:"unique"`
	Password          string
	Name              string
	FCMToken          string
	Role              int8
	XP                int
	GoalDays          *int
	GoalDailyDistance *int
	Birthdate         time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Activities        []*Activity                 `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Friends           []*User                     `gorm:"many2many:user_friends;constraint:OnDelete:CASCADE"`
	Challenges        []*Challenge                `gorm:"many2many:user_challenges;constraint:OnDelete:CASCADE"`
	CreatedChallenges []*Challenge                `gorm:"foreignKey:CreatedBy;constraint:OnDelete:CASCADE"`
	ChallengeEvents   []*ChallengeEvent           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Events            []*Event                    `gorm:"many2many:user_events;constraint:OnDelete:CASCADE"`
	WeekActivities    map[string]*UserDayActitivy `gorm:"-:all"`
}

func (u *User) Validate() error {
	if !usernameRegex.MatchString(u.Username) {
		return errors.New("invalid username")
	}

	return nil
}

func (u *User) CurrentLevel() int {
	return int(math.Sqrt(float64(1000*(2*u.XP+250)))+500) / 1000
}

func (u *User) levelXP(level int) int {
	return int((math.Pow(float64(level), 2)+float64(level))/2*1000) - (level * 1000)
}

func (u *User) NextLevelXP() int {
	return u.levelXP(u.CurrentLevel()+1) - u.XP
}
