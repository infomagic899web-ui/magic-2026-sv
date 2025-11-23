package controllers

import (
	"context"
	"encoding/json"
	"log"
	"magic-server-2026/src/db"
	"magic-server-2026/src/models"
	"magic-server-2026/src/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Music Collection Initialization
func MusicCollectionInit() *mongo.Collection {
	return db.GetCollection("magic899_db", "music")
}

func VoteCollectionInit() *mongo.Collection {
	return db.GetCollection("magic899_db", "votes")
}

// IncrementVote increments vote count for a track
func IncrementVote(c fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID parameter is required"})
	}

	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	// Check if it's Mon–Fri (weekday)
	isWeekday := utils.IsVotingOpen()

	var update bson.M

	if isWeekday {
		// Weekdays → ONLY increment upcoming_votes
		update = bson.M{
			"$inc": bson.M{"upcoming_votes": 1},
			"$set": bson.M{"updated_at": time.Now()},
		}
	} else {
		// Weekends → ONLY increment votes
		update = bson.M{
			"$inc": bson.M{"votes": 1},
			"$set": bson.M{"updated_at": time.Now()},
		}
	}

	result, err := MusicCollectionInit().UpdateOne(c.Context(), bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to increment vote"})
	}
	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Track not found"})
	}

	return c.JSON(fiber.Map{
		"message": "Vote recorded successfully",
	})
}

// Can Vote Check

func CanUserVote(c fiber.Ctx) error {
	musicID := c.Params("id")
	ip := c.IP()

	objectID, err := primitive.ObjectIDFromHex(musicID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid music ID",
		})
	}

	currentDate := time.Now().Format("2006-01-02")

	var record models.VoteRecord
	err = VoteCollectionInit().FindOne(context.TODO(), bson.M{
		"music_id":   objectID,
		"ip_address": ip,
		"date":       currentDate,
	}).Decode(&record)

	if err == mongo.ErrNoDocuments {
		// No record yet, so can vote
		return c.JSON(fiber.Map{
			"canVote":    true,
			"votesToday": 0,
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check vote record",
		})
	}

	// Return can_vote from record
	return c.JSON(fiber.Map{
		"canVote":    record.CanVote,
		"votesToday": record.Votes,
	})
}

// Get All Music
func GetAllMusic(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	musicCollection := MusicCollectionInit()
	var music []models.Music

	cursor, err := musicCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error fetching music:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch music"})
	}

	if err = cursor.All(ctx, &music); err != nil {
		log.Println("Error decoding music:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process music"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "All music fetched successfully", "music": music})
}

// Get a single music by ID
func GetMusic(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	musicID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(musicID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Music ID"})
	}

	musicCollection := MusicCollectionInit()
	var music models.Music

	err = musicCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&music)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Music not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch music"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Music fetched successfully", "music": music})
}

// Create a new music entry with Spotify Image validation
func CreateMusic(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var music models.Music
	if err := json.Unmarshal(c.Body(), &music); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate Image URL (Must be from Spotify: i.scdn.co)
	if !strings.HasPrefix(music.Music_image, "https://i.scdn.co/") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid image URL. Only Spotify images are allowed."})
	}

	music.ID = primitive.NewObjectID()
	music.Created_at = primitive.NewDateTimeFromTime(time.Now())
	music.Updated_at = primitive.NewDateTimeFromTime(time.Now())

	musicCollection := MusicCollectionInit()

	_, err := musicCollection.InsertOne(ctx, music)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create music"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Music created successfully",
	})
}

// Update Music (Prevents Updating `created_at`)
func UpdateMusic(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	musicID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(musicID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Music ID"})
	}

	var updateData map[string]interface{}
	if err := json.Unmarshal(c.Body(), &updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Remove fields that should not be updated
	delete(updateData, "created_at")

	// Validate updated image URL if provided
	if img, exists := updateData["music_image"].(string); exists {
		if !strings.HasPrefix(img, "https://i.scdn.co/") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid image URL. Only Spotify images are allowed."})
		}
	}

	updateData["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	musicCollection := MusicCollectionInit()
	updateResult, err := musicCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil || updateResult.MatchedCount == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update music or music not found"})
	}

	// Retrieve updated music
	var updatedMusic models.Music
	err = musicCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedMusic)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated music"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Music updated successfully", "music": updatedMusic})
}

// Delete Music
func DeleteMusic(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	musicID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(musicID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Music ID"})
	}

	musicCollection := MusicCollectionInit()
	_, err = musicCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete music"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Music deleted successfully"})
}

// Vote for a Song (Ensures only 1 vote per user)
func VoteMusic(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	musicID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(musicID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Music ID"})
	}

	userID := c.Locals("email") // Ensure this is set by authentication middleware
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	musicCollection := MusicCollectionInit()

	// Check if user already voted
	var music models.Music
	err = musicCollection.FindOne(ctx, bson.M{"_id": objID, "voters": userID}).Decode(&music)
	if err == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You can only vote once per song"})
	}

	// Add vote and prevent duplicates
	update := bson.M{
		"$inc":      bson.M{"votes": 1},
		"$addToSet": bson.M{"voters": userID},
	}

	_, err = musicCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to cast vote"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Vote cast successfully"})
}

func CheckEligibility(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	voteCollection := VoteCollectionInit()
	today := time.Now().Format("2006-01-02")
	userIP := c.IP()

	// Sum all votes for today's date from this IP
	cursor, err := voteCollection.Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"date": today, "ip_address": userIP}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$votes"},
		}}},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check eligibility"})
	}
	defer cursor.Close(ctx)

	var result struct {
		Total int `bson:"total"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse vote data"})
		}
	}

	// Default to 0 if no records
	if result.Total >= 3 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"can_vote": false, "remaining": 0})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"can_vote":  true,
		"remaining": 3 - result.Total,
	})
}
