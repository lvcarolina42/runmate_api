package entity

import (
	"errors"
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
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username  string    `gorm:"unique"`
	Email     string    `gorm:"unique"`
	Password  string
	Name      string
	Role      int8
	Birthdate time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Friends   []*User `gorm:"many2many:user_friends;constraint:OnDelete:CASCADE"`
}

func (u *User) Validate() error {
	if !usernameRegex.MatchString(u.Username) {
		return errors.New("invalid username")
	}

	return nil
}
