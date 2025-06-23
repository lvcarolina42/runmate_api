package repository

import (
	"context"
	"fmt"

	"runmate_api/internal/entity"

	"gorm.io/gorm"
)

type Challenge struct {
	db *gorm.DB
}

func NewChallenge(db *gorm.DB) *Challenge {
	return &Challenge{db: db}
}

func (c *Challenge) Create(ctx context.Context, challenge *entity.Challenge) error {
	result := c.db.WithContext(ctx).Create(challenge)
	if result.Error != nil {
		return fmt.Errorf("failed to create challenge: %v", result.Error)
	}

	return nil
}

func (c *Challenge) GetByID(ctx context.Context, id string) (*entity.Challenge, error) {
	var challenge entity.Challenge
	result := c.db.WithContext(ctx).Preload("Users").Where("id = ?", id).First(&challenge)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get challenge %s: %v", id, result.Error)
	}

	return &challenge, nil
}

func (c *Challenge) GetAllActive(ctx context.Context) ([]*entity.Challenge, error) {
	var challenges []*entity.Challenge
	result := c.db.WithContext(ctx).Preload("Users").Where("end_date IS NULL OR end_date > NOW()").Find(&challenges)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get active challenges: %v", result.Error)
	}

	return challenges, nil
}

func (c *Challenge) Update(ctx context.Context, challenge *entity.Challenge) error {
	result := c.db.WithContext(ctx).Save(challenge)
	if result.Error != nil {
		return fmt.Errorf("failed to update challenge: %v", result.Error)
	}

	return nil
}

func (c *Challenge) GetAllActiveWithoutUser(ctx context.Context, user *entity.User) ([]*entity.Challenge, error) {
	var challenges []*entity.Challenge
	err := c.db.WithContext(ctx).
		Preload("Users").
		Table("challenges").
		Select("challenges.*").
		Joins("LEFT JOIN user_challenges ON challenges.id = user_challenges.challenge_id AND user_challenges.user_id = ?", user.ID).
		Where("user_challenges.challenge_id IS NULL AND (challenges.end_date IS NULL OR challenges.end_date > NOW())").
		Find(&challenges).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to get challenges: %v", err)
	}

	return challenges, nil
}

func (c *Challenge) GetAllActiveByUser(ctx context.Context, user *entity.User) ([]*entity.Challenge, error) {
	var challenges []*entity.Challenge
	err := c.db.WithContext(ctx).Model(&user).Where("end_date IS NULL OR end_date > NOW()").Preload("Users").Association("Challenges").Find(&challenges)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenges: %v", err)
	}

	return challenges, nil
}

func (c *Challenge) GetAllByUser(ctx context.Context, user *entity.User) ([]*entity.Challenge, error) {
	var challenges []*entity.Challenge
	err := c.db.WithContext(ctx).Model(&user).Preload("Users").Association("Challenges").Find(&challenges)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenges: %v", err)
	}

	return challenges, nil
}

func (c *Challenge) AddEvent(ctx context.Context, challenge *entity.Challenge, event *entity.ChallengeEvent) error {
	err := c.db.WithContext(ctx).Model(&challenge).Association("Events").Append(event)
	if err != nil {
		return fmt.Errorf("failed to add event to challenge: %v", err)
	}

	return nil
}

func (c *Challenge) GetAllEventsByUser(ctx context.Context, challenge *entity.Challenge, user *entity.User) ([]*entity.ChallengeEvent, error) {
	var events []*entity.ChallengeEvent
	err := c.db.WithContext(ctx).Model(&challenge).Association("Events").Find(&events)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %v", err)
	}

	return events, nil
}

func (c *Challenge) AddUser(ctx context.Context, challenge *entity.Challenge, user *entity.User) error {
	err := c.db.WithContext(ctx).Model(&challenge).Association("Users").Append(user)
	if err != nil {
		return fmt.Errorf("failed to add user to challenge: %v", err)
	}

	return nil
}

func (c *Challenge) GetRanking(ctx context.Context, challenge *entity.Challenge) ([]*entity.ChallengeRankingResult, error) {
	var results []*entity.ChallengeRankingResult
	err := c.db.
		WithContext(ctx).
		Table("challenge_events").
		Select("user_id, SUM(distance) AS distance").
		Where("challenge_id = ?", challenge.ID).
		Group("user_id").
		Group("challenge_id").
		Order("distance DESC").
		Scan(&results).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to get rankings: %v", err)
	}

	return results, nil
}
