package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username  string
	Email     string
	Password  string
	Name      string
	Role      int8
	Birthdate time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
