package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type Client struct {
	Auth *auth.Client
}

func NewClient(ctx context.Context, credFile string) *Client {
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credFile))
	if err != nil {
		log.Fatal("error initializing Firebase app:", err)
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatal("error initializing Firebase auth client:", err)
	}
	return &Client{
		Auth: authClient,
	}
}
