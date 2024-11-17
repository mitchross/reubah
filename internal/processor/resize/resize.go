package resize

import (
	"fmt"
	"image"
	"strings"

	"github.com/disintegration/imaging"
)

// ResizeMode defines how the image should be resized
type ResizeMode int

const (
	ModeAspectFit ResizeMode = iota // Maintain aspect ratio, fit within dimensions
	ModeFill                        // Fill the dimensions, crop if necessary
	ModeStretch                     // Stretch/squish to exactly match dimensions
)

// String representations of resize modes
const (
	ModeAspectFitStr = "fit"
	ModeFillStr      = "fill"
	ModeStretchStr   = "stretch"
)

// ParseResizeMode converts a string to ResizeMode
func ParseResizeMode(mode string) (ResizeMode, error) {
	switch strings.ToLower(mode) {
	case ModeAspectFitStr, "aspect", "aspectfit":
		return ModeAspectFit, nil
	case ModeFillStr, "cover":
		return ModeFill, nil
	case ModeStretchStr, "exact":
		return ModeStretch, nil
	default:
		return ModeAspectFit, fmt.Errorf("invalid resize mode: %s", mode)
	}
}

// ResizeOptions contains all options for image resizing
type ResizeOptions struct {
	Width  int
	Height int
	Mode   ResizeMode
	Filter imaging.ResampleFilter
}

// Resize resizes the image according to the specified options
func Resize(img image.Image, opts ResizeOptions) (image.Image, error) {
	if opts.Width == 0 && opts.Height == 0 {
		return img, nil // No resize needed
	}

	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	switch opts.Mode {
	case ModeAspectFit:
		return aspectFit(img, opts, origWidth, origHeight)
	case ModeFill:
		return fill(img, opts, origWidth, origHeight)
	case ModeStretch:
		return imaging.Resize(img, opts.Width, opts.Height, opts.Filter), nil
	default:
		return nil, fmt.Errorf("unsupported resize mode: %d", opts.Mode)
	}
}

func aspectFit(img image.Image, opts ResizeOptions, origWidth, origHeight int) (image.Image, error) {
	if opts.Width == 0 {
		// Calculate width to maintain aspect ratio
		opts.Width = int(float64(origWidth) * float64(opts.Height) / float64(origHeight))
	} else if opts.Height == 0 {
		// Calculate height to maintain aspect ratio
		opts.Height = int(float64(origHeight) * float64(opts.Width) / float64(origWidth))
	} else {
		// Both dimensions specified, maintain aspect ratio within bounds
		widthRatio := float64(opts.Width) / float64(origWidth)
		heightRatio := float64(opts.Height) / float64(origHeight)
		
		if widthRatio < heightRatio {
			opts.Height = int(float64(origHeight) * widthRatio)
		} else {
			opts.Width = int(float64(origWidth) * heightRatio)
		}
	}

	return imaging.Resize(img, opts.Width, opts.Height, opts.Filter), nil
}

func fill(img image.Image, opts ResizeOptions, origWidth, origHeight int) (image.Image, error) {
	// First resize to cover the target dimensions while maintaining aspect ratio
	widthRatio := float64(opts.Width) / float64(origWidth)
	heightRatio := float64(opts.Height) / float64(origHeight)
	
	ratio := widthRatio
	if heightRatio > widthRatio {
		ratio = heightRatio
	}
	
	resizedWidth := int(float64(origWidth) * ratio)
	resizedHeight := int(float64(origHeight) * ratio)
	
	resized := imaging.Resize(img, resizedWidth, resizedHeight, opts.Filter)
	
	// Then crop to exact dimensions
	return imaging.CropCenter(resized, opts.Width, opts.Height), nil
}
