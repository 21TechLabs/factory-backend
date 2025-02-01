package controllers

import (
	"github.com/21TechLabs/factory-be/models"
	"github.com/21TechLabs/factory-be/utils"
	"github.com/gofiber/fiber/v2"
)

func FileUpload(c *fiber.Ctx) error {

	currentUser, ok := c.Locals("user").(models.User)

	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user not found")
	}

	form, err := c.MultipartForm()

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	files := form.File["files"]

	title := form.Value["title"]

	if len(files) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "no file uploaded")
	}

	if len(title) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "title is required")
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
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"files":   uploadedFiles,
	})
}
