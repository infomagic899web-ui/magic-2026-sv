package helpers

import (
	"context"
	"magic899-server/src/models"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func isValidString(str string) bool {
	// Allow letters, numbers, whitespace characters, hyphens, apostrophes
	re := regexp.MustCompile(`^[a-zA-Z0-9\s\-']+$`)
	trimmed := strings.TrimSpace(str)
	return re.MatchString(trimmed) && trimmed != ""
}

func IncreaseRequestedCount(collection *mongo.Collection, song models.Track, requester string) error {
	// Fetch the existing song by ID
	var existing models.Track
	if err := collection.FindOne(context.Background(), bson.M{"_id": song.ID}).Decode(&existing); err != nil {
		return err
	}

	// Update request count
	existing.RequestCount++

	// Use map to ensure uniqueness
	requestedBySet := make(map[string]struct{})
	for _, r := range existing.RequestedBy {
		requestedBySet[strings.TrimSpace(r)] = struct{}{}
	}

	// Add new requester
	newRequester := strings.TrimSpace(requester)
	if isValidString(newRequester) {
		requestedBySet[newRequester] = struct{}{}
	}

	// Convert back to slice
	updatedRequesters := make([]string, 0, len(requestedBySet))
	for r := range requestedBySet {
		updatedRequesters = append(updatedRequesters, r)
	}

	// Update DB
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": existing.ID},
		bson.M{
			"$set": bson.M{
				"requestCount": existing.RequestCount,
				"requestedBy":  updatedRequesters,
				"updated_at":   primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)

	return err
}
