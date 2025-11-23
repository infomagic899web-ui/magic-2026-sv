package connection

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func TryConnect(uri string, connectionType string) *mongo.Client {
	if uri == "" {
		log.Printf("No URI provided for %s database.", connectionType)
		return nil
	}

	client, err := connectToDatabase(uri, connectionType)
	if err != nil {
		log.Printf("Error connecting to %s database: %v", connectionType, err)
		return nil
	}
	return client
}
