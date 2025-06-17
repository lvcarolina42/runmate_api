package repository

import (
	"context"

	"runmate_api/internal/entity"

	"gorm.io/gorm"
)

type Message struct {
	db *gorm.DB
}

func NewMessage(db *gorm.DB) *Message {
	return &Message{db}
}

func (r *Message) Save(ctx context.Context, message *entity.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *Message) GetAllByChallengeID(ctx context.Context, challengeID string) ([]*entity.Message, error) {
	var messages []*entity.Message
	err := r.db.WithContext(ctx).Where("challenge_id = ?", challengeID).Order("created_at ASC").Find(&messages).Error
	return messages, err
}
