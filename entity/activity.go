package entity

import (
	"github.com/google/uuid"
	"time"
)

type Coordinate struct {
	ID         uuid.UUID `json:"-" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ActivityID uuid.UUID `json:"-"` // This field is used as a foreign key.
	Lat        float64   `json:"lat"`
	Long       float64   `json:"long"`
	Order      int       `json:"-"`
}

type Activity struct {
	ID          uuid.UUID     `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID     `json:"user_id"`
	Title       string        `json:"title"`
	Date        time.Time     `json:"date"`
	Duration    int           `json:"duration"`
	Distance    int           `json:"distance"`
	Coordinates []*Coordinate `json:"coordinates" gorm:"foreignKey:ActivityID;constraint:OnDelete:CASCADE"`
}
