package service

import (
	"context"
	"fmt"
	"maps"
	"runmate_api/internal/entity"
	"runmate_api/internal/firebase"
	"runmate_api/internal/repository"
	"slices"
)

func newChallengeMessageNotification(userName, challengeTitle, message string) *firebase.Notification {
	return &firebase.Notification{
		Title: fmt.Sprintf("💬 %s ▸ %s", userName, challengeTitle),
		Body:  message,
	}
}

type Message struct {
	challengeRepo *repository.Challenge
	messageRepo   *repository.Message
	userRepo      *repository.User

	firebaseClient *firebase.Client
}

func NewMessage(challengeRepo *repository.Challenge, messageRepo *repository.Message, userRepo *repository.User, firebaseClient *firebase.Client) *Message {
	return &Message{
		challengeRepo: challengeRepo,
		messageRepo:   messageRepo,
		userRepo:      userRepo,

		firebaseClient: firebaseClient,
	}
}

func (m *Message) Create(ctx context.Context, message *entity.Message, sender *entity.User) error {
	err := m.messageRepo.Save(ctx, message)
	if err != nil {
		return err
	}

	challenge, err := m.challengeRepo.GetByID(ctx, message.ChallengeID.String())
	if err != nil {
		return err
	}

	if challenge == nil {
		return ErrChallengeNotFound
	}

	tokens := make(map[string]any, len(challenge.Users)-1)
	for _, user := range challenge.Users {
		if user.ID == message.UserID || user.FCMToken == "" {
			continue
		}

		tokens[user.FCMToken] = struct{}{}
	}

	notification := newChallengeMessageNotification(sender.Name, challenge.Title, message.Content)
	err = m.firebaseClient.SendNotification(ctx, notification, slices.Collect(maps.Keys(tokens)))
	if err != nil {
		return err
	}

	return nil
}

func (m *Message) ListByChallengeID(ctx context.Context, challengeID string) ([]*entity.Message, error) {
	return m.messageRepo.GetAllByChallengeID(ctx, challengeID)
}
