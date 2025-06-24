package chat

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"runmate_api/config"
	"runmate_api/http/model"
	"runmate_api/internal/entity"
	"runmate_api/internal/service"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type messagePayload struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
	Type    int    `json:"type"`
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
		Type:        p.Type,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}, nil
}

func getTopic(challengeID string) string {
	return fmt.Sprintf("chat-challenge-%s", challengeID)
}

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(challengeID string) *Publisher {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.KafkaURL()),
		Topic:                  getTopic(challengeID),
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
	}
	if config.Production() {
		writer.Transport = &kafka.Transport{
			DialTimeout: 10 * time.Second,
			TLS:         &tls.Config{},
			SASL: plain.Mechanism{
				Username: config.KafkaUsername(),
				Password: config.KafkaPassword(),
			},
		}
	}

	return &Publisher{
		writer: writer,
	}
}

func (p *Publisher) Start() error {
	fmt.Println("STARTING")

	message, err := json.Marshal(messagePayload{
		UserID:  "00000000-0000-0000-0000-000000000000",
		Content: "Publisher started",
		Type:    entity.MessageTypeSystem,
	})
	if err != nil {
		return err
	}

	retries := 5
	for i := 0; i < retries; i++ {
		err = p.writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(uuid.New().String()),
				Value: message,
			},
		)
		if err != nil {
			log.Println("Error strating publisher:", err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		fmt.Println("Publisher started")
		break
	}

	return err
}

func (p *Publisher) Close() {
	p.writer.Close()
}

func (p *Publisher) Publish(message []byte) error {
	err := p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(uuid.New().String()),
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
	userService    *service.User
}

func NewConsumer(hub *Hub, messageService *service.Message, userService *service.User) *Consumer {
	return &Consumer{
		hub:            hub,
		messageService: messageService,
		userService:    userService,
	}
}

func (c *Consumer) Start(ctx context.Context, challengeID string) {
	topic := getTopic(challengeID)

	readerConfig := kafka.ReaderConfig{
		Brokers: []string{config.KafkaURL()},
		Topic:   topic,
		GroupID: "chat-group-" + challengeID,
	}
	if config.Production() {
		readerConfig.Dialer = &kafka.Dialer{
			Timeout: 10 * time.Second,
			TLS:     &tls.Config{},
			SASLMechanism: plain.Mechanism{
				Username: config.KafkaUsername(),
				Password: config.KafkaPassword(),
			},
		}
	}

	reader := kafka.NewReader(readerConfig)
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

			msg, err := payload.ToEntity(challengeID)
			if err != nil {
				log.Println("Failed to create message entity:", err)
				continue
			}

			if msg.Type == entity.MessageTypeUser {
				user, err := c.userService.GetByID(ctx, msg.UserID.String())
				if err != nil {
					log.Println("Failed to get user:", err)
					continue
				}

				if err := c.messageService.Create(ctx, msg, user); err != nil {
					log.Println("Failed to save message:", err)
				}

				messageData, err := json.Marshal(model.NewMessageFromEntity(msg, user))
				if err != nil {
					log.Println("Failed to marshal message:", err)
					continue
				}

				c.hub.Broadcast(challengeID, messageData)
			}

			if err := reader.CommitMessages(ctx, m); err != nil {
				log.Println("Failed to commit message:", err)
			}
		}
	}()
}
