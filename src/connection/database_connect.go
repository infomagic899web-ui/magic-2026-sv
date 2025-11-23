package connection

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectToDatabase(uri string, connectionType string) (*mongo.Client, error) {
	log.Printf("Attempting to connect to the %s database...", connectionType)

	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Printf("Successfully connected to the %s database!", connectionType)
	return client, nil
}
