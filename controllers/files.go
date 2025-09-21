package controllers

import (
	"fmt"
	"log"

	"github.com/21TechLabs/musiclms-backend/models"
	"github.com/21TechLabs/musiclms-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type FileController struct {
	Logger    *log.Logger
	FileStore *models.FileStore
	UserStore *models.UserStore
}

func NewFileController(log *log.Logger, fs *models.FileStore, us *models.UserStore) *FileController {
	return &FileController{
		Logger:    log,
		FileStore: fs,
		UserStore: us,
	}
}

func (fc *FileController) FileUpload(c *fiber.Ctx) error {

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

	uploadedFiles, err := fc.UserStore.UploadFile(&currentUser, fileUploads)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"files":   uploadedFiles,
	})
}

func (fc *FileController) FileStreamS3(c *fiber.Ctx) error {
	fileKey := c.Params("fileKey")

	if fileKey == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "file id is required")
	}

	file, err := fc.FileStore.FileGetBy(map[string]interface{}{"key": fileKey})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "file not found")
	}

	fileObject, err := file.GetObject()

	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}
	defer fileObject.Close()
	c.Response().Header.Set("Content-Type", file.Type)
	c.Response().Header.Set("Content-Disposition", "inline; filename="+fileKey) // Or "attachment; filename=" for download
	c.Response().Header.Set("Content-Length", fmt.Sprint(file.Size))
	c.Response().Header.Set("Transfer-Encoding", "chunked")

	// if err := c.SendStream(fileObject, 2048); err != nil {
	// 	return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to stream file")
	// }
	fileByte := make([]byte, 2048) // Adjust the size as needed

	for {
		n, err := fileObject.Read(fileByte)
		if err != nil {
			if err.Error() == "EOF" {
				break // End of file reached
			}
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to read file")
		}
		if n == 0 {
			break // No more data to read
		}
		if _, err := c.Write(fileByte[:n]); err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to write file")
		}
	}
	return nil
}
