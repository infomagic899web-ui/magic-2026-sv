package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Nonce struct {
	Value     string    `bson:"value"`
	UserID    string    `bson:"user_id"`
	ExpiresAt time.Time `bson:"expires_at"`
}

func StoreNonce(coll *mongo.Collection, userID, nonce string, ttl time.Duration) error {
	_, err := coll.InsertOne(context.TODO(), Nonce{
		Value:     nonce,
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
	})
	return err
}

func ValidateNonce(coll *mongo.Collection, userID, nonce string) (bool, error) {
	res := coll.FindOne(context.TODO(), bson.M{
		"value":      nonce,
		"user_id":    userID,
		"expires_at": bson.M{"$gt": time.Now()},
	})
	if res.Err() != nil {
		return false, res.Err()
	}
	return true, nil
}

func InvalidateNonce(coll *mongo.Collection, userID, nonce string) error {
	_, err := coll.DeleteOne(context.TODO(), bson.M{
		"value":   nonce,
		"user_id": userID,
	})
	return err
}
