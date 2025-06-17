package service

import (
	"context"
	"errors"
	"fmt"

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

func (c *Challenge) ListAllActive(ctx context.Context) ([]*entity.Challenge, error) {
	return c.challengeRepo.GetAllActive(ctx)
}

func (c *Challenge) ListAllActiveWithoutUserID(ctx context.Context, userID string) ([]*entity.Challenge, error) {
	user, err := c.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return c.challengeRepo.GetAllActiveWithoutUser(ctx, user)
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

func (c *Challenge) GetRanking(ctx context.Context, challenge *entity.Challenge) ([]*entity.ChallengeRanking, error) {
	rankingItems, err := c.challengeRepo.GetRanking(ctx, challenge)
	if err != nil {
		return nil, err
	}

	ranking := make([]*entity.ChallengeRanking, 0, len(rankingItems))
	for i, item := range rankingItems {
		fmt.Printf("item: %v\n", item)
		user, err := c.userRepo.GetByID(ctx, item.UserID.String())
		if err != nil {
			return nil, err
		}

		if user == nil {
			return nil, ErrUserNotFound
		}

		ranking = append(ranking, &entity.ChallengeRanking{
			User:     user,
			Position: i,
			Distance: item.Distance,
		})
	}

	return ranking, nil
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
