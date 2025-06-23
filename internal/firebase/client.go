package firebase

import (
	"context"
	"fmt"
	"runmate_api/config"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type Client struct {
	app *firebase.App
}

func NewClient() (*Client, error) {
	var options []option.ClientOption
	if config.Production() {
		options = append(options, option.WithCredentialsJSON(config.FirebaseCredentials()))
	}

	app, err := firebase.NewApp(context.Background(), nil, options...)
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
			return fmt.Errorf("error sending notification: %v", err)
		}
	}

	return nil
}
