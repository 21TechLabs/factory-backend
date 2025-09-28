package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/21TechLabs/factory-backend/models"
	"github.com/21TechLabs/factory-backend/utils"
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

func (fc *FileController) FileUpload(w http.ResponseWriter, r *http.Request) {
	currentUser, err := utils.ReadContextValue[*models.User](r, utils.UserContextKey)

	if err != nil || currentUser == nil {
		utils.ErrorResponse(fc.Logger, w, http.StatusUnauthorized, []byte("User not found"))
		return
	}

	form := r.MultipartForm

	files := form.File["files"]
	title := form.Value["title"]

	if len(files) == 0 {
		utils.ErrorResponse(fc.Logger, w, http.StatusBadRequest, []byte("No file uploaded"))
		return
	}

	if len(title) == 0 {
		utils.ErrorResponse(fc.Logger, w, http.StatusBadRequest, []byte("Title is required"))
		return
	}

	var fileUploads []models.FileUpload

	for _, file := range files {
		fileUploads = append(fileUploads, models.FileUpload{
			Title: title[0],
			File:  *file,
		})
	}

	uploadedFiles, err := fc.UserStore.UploadFile(currentUser, fileUploads)

	if err != nil {
		utils.ErrorResponse(fc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	utils.ResponseWithJSON(fc.Logger, w, http.StatusOK, utils.Map{
		"success": true,
		"files":   uploadedFiles,
	})
}

func (fc *FileController) FileStreamS3(w http.ResponseWriter, r *http.Request) {
	fileKey := r.PathValue("fileKey")

	if fileKey == "" {
		utils.ErrorResponse(fc.Logger, w, http.StatusBadRequest, []byte("file id is required"))
		return
	}

	file, err := fc.FileStore.FileGetBy(map[string]interface{}{"key": fileKey})
	if err != nil {
		utils.ErrorResponse(fc.Logger, w, http.StatusNotFound, []byte("file not found"))
		return
	}

	fileObject, err := file.GetObject()

	if err != nil {
		utils.ErrorResponse(fc.Logger, w, http.StatusNotFound, []byte(err.Error()))
		return
	}
	defer fileObject.Close()

	w.Header().Set("Content-Type", file.Type)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", fileKey)) // Or "attachment; filename=" for download
	w.Header().Set("Content-Length", fmt.Sprint(file.Size))
	w.Header().Set("Transfer-Encoding", "chunked")

	fileByte := make([]byte, 2048) // Adjust the size as needed

	for {
		n, err := fileObject.Read(fileByte)
		if err != nil {
			if err.Error() == "EOF" {
				break // End of file reached
			}
			utils.ErrorResponse(fc.Logger, w, http.StatusInternalServerError, []byte("failed to read file"))
			return
		}
		if n == 0 {
			break // No more data to read
		}
		if _, err := w.Write(fileByte[:n]); err != nil {
			utils.ErrorResponse(fc.Logger, w, http.StatusInternalServerError, []byte("failed to write file"))
			return
		}
	}
}
