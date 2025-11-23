package db

import (
	"context"
	"log"
	"magic-server-2026/src/connection"
	"magic-server-2026/src/utils"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var Client *mongo.Client // ✅ Global client

func Init() {
	onlineURI, offlineURI := utils.GetVariables()

	// ✅ Try connecting to the online database first
	log.Println("Attempting to connect to the ONLINE database...")
	Client = connection.TryConnect(onlineURI, "online")

	// Optional: validate if the client is actually alive
	if !isConnected(Client) {
		log.Println("Online connection failed, attempting OFFLINE database...")
		Client = connection.TryConnect(offlineURI, "offline")
	}

	if !isConnected(Client) {
		log.Fatal("❌ Failed to connect to both online and offline databases. Shutting down...")
	}

	log.Println("✅ Database connection successfully initialized!")
}

// Helper: Verify if client is connected
func isConnected(client *mongo.Client) bool {
	if client == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return false
	}
	return true
}

func GetCollection(dbName, collName string) *mongo.Collection {
	return Client.Database(dbName).Collection(collName)
}
