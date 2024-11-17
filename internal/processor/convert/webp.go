package convert

import (
	"bytes"
	"image"

	"github.com/chai2010/webp"
)

type WebPConverter struct {
	quality  int
	lossless bool
}

func NewWebPConverter() *WebPConverter {
	return &WebPConverter{
		quality:  DefaultWebPQuality,
		lossless: false,
	}
}

func (c *WebPConverter) Convert(img image.Image) (image.Image, error) {
	// Convert to RGBA first for consistent encoding
	rgba := convertToRGBA(img)

	// Encode to WebP
	var buf bytes.Buffer
	var err error

	if c.lossless {
		err = webp.Encode(&buf, rgba, &webp.Options{
			Lossless: true,
		})
	} else {
		err = webp.Encode(&buf, rgba, &webp.Options{
			Lossless: false,
			Quality:  float32(c.quality),
		})
	}
	if err != nil {
		return nil, err
	}

	// Decode back to image.Image
	decoded, err := webp.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func (c *WebPConverter) Format() string {
	return FormatWebP
}

func (c *WebPConverter) Quality() int {
	if c.lossless {
		return MaxQuality
	}
	return c.quality
}

func (c *WebPConverter) SetQuality(quality int) {
	if quality < MinQuality {
		quality = MinQuality
	} else if quality > MaxQuality {
		quality = MaxQuality
	}
	c.quality = quality
	c.lossless = quality == MaxQuality
}

func (c *WebPConverter) SupportsTransparency() bool {
	return true
}

// Additional methods specific to WebP
func (c *WebPConverter) SetLossless(lossless bool) {
	c.lossless = lossless
	if lossless {
		c.quality = MaxQuality
	}
}

func (c *WebPConverter) IsLossless() bool {
	return c.lossless
}
