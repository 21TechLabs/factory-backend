package models

import (
	"mime/multipart"

	"github.com/kamva/mgm/v3"
)

type File struct {
	mgm.DefaultModel `bson:",inline"`
	ID               uint   `json:"id" gorm:"primaryKey"`
	UserId           string `json:"user_id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	RelativeURL      string `json:"url"`
}

type FileUpload struct {
	Title string
	File  multipart.FileHeader
}
