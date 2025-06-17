package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"runmate_api/http/model"
	"runmate_api/internal/chat"
	"runmate_api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type chatHandler struct {
	activityService  *service.Activity
	challengeService *service.Challenge
	messageService   *service.Message
	userService      *service.User

	hub      *chat.Hub
	consumer *chat.Consumer
	upgrader websocket.Upgrader
}

func NewChat(
	activityService *service.Activity,
	challengeService *service.Challenge,
	messageService *service.Message,
	userService *service.User,
	hub *chat.Hub,
	consumer *chat.Consumer,
) *chatHandler {
	return &chatHandler{
		activityService:  activityService,
		challengeService: challengeService,
		messageService:   messageService,
		userService:      userService,

		hub:      hub,
		consumer: consumer,
		upgrader: websocket.Upgrader{},
	}
}

func (c *chatHandler) Routes(r *chi.Mux) {
	r.Route("/chat", func(r chi.Router) {
		r.Get("/{id}", c.handle)
		r.Get("/{id}/messages", c.getMessages)
	})
}

func (c *chatHandler) handle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	publisher := chat.NewPublisher(id)
	defer publisher.Close()

	c.hub.AddConnection(id, conn, func() {
		ctx, cancel := context.WithCancel(context.Background())
		c.hub.Consumers[id] = cancel
		publisher.Start()
		go c.consumer.Start(ctx, id)
	})

	defer func() {
		c.hub.RemoveConnection(id, conn)
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		publisher.Publish(msg)
	}
}

func (c *chatHandler) getMessages(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	messages, err := c.messageService.ListByChallengeID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]*model.Message, 0, len(messages))
	for _, message := range messages {
		user, err := c.userService.GetByID(r.Context(), message.UserID.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result = append(result, model.NewMessageFromEntity(message, user))
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
