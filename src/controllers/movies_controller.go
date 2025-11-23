package controllers

import (
	"context"
	"fmt"
	"magic-server-2026/src/db"
	"magic-server-2026/src/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// -------------------------------
// Collection Initialization
// -------------------------------

func MoviesCollectionInit() *mongo.Collection {

	return db.GetCollection("magic899_db", "movies")
}

// -------------------------------
// Helper Functions
// -------------------------------

func jsonResponse(c fiber.Ctx, status int, message string, data fiber.Map) error {
	resp := fiber.Map{"message": message}
	for k, v := range data {
		resp[k] = v
	}
	return c.Status(status).JSON(resp)
}

func errorResponse(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}

// -------------------------------
// Handlers
// -------------------------------

// Get all movies
func GetMovies(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var movies []models.Movies
	moviesCollection := MoviesCollectionInit()
	cursor, err := moviesCollection.Find(ctx, bson.M{})
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, "Failed to fetch movies")
	}
	if err = cursor.All(ctx, &movies); err != nil {
		return errorResponse(c, http.StatusInternalServerError, "Failed to decode movies")
	}

	return jsonResponse(c, http.StatusOK, "Movies fetched successfully", fiber.Map{"movies": movies})
}

// Get single movie by ID
func GetMovie(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, "Invalid Movie ID")
	}

	var movie models.Movies
	moviesCollection := MoviesCollectionInit()

	err = moviesCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errorResponse(c, http.StatusNotFound, "Movie not found")
		}
		return errorResponse(c, http.StatusInternalServerError, "Failed to fetch movie")
	}

	return jsonResponse(c, http.StatusOK, "Movie fetched successfully", fiber.Map{"movie": movie})
}

// Get movie by name (case-insensitive, ignore spaces)
func GetByMovieName(c fiber.Ctx) error {
	movieName := strings.TrimSpace(c.Params("name"))
	if movieName == "" {
		return errorResponse(c, http.StatusBadRequest, "Movie name is required")
	}

	trimmedName := strings.ReplaceAll(movieName, " ", "")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var movie models.Movies
	moviesCollection := MoviesCollectionInit()

	err := moviesCollection.FindOne(ctx, bson.M{
		"$expr": bson.M{
			"$eq": []interface{}{
				bson.M{"$replaceAll": bson.M{"input": "$title", "find": " ", "replacement": ""}},
				trimmedName,
			},
		},
	}).Decode(&movie)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errorResponse(c, http.StatusNotFound, "Movie not found")
		}
		return errorResponse(c, http.StatusInternalServerError, "Error fetching movie")
	}

	return jsonResponse(c, http.StatusOK, "Movie fetched successfully", fiber.Map{"movie": movie})
}

// Get the latest movie
func GetLatestMovie(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var movie models.Movies
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	moviesCollection := MoviesCollectionInit()

	err := moviesCollection.FindOne(ctx, bson.M{}, opts).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errorResponse(c, http.StatusNotFound, "No movies found")
		}
		return errorResponse(c, http.StatusInternalServerError, "Failed to fetch latest movie")
	}

	return jsonResponse(c, http.StatusOK, "Latest movie fetched successfully", fiber.Map{"movie": movie})
}

// Get upcoming movies after the latest release
func GetUpNextMovies(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var latestMovie models.Movies
	opts := options.FindOne().SetSort(bson.D{{Key: "release_date", Value: -1}})
	moviesCollection := MoviesCollectionInit()

	if err := moviesCollection.FindOne(ctx, bson.M{}, opts).Decode(&latestMovie); err != nil {
		if err == mongo.ErrNoDocuments {
			return errorResponse(c, http.StatusNotFound, "No movies found")
		}
		return errorResponse(c, http.StatusInternalServerError, "Failed to fetch latest movie")
	}

	filter := bson.M{"release_date": bson.M{"$gt": latestMovie.Release_date}}
	findOpts := options.Find().SetSort(bson.D{{Key: "release_date", Value: 1}}).SetLimit(5)
	cursor, err := moviesCollection.Find(ctx, filter, findOpts)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, "Failed to fetch upcoming movies")
	}
	defer cursor.Close(ctx)

	var upNext []models.Movies
	if err := cursor.All(ctx, &upNext); err != nil {
		return errorResponse(c, http.StatusInternalServerError, "Failed to decode upcoming movies")
	}

	return jsonResponse(c, http.StatusOK, "Upcoming movies fetched successfully", fiber.Map{"movies": upNext})
}

// Get all previous movies except the latest one
func GetPreviousMovies(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(1)
	moviesCollection := MoviesCollectionInit()

	cursor, err := moviesCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, "Failed to fetch previous movies")
	}
	defer cursor.Close(ctx)

	var movies []models.Movies
	if err := cursor.All(ctx, &movies); err != nil {
		return errorResponse(c, http.StatusInternalServerError, "Failed to decode previous movies")
	}

	return jsonResponse(c, http.StatusOK, "Previous movies fetched successfully", fiber.Map{"movies": movies})
}

// Create a new movie
func CreateMovie(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var movie models.Movies
	if err := c.Bind().Body(&movie); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	movie.ID = primitive.NewObjectID()
	movie.Created_at = primitive.NewDateTimeFromTime(time.Now())
	movie.Updated_at = movie.Created_at
	moviesCollection := MoviesCollectionInit()

	if _, err := moviesCollection.InsertOne(ctx, movie); err != nil {
		return errorResponse(c, http.StatusInternalServerError, "Failed to create movie")
	}

	return jsonResponse(c, http.StatusCreated, "Movie created successfully", fiber.Map{"movie": movie})
}

// Update a movie
func UpdateMovie(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, "Invalid Movie ID")
	}

	var updateData map[string]interface{}
	if err := c.Bind().Body(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	delete(updateData, "created_at")
	updateData["updated_at"] = primitive.NewDateTimeFromTime(time.Now())
	moviesCollection := MoviesCollectionInit()

	updateResult, err := moviesCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, "Failed to update movie")
	}

	if updateResult.MatchedCount == 0 {
		return errorResponse(c, http.StatusNotFound, "Movie not found")
	}

	var updatedMovie models.Movies
	if err := moviesCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedMovie); err != nil {
		return errorResponse(c, http.StatusInternalServerError, "Failed to fetch updated movie")
	}

	return jsonResponse(c, http.StatusOK, "Movie updated successfully", fiber.Map{"movie": updatedMovie})
}

// Delete a movie
func DeleteMovie(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, "Invalid movie ID")
	}
	moviesCollection := MoviesCollectionInit()

	result, err := moviesCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, "Failed to delete movie")
	}

	if result.DeletedCount == 0 {
		return errorResponse(c, http.StatusNotFound, "Movie not found")
	}

	return jsonResponse(c, http.StatusOK, "Movie deleted successfully", fiber.Map{})
}

// Stream (play) a movie file securely
func GetMovieHandler(c fiber.Ctx) error {
	filename := c.Params("filename")

	if strings.Contains(filename, "..") {
		return errorResponse(c, http.StatusBadRequest, "Invalid file path")
	}

	uploadDir := "./server/uploads"
	filePath := filepath.Join(uploadDir, filename)

	if filepath.Ext(filename) != ".mp4" {
		return errorResponse(c, http.StatusBadRequest, "Invalid file type. Only MP4 files are allowed.")
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errorResponse(c, http.StatusNotFound, "File not found")
	}

	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	c.Set("Content-Type", "video/mp4")
	return c.SendFile(filePath) // Enable range streaming
}
