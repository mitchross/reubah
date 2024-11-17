package convert

import (
	"image"
	"image/color"
)

type ConversionManager struct {
	options *ConversionOptions
}

func NewManager(options *ConversionOptions) *ConversionManager {
	if options == nil {
		options = DefaultOptions(FormatJPEG)
	}
	return &ConversionManager{options: options}
}

func (m *ConversionManager) Convert(img image.Image, targetFormat string) (image.Image, error) {
	// Get the appropriate converter
	converter, err := Factory(targetFormat)
	if err != nil {
		return nil, err
	}

	// Apply quality settings
	if m.options.Quality > 0 {
		converter.SetQuality(m.options.Quality)
	}

	// Pre-process image based on options
	img = m.preprocess(img)

	// Perform the conversion
	converted, err := converter.Convert(img)
	if err != nil {
		return nil, err
	}

	// Post-process image based on options
	return m.postprocess(converted)
}

func (m *ConversionManager) preprocess(img image.Image) image.Image {
	// Handle color quantization if needed
	if m.options.PNG != nil && m.options.PNG.QuantizeColors {
		img = quantizeColors(img, m.options.PNG.MaxColors)
	}

	// Handle transparency
	if !getConverter(m.options.Format).SupportsTransparency() {
		img = removeTransparency(img)
	}

	return img
}

func (m *ConversionManager) postprocess(img image.Image) (image.Image, error) {
	// Additional processing can be added here
	return img, nil
}

// Helper functions
func quantizeColors(img image.Image, maxColors int) image.Image {
	// Implementation of color quantization
	// This is a placeholder - you'd want to implement a proper color quantization algorithm
	return img
}

func removeTransparency(img image.Image) image.Image {
	bounds := img.Bounds()
	rgb := image.NewRGBA(bounds)

	bgColor := color.White // Default background color
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a == 0 {
				rgb.Set(x, y, bgColor)
			} else {
				rgb.Set(x, y, color.RGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: 255,
				})
			}
		}
	}
	return rgb
}

func getConverter(format string) Converter {
	conv, _ := Factory(format)
	return conv
}
