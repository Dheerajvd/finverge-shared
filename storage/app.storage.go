package storage

import (
	"io"
	"mime/multipart"
)

type Uploader interface {
	Upload(file io.Reader, fileHeader *multipart.FileHeader) (string, error)
}
