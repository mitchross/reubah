package convert

import (
	"fmt"
	"image"
)

// Converter defines the interface for image format conversion
type Converter interface {
	Convert(img image.Image) (image.Image, error)
	Format() string
	Quality() int
	SetQuality(quality int)
	// Adding new method to check if format supports transparency
	SupportsTransparency() bool
}

// Factory returns the appropriate converter for the given format
func Factory(format string) (Converter, error) {
	switch format {
	case FormatJPEG, FormatJPG:
		return NewJPEGConverter(), nil
	case FormatPNG:
		return NewPNGConverter(), nil
	case FormatWebP:
		return NewWebPConverter(), nil
	case FormatGIF:
		return NewGIFConverter(), nil
	case FormatBMP:
		return NewBMPConverter(), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
