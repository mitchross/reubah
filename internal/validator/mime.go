package validator

import (
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/dendianugerah/reubah/pkg/errors"
)

var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
	"image/bmp":  true,
	"application/pdf": true,
}

func ValidateMIMEType(file multipart.File) error {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return errors.New(errors.ErrInvalidMIME, "Failed to read file", err)
	}

	// Reset the file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		return errors.New(errors.ErrInvalidMIME, "Failed to reset file pointer", err)
	}

	mimeType := http.DetectContentType(buffer)
	if !allowedMIMETypes[strings.ToLower(mimeType)] {
		return errors.New(errors.ErrInvalidMIME, "Unsupported file type: "+mimeType, nil)
	}

	return nil
} 