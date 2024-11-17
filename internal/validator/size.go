package validator

import (
	"github.com/dendianugerah/reubah/pkg/errors"
)

const (
	MaxFileSize    = 10 << 20 // 10 MB
	MaxImageWidth  = 8192
	MaxImageHeight = 8192
)

func ValidateFileSize(size int64) error {
	if size > MaxFileSize {
		return errors.New(
			errors.ErrInvalidSize,
			"File size exceeds maximum allowed size (10MB)",
			nil,
		)
	}
	return nil
}

func ValidateImageDimensions(width, height int) error {
	if width > MaxImageWidth || height > MaxImageHeight {
		return errors.New(
			errors.ErrInvalidSize,
			"Image dimensions exceed maximum allowed size (8192x8192)",
			nil,
		)
	}
	return nil
} 