package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalUploader struct {
	BasePath string
}

func (l *LocalUploader) Upload(file io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %v", err)
	}

	uploadsDir := filepath.Join(workingDir, l.BasePath)
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create uploads directory: %v", err)
	}

	filePath := filepath.Join(uploadsDir, fileHeader.Filename)
	destFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return filePath, nil
}
