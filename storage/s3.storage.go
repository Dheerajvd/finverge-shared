package storage

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Uploader struct {
	Client     *s3.S3
	BucketName string
	BasePath   string
}

func (s *S3Uploader) Upload(file io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	fileKey := filepath.Join(s.BasePath, fileHeader.Filename)

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return "", fmt.Errorf("failed to read file into buffer: %v", err)
	}

	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fileKey),
		Body:   bytes.NewReader(buf.Bytes()),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.BucketName, fileKey)
	return fileURL, nil
}

func NewS3Uploader(region, accessKey, secretKey, bucketName, basePath string) (*S3Uploader, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	return &S3Uploader{
		Client:     s3.New(sess),
		BucketName: bucketName,
		BasePath:   basePath,
	}, nil
}
