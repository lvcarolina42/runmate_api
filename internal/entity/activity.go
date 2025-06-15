package entity

import (
	"time"

	"github.com/google/uuid"
)

type Activity struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID
	Title       string
	Date        time.Time
	Duration    int
	Distance    int
	Coordinates []*Coordinate `gorm:"foreignKey:ActivityID;constraint:OnDelete:CASCADE"`
}

type Coordinate struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ActivityID uuid.UUID // This field is used as a foreign key.
	Lat        float64
	Long       float64
	Order      int
}
