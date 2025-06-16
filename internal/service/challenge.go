package service

import (
	"context"
	"errors"

	"runmate_api/internal/entity"
	"runmate_api/internal/repository"
)

var (
	ErrChallengeNotFound = errors.New("challenge not found")
)

type Challenge struct {
	challengeRepo *repository.Challenge
	userRepo      *repository.User
}

func NewChallenge(challengeRepo *repository.Challenge, userRepo *repository.User) *Challenge {
	return &Challenge{challengeRepo: challengeRepo, userRepo: userRepo}
}

func (c *Challenge) Create(ctx context.Context, challenge *entity.Challenge) error {
	user, err := c.userRepo.GetByID(ctx, challenge.CreatedBy.String())
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	challenge.Users = []*entity.User{user}
	return c.challengeRepo.Create(ctx, challenge)
}

func (c *Challenge) GetByID(ctx context.Context, id string) (*entity.Challenge, error) {
	return c.challengeRepo.GetByID(ctx, id)
}

func (c *Challenge) ListAllActiveByUserID(ctx context.Context, userID string) ([]*entity.Challenge, error) {
	user, err := c.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return c.challengeRepo.GetAllActiveByUser(ctx, user)
}

func (c *Challenge) ListAllByUserID(ctx context.Context, userID string) ([]*entity.Challenge, error) {
	user, err := c.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return c.challengeRepo.GetAllByUser(ctx, user)
}

func (c *Challenge) AddEvent(ctx context.Context, challengeID string, event *entity.ChallengeEvent) error {
	challenge, err := c.challengeRepo.GetByID(ctx, challengeID)
	if err != nil {
		return err
	}

	if challenge == nil {
		return ErrChallengeNotFound
	}

	event.ChallengeID = challenge.ID
	return c.challengeRepo.AddEvent(ctx, challenge, event)
}

func (c *Challenge) Join(ctx context.Context, challengeID, userID string) error {
	user, err := c.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	challenge, err := c.challengeRepo.GetByID(ctx, challengeID)
	if err != nil {
		return err
	}

	if challenge == nil {
		return ErrChallengeNotFound
	}

	return c.challengeRepo.AddUser(ctx, challenge, user)
}
