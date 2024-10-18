package controllers

import (
	"github.com/21TechLabs/factory-be/models"
	"github.com/gofiber/fiber/v2"
)

func FileUpload(c *fiber.Ctx) error {

	currentUser := c.Locals("user").(models.User)

	form, err := c.MultipartForm()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	files := form.File["files"]

	title := form.Value["title"]

	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "No files uploaded",
			"success": false,
		})
	}

	if len(title) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Title is required",
			"success": false,
		})
	}

	var fileUploads []models.FileUpload

	for _, file := range files {
		fileUploads = append(fileUploads, models.FileUpload{
			Title: title[0],
			File:  *file,
		})
	}

	uploadedFiles, err := currentUser.UploadFile(fileUploads)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"files":   uploadedFiles,
	})
}
