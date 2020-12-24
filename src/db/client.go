package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client contains mongo.Client that is used during Connect
type Client struct {
	Client *mongo.Client
}

// Connect creates a new Client and initializes it for a given context.
func Connect(ctx context.Context, uri string) (*Client, error) {
	var err error
	var client *mongo.Client
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return &Client{Client: client}, err
}
