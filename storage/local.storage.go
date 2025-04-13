package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type LocalUploader struct {
	BasePath string
}

// Upload saves the file locally under BasePath
func (l *LocalUploader) Upload(file io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %v", err)
	}

	uploadsDir := filepath.Join(workingDir, l.BasePath)
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create uploads directory: %v", err)
	}

	// Add 13-digit timestamp to filename
	ext := filepath.Ext(fileHeader.Filename)
	name := fileHeader.Filename[:len(fileHeader.Filename)-len(ext)]
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	newFilename := fmt.Sprintf("%s-%d%s", name, timestamp, ext)

	filePath := filepath.Join(uploadsDir, newFilename)
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

// Get opens a file for reading
func (l *LocalUploader) Get(filename string) (io.ReadCloser, error) {
	fullPath := filepath.Join(l.BasePath, filename)
	return os.Open(fullPath)
}

// List returns all files under BasePath
func (l *LocalUploader) List() ([]string, error) {
	files, err := os.ReadDir(l.BasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %v", err)
	}

	var filenames []string
	for _, f := range files {
		if !f.IsDir() {
			filenames = append(filenames, f.Name())
		}
	}
	return filenames, nil
}

// Delete removes a single file
func (l *LocalUploader) Delete(filename string) error {
	fullPath := filepath.Join(l.BasePath, filename)
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

// DeleteMultiple removes multiple files
func (l *LocalUploader) DeleteMultiple(filenames []string) error {
	for _, name := range filenames {
		if err := l.Delete(name); err != nil {
			return err
		}
	}
	return nil
}
