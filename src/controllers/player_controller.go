package controllers

import (
	"fmt"
	"magic-server-2026/src/utils"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v3"
)

func GetVideoPlayer(c fiber.Ctx) error {
	filename := c.Params("filename")

	// Serve video
	uploadDir := "./src/uploads/videos"
	filePath := filepath.Join(uploadDir, filename)

	contentType := "application/octet-stream"
	switch filepath.Ext(filename) {
	case ".mp4":
		contentType = "video/mp4"
	case ".avi":
		contentType = "video/x-msvideo"
	case ".mov":
		contentType = "video/quicktime"
	}

	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	c.Set("Content-Type", contentType)
	return c.SendFile(filePath)
}

func GetVideoPlayerSignedURL(c fiber.Ctx) error {
	filename := c.Params("filename")

	// ResourceTokenMiddleware already validates X-RSP-Token
	// Validate file exists
	uploadDir := "./src/uploads/videos"
	filePath := filepath.Join(uploadDir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File not found"})
	}

	// Generate signed URL (expires in 5 min)
	signedURL, err := utils.GenerateSignedURL(filename, 5*time.Minute)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot generate URL"})
	}

	return c.JSON(fiber.Map{"url": signedURL})
}
