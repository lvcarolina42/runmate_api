package service

import (
	"context"
	"runmate_api/internal/entity"
	"runmate_api/internal/repository"
)

type Message struct {
	challengeRepo *repository.Challenge
	messageRepo   *repository.Message
	userRepo      *repository.User
}

func NewMessage(challengeRepo *repository.Challenge, messageRepo *repository.Message, userRepo *repository.User) *Message {
	return &Message{challengeRepo: challengeRepo, messageRepo: messageRepo, userRepo: userRepo}
}

func (m *Message) Create(ctx context.Context, message *entity.Message) error {
	return m.messageRepo.Save(ctx, message)
}

func (m *Message) ListByChallengeID(ctx context.Context, challengeID string) ([]*entity.Message, error) {
	return m.messageRepo.GetAllByChallengeID(ctx, challengeID)
}
