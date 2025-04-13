package storage

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

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

// Upload a file to S3
func (s *S3Uploader) Upload(file io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	// Extract original filename components
	ext := filepath.Ext(fileHeader.Filename)
	name := fileHeader.Filename[:len(fileHeader.Filename)-len(ext)]
	timestamp := time.Now().UnixNano() / int64(time.Millisecond) // 13-digit timestamp
	newFilename := fmt.Sprintf("%s-%d%s", name, timestamp, ext)

	// Build the S3 object key
	fileKey := filepath.Join(s.BasePath, newFilename)

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return "", fmt.Errorf("failed to read file into buffer: %v", err)
	}

	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fileKey),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.BucketName, fileKey)
	return fileURL, nil
}

// GetObject returns a file reader for a specific key
func (s *S3Uploader) GetObject(key string) (io.ReadCloser, error) {
	fullKey := filepath.Join(s.BasePath, key)
	resp, err := s.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fullKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %v", err)
	}
	return resp.Body, nil
}

// ListObjects lists all objects under the base path
func (s *S3Uploader) ListObjects() ([]string, error) {
	resp, err := s.Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(s.BucketName),
		Prefix: aws.String(s.BasePath),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects in S3: %v", err)
	}

	var keys []string
	for _, item := range resp.Contents {
		keys = append(keys, *item.Key)
	}
	return keys, nil
}

// DeleteObject deletes a single object from S3
func (s *S3Uploader) DeleteObject(key string) error {
	fullKey := filepath.Join(s.BasePath, key)
	_, err := s.Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fullKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %v", err)
	}
	return nil
}

// DeleteMultipleObjects deletes multiple objects from S3
func (s *S3Uploader) DeleteMultipleObjects(keys []string) error {
	var objects []*s3.ObjectIdentifier
	for _, key := range keys {
		fullKey := filepath.Join(s.BasePath, key)
		objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(fullKey)})
	}

	_, err := s.Client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(s.BucketName),
		Delete: &s3.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete multiple objects: %v", err)
	}
	return nil
}

// NewS3Uploader initializes the uploader with AWS credentials and config
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
