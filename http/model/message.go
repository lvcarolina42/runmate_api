package model

import (
	"time"

	"runmate_api/internal/entity"
)

type Message struct {
	User    *User     `json:"user"`
	Content string    `json:"message"`
	Date    time.Time `json:"date"`
}

func NewMessageFromEntity(message *entity.Message, user *entity.User) *Message {
	return &Message{
		User:    NewUserFromEntity(user),
		Content: message.Content,
		Date:    message.CreatedAt,
	}
}
