package optimize

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/dendianugerah/reubah/pkg/errors"
)

type OptimizeOptions struct {
	Quality       int  // 1-100 for JPEG/WebP
	Progressive   bool // For JPEG
	Compression   int  // 0-9 for PNG
	StripMetadata bool // Remove EXIF and other metadata
	Interlace     bool // For PNG
	AutoQuality   bool // Automatically determine quality based on image content
	LossyPNG      bool // Allow lossy compression for PNG (smaller size)
}

// QualityLevel represents predefined quality settings
type QualityLevel string

const (
	QualityLow      QualityLevel = "low"      // Prioritize size over quality
	QualityMedium   QualityLevel = "medium"   // Balanced
	QualityHigh     QualityLevel = "high"     // Prioritize quality over size
	QualityLossless QualityLevel = "lossless" // No quality loss
)

// DefaultOptions returns recommended optimization options per format
func DefaultOptions(format string) OptimizeOptions {
	switch format {
	case "jpeg", "jpg":
		return OptimizeOptions{
			Quality:       85,
			Progressive:   true,
			StripMetadata: true,
			AutoQuality:   true,
		}
	case "png":
		return OptimizeOptions{
			Compression:   9,
			StripMetadata: true,
			Interlace:     false,
			LossyPNG:      false,
		}
	case "webp":
		return OptimizeOptions{
			Quality:       85,
			StripMetadata: true,
			AutoQuality:   true,
		}
	default:
		return OptimizeOptions{
			Quality:       85,
			Progressive:   false,
			Compression:   5,
			StripMetadata: false,
		}
	}
}

// GetOptionsForQuality returns optimization options for a specific quality level
func GetOptionsForQuality(format string, level QualityLevel) OptimizeOptions {
	opts := DefaultOptions(format)

	switch level {
	case QualityLow:
		opts.Quality = 60
		opts.LossyPNG = true
		opts.StripMetadata = true
	case QualityMedium:
		opts.Quality = 75
		opts.LossyPNG = false
		opts.StripMetadata = true
	case QualityHigh:
		opts.Quality = 90
		opts.LossyPNG = false
		opts.StripMetadata = false
	case QualityLossless:
		opts.Quality = 100
		opts.LossyPNG = false
		opts.Compression = 9
		opts.StripMetadata = false
	}

	return opts
}

// Optimize writes the image to the writer with optimization options
func Optimize(w io.Writer, img image.Image, format string, opts OptimizeOptions) error {
	if opts.AutoQuality {
		opts = autoAdjustQuality(img, opts)
	}

	switch format {
	case "jpeg", "jpg":
		return optimizeJPEG(w, img, opts)
	case "png":
		return optimizePNG(w, img, opts)
	default:
		return errors.New(errors.ErrInvalidFormat, "unsupported format for optimization", nil)
	}
}

func optimizeJPEG(w io.Writer, img image.Image, opts OptimizeOptions) error {
	options := &jpeg.Options{
		Quality: opts.Quality,
	}
	if err := jpeg.Encode(w, img, options); err != nil {
		return errors.New(errors.ErrProcessingFailed, "failed to optimize JPEG", err)
	}
	return nil
}

func optimizePNG(w io.Writer, img image.Image, opts OptimizeOptions) error {
	encoder := png.Encoder{
		CompressionLevel: png.CompressionLevel(opts.Compression),
	}
	if err := encoder.Encode(w, img); err != nil {
		return errors.New(errors.ErrProcessingFailed, "failed to optimize PNG", err)
	}
	return nil
}

// autoAdjustQuality analyzes image content and adjusts quality settings
func autoAdjustQuality(img image.Image, opts OptimizeOptions) OptimizeOptions {
	// Calculate image complexity
	complexity := calculateImageComplexity(img)
	
	// Adjust quality based on complexity
	if complexity < 0.3 {
		opts.Quality = 70 // Simple images can use lower quality
	} else if complexity < 0.6 {
		opts.Quality = 80 // Medium complexity
	} else {
		opts.Quality = 90 // Complex images need higher quality
	}

	return opts
}

// calculateImageComplexity returns a value between 0 and 1
// indicating image complexity based on color variance
func calculateImageComplexity(img image.Image) float64 {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	if width == 0 || height == 0 {
		return 0
	}

	// Sample points for efficiency
	sampleSize := 100
	variance := 0.0
	lastColor := img.At(0, 0)

	for i := 0; i < sampleSize; i++ {
		x := (i * width) / sampleSize
		for j := 0; j < sampleSize; j++ {
			y := (j * height) / sampleSize
			currentColor := img.At(x, y)
			
			// Calculate color difference
			r1, g1, b1, _ := lastColor.RGBA()
			r2, g2, b2, _ := currentColor.RGBA()
			
			diff := abs(int(r1)-int(r2)) + abs(int(g1)-int(g2)) + abs(int(b1)-int(b2))
			variance += float64(diff) / float64(0xffff*3)
			
			lastColor = currentColor
		}
	}

	// Normalize variance
	return variance / float64(sampleSize*sampleSize)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
} 