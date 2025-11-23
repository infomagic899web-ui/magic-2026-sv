package controllers

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
)

// UploadImageHandler handles image uploads
func UploadImageHandler(c fiber.Ctx) error {
	// Extract tokens

	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Save file
	err = c.SaveFile(file, "./src/uploads/images/"+file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "File uploaded successfully",
		"file_name": file.Filename,
	})
}

// GetImageHandler serves an image from the uploads directory
func GetImageHandler(c fiber.Ctx) error {

	// Get filename
	filename := c.Params("filename")
	filePath := filepath.Join("./src/uploads/images", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "File not found",
		})
	}

	return c.SendFile(filePath)
}
