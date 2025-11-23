package controllers

import (
	"context"
	"log"
	"magic-server-2026/src/db"
	"magic-server-2026/src/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ShowsCollectionInit() *mongo.Collection {
	return db.GetCollection("magic899_db", "shows")
}

// Get all Shows
func GetShows(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	showCollection := ShowsCollectionInit()
	cursor, err := showCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Unable to fetch shows"})
	}
	defer cursor.Close(ctx)

	var shows []models.Shows
	for cursor.Next(ctx) {
		var show models.Shows
		if err := cursor.Decode(&show); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error decoding shows"})
		}
		shows = append(shows, show)
	}

	if err := cursor.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error reading cursor"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Shows fetched successfully",
		"shows":   shows,
	})
}

// Get a Show
func GetShow(c fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var show models.Shows
	showCollection := ShowsCollectionInit()

	err = showCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&show)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Show not found"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Show fetched successfully",
		"show":    show,
	})
}

func GetMagicVideosByShowID(c fiber.Ctx) error {
	// Get the showID parameter
	showID := c.Params("id")
	if showID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Show ID is required"})
	}

	// Convert string ID to MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(showID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid show ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the show by ID
	var show models.Shows
	showCollection := ShowsCollectionInit()
	err = showCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&show)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Show not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Error fetching show"})
	}

	// Find related magic videos by show_id
	var magicVideos []models.MagicVideos
	cursor, err := MagicVideosCollectionInit().Find(ctx, bson.M{"show_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error fetching magic videos"})
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &magicVideos); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error decoding magic videos"})
	}

	// Return the result
	return c.Status(200).JSON(fiber.Map{
		"show":         show,
		"magic_videos": magicVideos,
	})
}

func GetByShowName(c fiber.Ctx) error {
	// Get the showName parameter from the URL
	showName := c.Params("showName")
	if showName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Show name is required"})
	}

	// Remove spaces from the show name in the request
	trimmedShowName := strings.ReplaceAll(showName, " ", "") // Remove all spaces from showName

	// MongoDB query to match the show name case-insensitively
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the show in the database, stripping spaces from the stored show_name field as well
	var show models.Shows
	showCollection := ShowsCollectionInit()

	err := showCollection.FindOne(ctx, bson.M{
		"$expr": bson.M{
			"$eq": []interface{}{
				bson.M{"$replaceAll": bson.M{"input": "$show_name", "find": " ", "replacement": ""}},
				trimmedShowName, // Compare without spaces
			},
		},
	}).Decode(&show)

	// If no show is found, return a 404
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Show not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Error fetching show"})
	}

	return c.Status(200).JSON(show)
}

// Get show_name from another collection from magic_videos

//  a dj using first_name and last_name that is in a show in show_host collection

// GetShowByDjName - Find a show where the DJ is listed as a host
func GetShowByDjName(c fiber.Ctx) error {
	// Get first name and last name from request parameters
	firstName := c.Params("first_name")
	lastName := c.Params("last_name")

	if firstName == "" || lastName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "First name and last name are required"})
	}

	// Full name of the DJ
	fullName := firstName + " " + lastName

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Search for shows where the DJ is listed in the show_host field (case-insensitive substring match)
	var shows []models.Shows
	showCollection := ShowsCollectionInit()

	cursor, err := showCollection.Find(ctx, bson.M{
		"show_host": bson.M{
			"$regex": primitive.Regex{Pattern: fullName, Options: "i"}, // Case-insensitive substring match
		},
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error searching for shows"})
	}

	// Decode all matching shows
	if err := cursor.All(ctx, &shows); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error decoding shows"})
	}

	if len(shows) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "This DJ is not hosting any show"})
	}

	// Return all shows where the DJ is hosting
	return c.Status(200).JSON(fiber.Map{
		"message": "DJ found as a show host",
		"shows":   shows,
	})
}

// Create a Show
func CreateShow(c fiber.Ctx) error {
	var show models.Shows
	showCollection := ShowsCollectionInit()

	if err := c.Bind().JSON(&show); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate Show_desc
	if len(show.Show_desc) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Show_desc cannot be empty"})
	}

	show.ID = primitive.NewObjectID()
	show.Created_at = primitive.NewDateTimeFromTime(time.Now())
	show.Updated_at = show.Created_at

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := showCollection.InsertOne(ctx, show)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error creating show"})
	}

	return c.Status(201).JSON(show)
}

// Update a Show
func UpdateShow(c fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Parse the request body into a generic map
	var updateData map[string]interface{}
	if err := c.Bind().JSON(&updateData); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Convert string dates to primitive.DateTime if necessary
	if createdAtStr, ok := updateData["created_at"].(string); ok {
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err == nil {
			updateData["created_at"] = primitive.NewDateTimeFromTime(createdAt)
		} else {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid created_at date format"})
		}
	}

	if updatedAtStr, ok := updateData["updated_at"].(string); ok {
		updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
		if err == nil {
			updateData["updated_at"] = primitive.NewDateTimeFromTime(updatedAt)
		} else {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid updated_at date format"})
		}
	}

	// Validate and process Show_desc if present
	if desc, ok := updateData["show_desc"].([]interface{}); ok {
		var showDesc []string
		for _, d := range desc {
			if str, valid := d.(string); valid {
				showDesc = append(showDesc, str)
			} else {
				return c.Status(400).JSON(fiber.Map{"error": "Invalid Show_desc format"})
			}
		}
		updateData["show_desc"] = showDesc
	}

	// Add the updated_at field if not passed as part of the request
	if _, ok := updateData["updated_at"]; !ok {
		updateData["updated_at"] = primitive.NewDateTimeFromTime(time.Now())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": updateData} // Only update the provided fields

	showCollection := ShowsCollectionInit()

	_, err = showCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error updating show"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Show updated successfully"})
}

// Delete a Show
func DeleteShow(c fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	showCollection := ShowsCollectionInit()

	_, err = showCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error deleting show"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Show deleted successfully"})
}
