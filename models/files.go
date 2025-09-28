package models

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/21TechLabs/factory-backend/utils"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type FileStore struct {
	DB *gorm.DB
}

func NewFileStore(db *gorm.DB) *FileStore {
	return &FileStore{
		DB: db,
	}
}

type File struct {
	gorm.Model
	ID     uint   `json:"id" gorm:"primaryKey"`
	UserID uint   `gorm:"not null;index"`
	Name   string `json:"name" gorm:"column:name"`
	Type   string `json:"type" gorm:"column:file_type"`
	Etag   string `json:"-" gorm:"column:etag"`
	Size   int64  `json:"size" gorm:"column:size"`
	Key    string `json:"key" gorm:"column:key"`
	Bucket string `json:"-" gorm:"column:bucket"`
}

type FileUpload struct {
	Title string
	File  multipart.FileHeader
}

func (fs *FileStore) FileGetByID(id primitive.ObjectID) (*File, error) {
	var file File
	result := fs.DB.Model(&file).Where("id = ?", id).First(&file)
	if result.Error != nil {
		return nil, result.Error
	}
	return &file, nil
}

func (fs *FileStore) FileGetBy(filter map[string]interface{}) (*File, error) {
	var file File
	result := fs.DB.Model(&file).Where(filter).First(&file)
	if result.Error != nil {
		return nil, result.Error
	}
	return &file, nil
}

func (f *File) GetObject() (*minio.Object, error) {
	return utils.S3GetObject(f.Bucket, f.Key)
}

func (fs *FileStore) GetUrl(f *File) (string, error) {
	var ctx = context.Background()
	location, err := utils.MinioClient.PresignedGetObject(ctx, f.Bucket, f.Key, time.Hour*24, nil)
	if err != nil {
		return "", err
	}
	return location.String(), nil
}

func (f *File) Delete() error {
	err := utils.S3DeleteFile(f.Bucket, f.Key)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStore) UploadFile(files []FileUpload, userId uint) ([]File, error) {
	var fileObjs []File
	for _, file := range files {
		fileObj := File{
			UserID: userId,
			Name:   file.Title,
			Type:   file.File.Header.Get("Content-Type"),
			Size:   file.File.Size,
		}

		uploadOutput, err := utils.S3UploadFile(utils.S3BucketName, fmt.Sprintf("%d", userId), &file.File)
		if err != nil {
			return nil, err
		}

		fileObj.Etag = uploadOutput.ETag
		fileObj.Key = uploadOutput.Key
		fileObj.Bucket = utils.S3BucketName

		// err = mgm.Coll(&fileObj).Create(&fileObj)
		result := fs.DB.Create(&fileObj)
		if result.Error != nil {
			return nil, result.Error
		}

		fileObjs = append(fileObjs, fileObj)
	}
	return fileObjs, nil
}
