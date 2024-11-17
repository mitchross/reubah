package processor

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/chai2010/webp"
	"github.com/dendianugerah/reubah/internal/processor/background"
	"github.com/dendianugerah/reubah/internal/processor/convert"
	"github.com/dendianugerah/reubah/internal/processor/optimize"
	"github.com/dendianugerah/reubah/internal/processor/resize"
	"github.com/disintegration/imaging"
)

type ProcessOptions struct {
	Width             int
	Height            int
	ResizeMode        resize.ResizeMode
	OutputFormat      string
	Quality           int // 1-100
	Optimize          bool
	QualityLevel      optimize.QualityLevel
	CustomOptions     *optimize.OptimizeOptions
	ConversionOptions *convert.ConversionOptions
	RemoveBackground  bool
}

type ImageProcessor struct{}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func isValidFormat(format string) bool {
	validFormats := map[string]bool{
		"jpeg": true,
		"jpg":  true,
		"png":  true,
		"webp": true,
		"gif":  true,
		"bmp":  true,
	}
	return validFormats[format]
}

func (p *ImageProcessor) ProcessImageData(img image.Image, opts ProcessOptions) (*ProcessedImage, error) {
	// Validate format
	if opts.OutputFormat == "" {
		opts.OutputFormat = "jpeg"
	}

	if !isValidFormat(opts.OutputFormat) {
		return nil, fmt.Errorf("unsupported format: %s", opts.OutputFormat)
	}

	// Remove background if requested
	if opts.RemoveBackground {
		var err error
		img, err = background.RemoveBackground(img)
		if err != nil {
			return nil, fmt.Errorf("failed to remove background: %w", err)
		}
	}

	// Create conversion manager with options
	conversionOpts := opts.ConversionOptions
	if conversionOpts == nil {
		conversionOpts = convert.DefaultOptions(opts.OutputFormat)
	}
	conversionOpts.Quality = opts.Quality

	manager := convert.NewManager(conversionOpts)

	// Resize if needed
	if opts.Width > 0 || opts.Height > 0 {
		resizeOpts := resize.ResizeOptions{
			Width:  opts.Width,
			Height: opts.Height,
			Mode:   opts.ResizeMode,
			Filter: imaging.Lanczos,
		}
		var err error
		img, err = resize.Resize(img, resizeOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to resize image: %w", err)
		}
	}

	// Convert the image using the manager
	img, err := manager.Convert(img, opts.OutputFormat)
	if err != nil {
		return nil, fmt.Errorf("failed to convert image: %w", err)
	}

	return NewProcessedImage(img, opts.OutputFormat, opts.Quality), nil
}

type ProcessedImage struct {
	Image   image.Image
	Format  string
	Quality int
}

func NewProcessedImage(img image.Image, format string, quality int) *ProcessedImage {
	return &ProcessedImage{
		Image:   img,
		Format:  format,
		Quality: quality,
	}
}

func (pi *ProcessedImage) Write(w io.Writer) error {
	switch pi.Format {
	case "jpeg", "jpg":
		conversionOpts := convert.DefaultOptions("jpeg")
		conversionOpts.Quality = pi.Quality
		manager := convert.NewManager(conversionOpts)

		converted, err := manager.Convert(pi.Image, "jpeg")
		if err != nil {
			return fmt.Errorf("failed to convert to jpeg: %w", err)
		}

		return jpeg.Encode(w, converted, &jpeg.Options{Quality: pi.Quality})
	case "png":
		conversionOpts := convert.DefaultOptions("png")
		conversionOpts.Quality = pi.Quality
		manager := convert.NewManager(conversionOpts)

		converted, err := manager.Convert(pi.Image, "png")
		if err != nil {
			return fmt.Errorf("failed to convert to png: %w", err)
		}

		encoder := &png.Encoder{
			CompressionLevel: png.CompressionLevel((pi.Quality * 9) / 100),
		}
		return encoder.Encode(w, converted)
	case "webp":
		conversionOpts := convert.DefaultOptions("webp")
		conversionOpts.Quality = pi.Quality
		manager := convert.NewManager(conversionOpts)

		converted, err := manager.Convert(pi.Image, "webp")
		if err != nil {
			return fmt.Errorf("failed to convert to webp: %w", err)
		}

		var buf bytes.Buffer
		if err := webp.Encode(&buf, converted, &webp.Options{
			Lossless: pi.Quality == 100,
			Quality:  float32(pi.Quality),
		}); err != nil {
			return fmt.Errorf("failed to encode webp: %w", err)
		}
		_, err = w.Write(buf.Bytes())
		return err
	case "gif":
		conversionOpts := convert.DefaultOptions("gif")
		conversionOpts.Quality = pi.Quality
		manager := convert.NewManager(conversionOpts)

		converted, err := manager.Convert(pi.Image, "gif")
		if err != nil {
			return fmt.Errorf("failed to convert to gif: %w", err)
		}

		return gif.Encode(w, converted, &gif.Options{
			NumColors: (pi.Quality * 256) / 100,
		})
	default:
		return fmt.Errorf("unsupported format for writing: %s", pi.Format)
	}
}
