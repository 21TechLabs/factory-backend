package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var CFG aws.Config
var Uploader *manager.Uploader

func init() {
	var err error
	CFG, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
}

func init() {
	client := s3.NewFromConfig(CFG)

	Uploader = manager.NewUploader(client)
}

func S3UploadFile(bucketName string, objectKey string, filePath *multipart.FileHeader) (manager.UploadOutput, error) {
	f, err := filePath.Open()

	if err != nil {
		return manager.UploadOutput{}, err
	}

	uploadOutput, err := Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%d-%s", time.Now().UnixNano(), objectKey)),
		Body:   f,
	})

	if err != nil {
		return manager.UploadOutput{}, err
	}

	return *uploadOutput, nil
}
