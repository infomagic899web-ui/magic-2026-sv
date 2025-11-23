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

// PostCollectionInit returns the MongoDB collection for posts
func PostCollectionInit() *mongo.Collection {
	return db.GetCollection("magic899_db", "posts")
}

// GetPosts - Retrieve all posts
func GetPosts(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postCollection := PostCollectionInit()

	var posts []models.Post
	cursor, err := postCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Find error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch posts"})
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &posts); err != nil {
		log.Println("Decode error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode posts"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Posts fetched successfully",
		"posts":   posts,
	})
}

// GetPost - Retrieve a single post by ID
func GetPost(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Post ID"})
	}

	postCollection := PostCollectionInit()
	var post models.Post

	err = postCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
		}
		log.Println("FindOne error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch post"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Post fetched successfully",
		"post":    post,
	})
}

// CreatePost - Create a new post
func CreatePost(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var post models.Post
	if err := c.Bind().JSON(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	post.ID = primitive.NewObjectID()
	post.Created_at = primitive.NewDateTimeFromTime(time.Now())
	post.Updated_at = primitive.NewDateTimeFromTime(time.Now())

	postCollection := PostCollectionInit()
	_, err := postCollection.InsertOne(ctx, post)
	if err != nil {
		log.Println("InsertOne error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create post"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Post created successfully",
		"post":    post,
	})
}

// UpdatePost - Update an existing post
func UpdatePost(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Post ID"})
	}

	var updateData map[string]interface{}
	if err := c.Bind().JSON(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	delete(updateData, "created_at")
	updateData["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	postCollection := PostCollectionInit()
	updateResult, err := postCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		log.Println("UpdateOne error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update post"})
	}

	if updateResult.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}

	var updatedPost models.Post
	err = postCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedPost)
	if err != nil {
		log.Println("FindOne after update error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated post"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Post updated successfully",
		"post":    updatedPost,
	})
}

// DeletePost - Delete a post by ID
func DeletePost(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Post ID"})
	}

	postCollection := PostCollectionInit()
	_, err = postCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Println("DeleteOne error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete post"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Post deleted successfully"})
}
