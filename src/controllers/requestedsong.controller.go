package controllers

import (
	"context"
	"log"
	"magic-server-2026/src/db"
	"magic-server-2026/src/helpers"
	"magic-server-2026/src/models"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RequestSongCollectionInit() *mongo.Collection {
	return db.GetCollection("magic899_db", "requestedsongs")
}

func VotesCollectionInit() *mongo.Collection {
	return db.GetCollection("magic899_db", "votes")
}

// isValidImageURL validates if the URL matches Spotify's image URL format
func isValidImageURL(url string) bool {
	// Must start with Spotify's known image URL prefix
	if url == "" || !strings.HasPrefix(url, "https://i.scdn.co/image/") {
		return false
	}

	// Match only the UUID-like segment after the base URL with no file extension
	re := regexp.MustCompile(`^https://i\.scdn\.co/image/[a-zA-Z0-9]{40}$`)
	return re.MatchString(url)
}

// Helper function to validate strings (allow alphanumeric, spaces, hyphens, apostrophes)
func isValidString(str string) bool {
	// Allow letters, numbers, whitespace characters, hyphens, apostrophes
	re := regexp.MustCompile(`^[a-zA-Z0-9\s\-']+$`)
	trimmed := strings.TrimSpace(str)
	return re.MatchString(trimmed) && trimmed != ""
}

func isValidAlbumName(str string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9\s\-\':&!?()\.]+$`)
	trimmed := strings.TrimSpace(str)
	return re.MatchString(trimmed) && trimmed != ""
}

// GetAllRequestSongs returns all request songs

func GetAllRequestSongs(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	requestSongCollection := RequestSongCollectionInit()

	var requestSongItems []models.Track

	cursor, err := requestSongCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch Requested Songs"})
	}

	if err = cursor.All(ctx, &requestSongItems); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode Requested Songs"})
	}

	// Convert `primitive.DateTime` to ISO8601 string format
	for i := range requestSongItems {
		requestSongItems[i].Created_at = primitive.NewDateTimeFromTime(time.Now())
		requestSongItems[i].Updated_at = primitive.NewDateTimeFromTime(time.Now())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Requested Songs fetched successfully",
		"songs":   requestSongItems,
	})
}

// GetRequestSong - Get a specific RequestSong by ID
func GetRequestSong(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	requestSongID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(requestSongID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Requested Song ID"})
	}

	// Initialize the collection dynamically
	requestSongCollection := RequestSongCollectionInit()

	var requestSong models.Track
	err = requestSongCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&requestSong)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Requested Song not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch Requested Song"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Requested Song fetched successfully", "song": requestSong})
}

// CreateRequestSong - Create a new RequestSong

func CreateRequestSong(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	requestSongCollection := RequestSongCollectionInit()
	votesCollection := VotesCollectionInit()

	// Parse and validate body
	var requestSong models.Track
	if err := c.Bind().JSON(&requestSong); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if !isValidString(requestSong.Name) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid song name"})
	}
	if len(requestSong.Artists) == 0 || !isValidString(requestSong.Artists[0].Name) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid artist name"})
	}
	if !isValidAlbumName(requestSong.Album.Name) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid album name"})
	}
	if len(requestSong.Album.Images) > 0 && requestSong.Album.Images[0].URL != "" && !isValidImageURL(requestSong.Album.Images[0].URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid image URL"})
	}
	if len(requestSong.RequestedBy) == 0 || !isValidString(requestSong.RequestedBy[0]) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid requester"})
	}

	// Requester and IP tracking
	requester := strings.TrimSpace(requestSong.RequestedBy[0])
	ip := c.IP()
	currentTime := time.Now()
	dateKey := currentTime.Format("2006-01-02")

	// Check eligibility in votes collection
	var voteRecord models.VoteRecord
	err := votesCollection.FindOne(ctx, bson.M{
		"ip_address": ip,
		"date":       dateKey,
	}).Decode(&voteRecord)

	if err == mongo.ErrNoDocuments || voteRecord.Votes < 3 {
		if err == mongo.ErrNoDocuments {
			voteRecord = models.VoteRecord{
				MusicID:    primitive.NilObjectID,
				Date:       dateKey,
				IPAddress:  ip,
				Votes:      1,
				LastVoteAt: currentTime,
				CanVote:    true,
				CanRequest: true,
			}
			_, err = votesCollection.InsertOne(ctx, voteRecord)
		} else {
			_, err = votesCollection.UpdateOne(ctx, bson.M{
				"ip_address": ip,
				"date":       dateKey,
			}, bson.M{
				"$inc": bson.M{"votes": 1},
				"$set": bson.M{
					"last_vote_at": currentTime,
					"can_request":  true,
				},
			})
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update voting record"})
		}
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Request limit reached. Try again after 72 hours."})
	}

	// Check if song exists
	filter := bson.M{
		"name":         requestSong.Name,
		"artists.name": requestSong.Artists[0].Name,
	}

	var existing models.Track
	err = requestSongCollection.FindOne(ctx, filter).Decode(&existing)
	switch err {
	case mongo.ErrNoDocuments:
		// New song
		requestSong.ID = primitive.NewObjectID()
		requestSong.RequestCount = 1
		requestSong.Created_at = primitive.NewDateTimeFromTime(currentTime)
		requestSong.Updated_at = primitive.NewDateTimeFromTime(currentTime)

		if _, err := requestSongCollection.InsertOne(ctx, requestSong); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create requested song"})
		}
	case nil:
		// Increase request count and append requester
		err := helpers.IncreaseRequestedCount(requestSongCollection, existing, requester)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update Requested Song"})
		}
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking for existing song"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Requested song created successfully",
	})
}

func UpdateRequestSong(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	songName := c.Params("name")
	if songName == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Song name is required"})
	}

	trimmedSongName := strings.ReplaceAll(songName, " ", "") // Remove all spaces from songName

	// Parse request body expecting JSON like: { "requestedBy": "requesterName" }
	var body struct {
		RequestedBy string `json:"requestedBy"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	requester := strings.TrimSpace(body.RequestedBy)
	if requester == "" || !isValidString(requester) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid requester"})
	}

	requestSongCollection := RequestSongCollectionInit()

	var existing models.Track
	err := requestSongCollection.FindOne(ctx, bson.M{
		"$expr": bson.M{
			"$eq": []interface{}{
				bson.M{"$replaceAll": bson.M{"input": "$name", "find": " ", "replacement": ""}},
				trimmedSongName,
			},
		},
	}).Decode(&existing)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Requested Song not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch Requested Song"})
	}

	// Prepare update operations:
	updateOps := bson.M{
		"$inc": bson.M{"requestcount": 1},
		"$set": bson.M{"updated_at": primitive.NewDateTimeFromTime(time.Now())},
	}

	// Add requester to RequestedBy array if not already present
	// Use $addToSet to prevent duplicates
	updateOps["$addToSet"] = bson.M{"requestedby": requester}

	updateResult, err := requestSongCollection.UpdateOne(ctx, bson.M{"_id": existing.ID}, updateOps)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update Requested Song"})
	}

	if updateResult.MatchedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Requested Song not found"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Requested Song updated successfully",
	})
}

// DeleteRequestSong - Delete a specific RequestSong by ID

func DeleteRequestSong(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	requestSongID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(requestSongID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Requested Song ID"})
	}

	// Initialize the collection dynamically
	requestSongCollection := RequestSongCollectionInit()

	// Delete the RequestSong document
	_, err = requestSongCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete Requested Song"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Requested Song deleted successfully"})
}

func OriginDomain() string {
	if os.Getenv("ENV") == "production" {
		return "onrender.com"
	}
	return "localhost"
}
