package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

var (
	S3BucketName = ""
)

func init() {
	LoadEnv()

	S3BucketName = GetEnv("AWS_S3_BUCKET", false)
	var err error
	// load AWS configuration from .env
	var (
		// AWS_REGION      = GetEnv("AWS_REGION", false)
		// bucket          = GetEnv("AWS_S3_BUCKET", false)
		endpoint        = GetEnv("AWS_S3_ORIGIN", false)
		accessKeyID     = GetEnv("AWS_ACCESS_KEY_ID", false)
		secretAccessKey = GetEnv("AWS_SECRET_ACCESS_KEY", false)
		useSSL          = GetEnv("AWS_S3_SECURE", false) == "true"
	)

	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

}

func S3GetObject(bucketName string, objectKey string) (*minio.Object, error) {
	var ctx = context.Background()

	object, err := MinioClient.GetObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return object, nil
}

func S3UploadFile(bucketName string, objectKey string, filePath *multipart.FileHeader) (minio.UploadInfo, error) {

	var ctx = context.Background()

	var objectName = fmt.Sprintf("%d-%s-%s", time.Now().UnixNano(), objectKey, filePath.Filename)
	var contentType = filePath.Header.Get("Content-Type")

	fp, err := filePath.Open()

	if err != nil {
		return minio.UploadInfo{}, err
	}

	u, err := url.Parse(objectName) // Example of parsing a URL with spaces
	if err != nil {
		return minio.UploadInfo{}, fmt.Errorf("failed to parse URL: %w", err)
	}
	objectName = u.String()

	uploadOutput, err := MinioClient.PutObject(ctx, bucketName, objectName, fp, filePath.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, uploadOutput.Size)

	if err != nil {
		return minio.UploadInfo{}, err
	}

	return uploadOutput, nil
}

func S3DeleteFile(bucketName string, objectKey string) error {
	var ctx = context.Background()

	err := MinioClient.RemoveObject(ctx, bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}

	log.Printf("Successfully deleted %s from bucket %s\n", objectKey, bucketName)
	return nil
}
