package convert

import (
	"bytes"
	"image"
	"image/png"
)

type PNGConverter struct {
	compression int // 0-9
}

func NewPNGConverter() *PNGConverter {
	return &PNGConverter{compression: DefaultPNGQuality}
}

func (c *PNGConverter) Convert(img image.Image) (image.Image, error) {
	// Convert to RGBA for consistent handling
	rgbaImg := convertToRGBA(img)
	
	// Encode to PNG with compression
	var buf bytes.Buffer
	encoder := png.Encoder{
		CompressionLevel: png.CompressionLevel(c.compression),
	}
	
	if err := encoder.Encode(&buf, rgbaImg); err != nil {
		return nil, err
	}

	// Decode back to image.Image
	decoded, err := png.Decode(&buf)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func (c *PNGConverter) Format() string {
	return FormatPNG
}

func (c *PNGConverter) Quality() int {
	return (c.compression * 100) / 9 // Convert compression (0-9) to quality (0-100)
}

func (c *PNGConverter) SetQuality(quality int) {
	// Convert quality (0-100) to compression (0-9)
	c.compression = (quality * 9) / 100
	if c.compression < 0 {
		c.compression = 0
	} else if c.compression > 9 {
		c.compression = 9
	}
}

func (c *PNGConverter) SupportsTransparency() bool {
	return true
}
