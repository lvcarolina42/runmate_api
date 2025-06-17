package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Content     string
	ChallengeID uuid.UUID
	UserID      uuid.UUID
	CreatedAt   time.Time
}
