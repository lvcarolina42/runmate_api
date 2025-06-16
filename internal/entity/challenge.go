package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChallengeType int8

const (
	ChallengeTypeDistance ChallengeType = 0
	ChallengeTypeDate     ChallengeType = 1
)

type Challenge struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title         string
	Description   string
	StartDate     time.Time
	EndDate       *time.Time
	Type          ChallengeType
	TotalDistance *int
	CreatedBy     uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Users         []*User           `gorm:"many2many:user_challenges;constraint:OnDelete:CASCADE"`
	Events        []*ChallengeEvent `gorm:"foreignKey:ChallengeID;constraint:OnDelete:CASCADE"`
}

type ChallengeEvent struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ChallengeID uuid.UUID `gorm:"type:uuid;not null"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Distance    int
	Date        time.Time
}
