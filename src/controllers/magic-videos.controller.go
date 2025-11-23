package controllers

import (
	"context"
	"log"
	"magic-server-2026/src/db"
	"magic-server-2026/src/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func MagicVideosCollectionInit() *mongo.Collection {
	return db.GetCollection("magic899_db", "magic_videos")
}

// Get All Magic Videos

func GetMagicVideos(c fiber.Ctx) error {
	// Set up a timeout for the context
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	// Initialize magic Videos collection
	magicVideoCollection := MagicVideosCollectionInit()

	var magicVideos []models.MagicVideos

	// Query MongoDB for all magic videos documents
	cursor, err := magicVideoCollection.Find(ctx, bson.M{})
	if err != nil {
		// Return an internal server error and log the issue
		log.Println("Error finding magic videos:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch magic videos",
		})
	}
	// Ensure the cursor is closed properly
	defer cursor.Close(ctx)

	// Decode all documents into the magic videos slice
	if err = cursor.All(ctx, &magicVideos); err != nil {
		log.Println("Error decoding magic videos documents:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode magic videos data",
		})
	}

	// Return the fetched magic videos in a response
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Magic Videos fetched successfully",
		"videos":  magicVideos,
	})
}

// Get Magic Video

func GetMagicVideo(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	eventID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Magic Video ID"})
	}

	// Initialize the collection dynamically
	magicVideoCollection := MagicVideosCollectionInit()

	var magicVideo models.MagicVideos
	err = magicVideoCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&magicVideo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Magic Video not found"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Magic Video fetched successfully", "video": magicVideo})
}

// Create Magic Video

func CreateMagicVideo(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var magicVideo models.MagicVideos
	if err := c.Bind().JSON(&magicVideo); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	magicVideo.ID = primitive.NewObjectID()
	magicVideo.Created_at = primitive.NewDateTimeFromTime(time.Now())
	magicVideo.Updated_at = primitive.NewDateTimeFromTime(time.Now())

	// Initialize the collection dynamically
	magicVideoCollection := MagicVideosCollectionInit()

	_, err := magicVideoCollection.InsertOne(ctx, magicVideo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create magic video"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Magic Video created successfully",
		"video":   magicVideo,
	})
}

// Update Magic Video

func UpdateMagicVideo(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	magicVideoID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(magicVideoID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Magic Video ID"})
	}

	var updateData map[string]interface{}
	if err := c.Bind().JSON(&updateData); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Remove fields that should not be updated
	delete(updateData, "created_at")

	// Ensure `updated_at` is stored as MongoDB DateTime
	updateData["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	magicVideoCollection := MagicVideosCollectionInit()

	// Update the concert document
	updateResult, err := magicVideoCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update magic video"})
	}

	if updateResult.MatchedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Magic Video not found"})
	}

	// Retrieve the updated document
	var updatedMagicVideo models.MagicVideos
	err = magicVideoCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedMagicVideo)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated Magic Video"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Magic Video updated successfully",
		"video":   updatedMagicVideo,
	})
}

// Delete Magic Video

func DeleteMagicVideo(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	magicVideoID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(magicVideoID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Magic Video ID"})
	}

	// Initialize the collection dynamically
	magicVideoCollection := MagicVideosCollectionInit()

	// Delete the event
	_, err = magicVideoCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete magic video"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Magic Video deleted successfully"})
}
