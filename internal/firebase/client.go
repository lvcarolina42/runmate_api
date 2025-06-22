package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type Client struct {
	app *firebase.App
}

func NewClient() (*Client, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return &Client{
		app: app,
	}, nil
}

func (c *Client) SendNotification(ctx context.Context, notification *Notification, tokens []string) error {
	messagingClient, err := c.app.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("error getting messaging client: %v", err)
	}

	for _, token := range tokens {
		_, err = messagingClient.Send(ctx, &messaging.Message{
			Notification: notification,
			Token:        token,
		})
		if err != nil {
			return fmt.Errorf("error sending message: %v", err)
		}
	}

	return nil
}
