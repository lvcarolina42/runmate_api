package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"runmate_api/config"
	"runmate_api/internal/entity"
	"runmate_api/internal/service"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type messagePayload struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

func (p messagePayload) ToEntity(challengeID string) (*entity.Message, error) {
	userID, err := uuid.Parse(p.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user id: %v", err)
	}

	cID, err := uuid.Parse(challengeID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse challenge id: %v", err)
	}

	return &entity.Message{
		Content:     p.Content,
		ChallengeID: cID,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}, nil
}

func getTopic(challengeID string) string {
	return fmt.Sprintf("chat-challenge-%s", challengeID)
}

func PublishMessage(challengeID string, message []byte) error {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.KafkaURL()),
		Topic:        getTopic(challengeID),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}
	defer writer.Close()

	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(challengeID),
			Value: message,
		},
	)
	if err != nil {
		log.Println("Error publishing message:", err)
	}
	return err
}

type Consumer struct {
	hub            *Hub
	messageService *service.Message
}

func NewConsumer(hub *Hub, messageService *service.Message) *Consumer {
	return &Consumer{hub: hub, messageService: messageService}
}

func (c *Consumer) Start(ctx context.Context, challengeID string) {
	topic := getTopic(challengeID)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.KafkaURL()},
		Topic:   topic,
		GroupID: "chat-group-" + challengeID,
	})
	go func() {
		defer reader.Close()

		for {
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Println("Kafka consumer closed for", topic)
				return
			}

			var payload messagePayload
			err = json.Unmarshal(m.Value, &payload)
			if err != nil {
				log.Println("Failed to unmarshal message:", err)
				continue
			}

			message, err := payload.ToEntity(challengeID)
			if err != nil {
				log.Println("Failed to create message entity:", err)
				continue
			}

			if err := c.messageService.Create(ctx, message); err != nil {
				log.Println("Failed to save message:", err)
			}

			c.hub.Broadcast(challengeID, m.Value)
		}
	}()
}
