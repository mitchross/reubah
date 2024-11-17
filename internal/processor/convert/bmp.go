package convert

import (
	"image"
)

type BMPConverter struct{}

func NewBMPConverter() *BMPConverter {
	return &BMPConverter{}
}

func (c *BMPConverter) Convert(img image.Image) (image.Image, error) {
	// BMP typically doesn't support alpha, convert to RGB
	return convertToRGB(img), nil
}

func (c *BMPConverter) Format() string {
	return FormatBMP
}

func (c *BMPConverter) Quality() int {
	return MaxQuality // BMP is lossless
}

func (c *BMPConverter) SetQuality(quality int) {
	// BMP doesn't support quality settings
}

func (c *BMPConverter) SupportsTransparency() bool {
	return false
}
