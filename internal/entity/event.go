package entity

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title     string
	Date      time.Time
	CreatedBy uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Users     []*User `gorm:"many2many:user_events;constraint:OnDelete:CASCADE"`
}
