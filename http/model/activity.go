package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"runmate_api/internal/entity"
)

type Activity struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	Title       string        `json:"title"`
	Date        time.Time     `json:"date"`
	Duration    int           `json:"duration"`
	Distance    int           `json:"distance"`
	Coordinates []*Coordinate `json:"coordinates"`
}

func NewActivityFromEntity(activity *entity.Activity) *Activity {
	return &Activity{
		ID:          activity.ID.String(),
		UserID:      activity.UserID.String(),
		Title:       activity.Title,
		Date:        activity.Date,
		Duration:    activity.Duration,
		Distance:    activity.Distance,
		Coordinates: newCoordinatesFromEntity(activity.Coordinates),
	}
}

func (a *Activity) ToEntity() (*entity.Activity, error) {
	id, err := uuid.Parse(a.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse id: %v", err)
	}

	userID, err := uuid.Parse(a.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user id: %v", err)
	}

	coordinates := make([]*entity.Coordinate, 0, len(a.Coordinates))
	for i, coordinate := range a.Coordinates {
		coordinates = append(coordinates, coordinate.ToEntity(i))
	}

	return &entity.Activity{
		ID:          id,
		UserID:      userID,
		Title:       a.Title,
		Date:        a.Date,
		Duration:    a.Duration,
		Distance:    a.Distance,
		Coordinates: coordinates,
	}, nil
}

type Coordinate struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func newCoordinateFromEntity(coordinate *entity.Coordinate) *Coordinate {
	return &Coordinate{
		Lat:  coordinate.Lat,
		Long: coordinate.Long,
	}
}

func newCoordinatesFromEntity(coordinates []*entity.Coordinate) []*Coordinate {
	var result []*Coordinate
	for _, coordinate := range coordinates {
		result = append(result, newCoordinateFromEntity(coordinate))
	}
	return result
}

func (c *Coordinate) ToEntity(order int) *entity.Coordinate {
	return &entity.Coordinate{
		Lat:   c.Lat,
		Long:  c.Long,
		Order: order,
	}
}

type CreateActivityCoordinateInput struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func (c *CreateActivityCoordinateInput) ToEntity(order int) *entity.Coordinate {
	return &entity.Coordinate{
		Lat:   c.Lat,
		Long:  c.Long,
		Order: order,
	}
}

type CreateActivityInput struct {
	UserID      string                           `json:"user_id"`
	Title       string                           `json:"title"`
	Date        time.Time                        `json:"date"`
	Duration    int                              `json:"duration"`
	Distance    int                              `json:"distance"`
	Coordinates []*CreateActivityCoordinateInput `json:"coordinates"`
}

func (c *CreateActivityInput) ToEntity() (*entity.Activity, error) {
	userID, err := uuid.Parse(c.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user id: %v", err)
	}

	coordinates := make([]*entity.Coordinate, 0, len(c.Coordinates))
	for i, coordinate := range c.Coordinates {
		coordinates = append(coordinates, coordinate.ToEntity(i))
	}

	return &entity.Activity{
		UserID:      userID,
		Title:       c.Title,
		Date:        c.Date,
		Duration:    c.Duration,
		Distance:    c.Distance,
		Coordinates: coordinates,
	}, nil
}
