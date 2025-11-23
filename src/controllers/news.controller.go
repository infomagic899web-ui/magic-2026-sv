package controllers

import (
	"context"
	"log"
	"magic-server-2026/src/db"
	"magic-server-2026/src/helpers"
	"magic-server-2026/src/models"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
   News Controller for News Posts
   -----------------------------------
   1. Get All Posts
   2. Get a Post
   3. Create a Post
   4. Update a Post
   5. Delete a Post
   -----------------------------------
   PATH: /api/v1/news
*/

func NormalizeTitle(title string) string {
	normalized := strings.ToLower(title)
	normalized = strings.ReplaceAll(normalized, " ", "")
	reg := regexp.MustCompile("[^a-z0-9]")
	normalized = reg.ReplaceAllString(normalized, "")
	return normalized
}

func NewsCollectionInit() *mongo.Collection {
	return db.GetCollection("magic899_db", "news")
}

// GetNews - Get all news items
func GetNews(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var news []models.News

	newsCollection := NewsCollectionInit()
	cursor, err := newsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Find error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch news"})
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &news); err != nil {
		log.Println("Cursor decode error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse news"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "News fetched successfully",
		"news":    news,
	})
}

// GetNewsItem - Get a single news item by ID
func GetNewsItem(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newsID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(newsID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid News ID"})
	}

	var news models.News
	newsCollection := NewsCollectionInit()

	err = newsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&news)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "News not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch news item"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "News item fetched successfully",
		"news":    news,
	})
}

// GetNewsBySlug - Get a news item based on slugified title
func GetNewsBySlug(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing slug parameter",
		})
	}

	// Query by exact slug
	filter := bson.M{"slug": slug}

	var news models.News
	newsCollection := NewsCollectionInit()

	err := newsCollection.FindOne(ctx, filter).Decode(&news)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "News item not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch news",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "News fetched successfully",
		"news":    news,
	})
}

// GetNewsByTitle - Get a single news item by title
func GetNewsByTitle(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rawTitle := c.Params("title")
	if rawTitle == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing title parameter"})
	}

	normalizedParam := NormalizeTitle(rawTitle)

	var news models.News
	newsCollection := NewsCollectionInit()

	err := newsCollection.FindOne(ctx, bson.M{"normalized_title": normalizedParam}).Decode(&news)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "News item not found"})
		}
		log.Println("DB error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch news"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "News fetched successfully",
		"news":    news,
	})
}

// GetNewsByTitleAndCategory - Get a single news item by title and category
func GetNewsByTitleAndCategory(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newsTitle := c.Params("name")
	newsCategory := c.Params("category")

	decodedTitle, err := url.QueryUnescape(strings.ReplaceAll(newsTitle, "-", " "))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid title"})
	}

	normalizedTitle := helpers.NormalizeName(decodedTitle)
	var news models.News
	newsCollection := NewsCollectionInit()

	err = newsCollection.FindOne(ctx, bson.M{
		"normalized_title": normalizedTitle,
		"category":         bson.M{"$regex": newsCategory, "$options": "i"},
	}).Decode(&news)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "News item not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch news"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "News fetched successfully",
		"news":    news,
	})
}

// CreateNewsItem - Create a new news item
func CreateNewsItem(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var news models.News
	if err := c.Bind().JSON(&news); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if news.Title == "" || news.Writer == "" || news.Category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title, Writer, and Category are required"})
	}

	news.ID = primitive.NewObjectID()
	news.Created_at = primitive.NewDateTimeFromTime(time.Now())
	news.Updated_at = primitive.NewDateTimeFromTime(time.Now())
	newsCollection := NewsCollectionInit()

	result, err := newsCollection.InsertOne(ctx, news)
	if err != nil {
		log.Println("Error inserting news:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create news item"})
	}

	log.Println("News created with ID:", result.InsertedID)
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "News created successfully",
		"news":    news,
	})
}

// UpdateNewsItem - Update an existing news item
func UpdateNewsItem(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newsID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(newsID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid News ID"})
	}

	var updateData map[string]interface{}
	if err := c.Bind().JSON(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	delete(updateData, "created_at")
	updateData["updated_at"] = primitive.NewDateTimeFromTime(time.Now())
	newsCollection := NewsCollectionInit()

	updateResult, err := newsCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update news item"})
	}
	if updateResult.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "News item not found"})
	}

	var updatedNews models.News
	err = newsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedNews)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated news"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "News item updated successfully",
		"news":    updatedNews,
	})
}

// DeleteNewsItem - Delete a news item by ID
func DeleteNews(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newsID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(newsID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid News ID"})
	}
	newsCollection := NewsCollectionInit()

	_, err = newsCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete news item"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "News deleted successfully"})
}
