package convert

import (
	"bytes"
	"image"
	"image/jpeg"
)

type JPEGConverter struct {
	quality int
}

func NewJPEGConverter() *JPEGConverter {
	return &JPEGConverter{quality: DefaultJPEGQuality}
}

func (c *JPEGConverter) Convert(img image.Image) (image.Image, error) {
	// Convert to RGB and encode to JPEG
	rgbImg := convertToRGB(img)
	
	// Encode to JPEG with quality setting
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, rgbImg, &jpeg.Options{
		Quality: c.quality,
	}); err != nil {
		return nil, err
	}

	// Decode back to image.Image
	decoded, err := jpeg.Decode(&buf)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func (c *JPEGConverter) Format() string {
	return FormatJPEG
}

func (c *JPEGConverter) Quality() int {
	return c.quality
}

func (c *JPEGConverter) SetQuality(quality int) {
	if quality < MinQuality {
		quality = MinQuality
	} else if quality > MaxQuality {
		quality = MaxQuality
	}
	c.quality = quality
}

func (c *JPEGConverter) SupportsTransparency() bool {
	return false
}
