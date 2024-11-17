package convert

import (
	"bytes"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
)

type GIFConverter struct {
	numColors int
}

func NewGIFConverter() *GIFConverter {
	return &GIFConverter{numColors: DefaultGIFColors}
}

func (c *GIFConverter) Convert(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	paletted := image.NewPaletted(bounds, palette.Plan9)
	draw.FloydSteinberg.Draw(paletted, bounds, img, bounds.Min)
	
	// Encode to GIF and decode back to ensure proper conversion
	var buf bytes.Buffer
	if err := gif.Encode(&buf, paletted, &gif.Options{
		NumColors: c.numColors,
	}); err != nil {
		return nil, err
	}

	// Decode back to image.Image
	decoded, err := gif.Decode(&buf)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func (c *GIFConverter) Format() string {
	return FormatGIF
}

func (c *GIFConverter) Quality() int {
	return (c.numColors * 100) / 256 // Convert colors (0-256) to quality (0-100)
}

func (c *GIFConverter) SetQuality(quality int) {
	// Convert quality (0-100) to number of colors (2-256)
	c.numColors = (quality * 256) / 100
	if c.numColors < 2 {
		c.numColors = 2
	} else if c.numColors > 256 {
		c.numColors = 256
	}
}

func (c *GIFConverter) SupportsTransparency() bool {
	return true
} 