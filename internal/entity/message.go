package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	MessageTypeUser   = 0
	MessageTypeSystem = 1
)

type Message struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Content     string
	ChallengeID uuid.UUID
	UserID      uuid.UUID
	Type        int
	CreatedAt   time.Time
}
