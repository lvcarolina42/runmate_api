package model

import (
	"errors"
	"fmt"
	"runmate_api/internal/entity"
	"time"

	"github.com/google/uuid"
)

type ChallengeType string

func NewChallengeTypeFromEntity(c entity.ChallengeType) ChallengeType {
	switch c {
	case entity.ChallengeTypeDistance:
		return ChallengeTypeDistance
	case entity.ChallengeTypeDate:
		return ChallengeTypeDate
	default:
		return ChallengeTypeDistance
	}
}

func (c ChallengeType) ToEntity() entity.ChallengeType {
	switch c {
	case ChallengeTypeDistance:
		return entity.ChallengeTypeDistance
	case ChallengeTypeDate:
		return entity.ChallengeTypeDate
	default:
		return entity.ChallengeTypeDistance
	}
}

const (
	ChallengeTypeDistance ChallengeType = "distance"
	ChallengeTypeDate     ChallengeType = "date"
)

var (
	ErrStartDateRequired        = errors.New("start date is required")
	ErrEndDateNotRequired       = errors.New("end date is not required")
	ErrTotalDistanceRequired    = errors.New("total distance is required")
	ErrInvalidChallengeType     = errors.New("invalid challenge type")
	ErrEndDateRequired          = errors.New("end date is required")
	ErrTotalDistanceNotRequired = errors.New("total distance is not required")
	ErrEndDateBeforeStartDate   = errors.New("end date must be after start date")
)

type Challenge struct {
	ID            string        `json:"id"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	StartDate     time.Time     `json:"start_date"`
	EndDate       *time.Time    `json:"end_date,omitempty"`
	TotalDistance *int          `json:"total_distance,omitempty"`
	Type          ChallengeType `json:"type"`
	Finished      bool          `json:"finished"`
}

func NewChallengeFromEntity(c *entity.Challenge) *Challenge {
	var finished bool
	if c.Type == entity.ChallengeTypeDistance {
		finished = c.EndDate != nil
	} else if c.Type == entity.ChallengeTypeDate {
		finished = c.EndDate != nil && c.EndDate.Before(time.Now())
	}

	return &Challenge{
		ID:            c.ID.String(),
		Title:         c.Title,
		Description:   c.Description,
		StartDate:     c.StartDate,
		EndDate:       c.EndDate,
		TotalDistance: c.TotalDistance,
		Type:          NewChallengeTypeFromEntity(c.Type),
		Finished:      finished,
	}
}

type CreateChallengeInput struct {
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	StartDate     time.Time     `json:"start_date"`
	EndDate       *time.Time    `json:"end_date,omitempty"`
	TotalDistance *int          `json:"total_distance,omitempty"`
	Type          ChallengeType `json:"type"`
	UserID        string        `json:"created_by"`
}

func (c *CreateChallengeInput) Validate() error {
	if c.StartDate.IsZero() {
		return ErrStartDateRequired
	}

	if c.Type != ChallengeTypeDistance && c.Type != ChallengeTypeDate {
		return ErrInvalidChallengeType
	}

	if c.Type == ChallengeTypeDistance {
		if c.TotalDistance == nil || *c.TotalDistance <= 0 {
			return ErrTotalDistanceRequired
		}

		if c.EndDate != nil {
			return ErrEndDateNotRequired
		}
	}

	if c.Type == ChallengeTypeDate {
		if c.EndDate == nil || c.EndDate.IsZero() {
			return ErrEndDateRequired
		}

		if c.TotalDistance != nil {
			return ErrTotalDistanceNotRequired
		}

		if c.EndDate.Before(c.StartDate) {
			return ErrEndDateBeforeStartDate
		}
	}

	return nil
}

func (c *CreateChallengeInput) ToEntity() (*entity.Challenge, error) {
	userID, err := uuid.Parse(c.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user id: %v", err)
	}

	return &entity.Challenge{
		Title:         c.Title,
		Description:   c.Description,
		StartDate:     c.StartDate,
		EndDate:       c.EndDate,
		TotalDistance: c.TotalDistance,
		Type:          c.Type.ToEntity(),
		CreatedBy:     userID,
	}, nil
}

type JoinChallengeInput struct {
	UserID      string `json:"user_id"`
	ChallengeID string `json:"challenge_id"`
}
