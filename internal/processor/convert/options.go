package convert

// ConversionOptions contains all possible options for image conversion
type ConversionOptions struct {
	// Basic options
	Quality int
	Format  string

	// Advanced options
	RemoveMetadata       bool
	PreserveColorProfile bool
	Interlace            bool
	Progressive          bool // For JPEG
	OptimizeSize         bool // Enable additional size optimization

	// Format-specific options
	WebP *WebPOptions
	PNG  *PNGOptions
	GIF  *GIFOptions
}

type WebPOptions struct {
	Lossless bool
	Exact    bool // Preserve RGB values in transparent area
}

type PNGOptions struct {
	CompressLevel  int  // 0-9
	QuantizeColors bool // Reduce number of colors
	MaxColors      int  // Maximum number of colors when quantizing
}

type GIFOptions struct {
	NumColors      int  // 2-256
	Dither         bool // Enable dithering
	PreserveFrames bool // Preserve animation frames
}

// DefaultOptions returns default conversion options
func DefaultOptions(format string) *ConversionOptions {
	opts := &ConversionOptions{
		Quality:              85,
		Format:               format,
		RemoveMetadata:       true,
		PreserveColorProfile: false,
		OptimizeSize:         true,
	}

	switch format {
	case FormatWebP:
		opts.WebP = &WebPOptions{
			Lossless: false,
			Exact:    false,
		}
	case FormatPNG:
		opts.PNG = &PNGOptions{
			CompressLevel:  6,
			QuantizeColors: false,
			MaxColors:      256,
		}
	case FormatGIF:
		opts.GIF = &GIFOptions{
			NumColors:      256,
			Dither:         true,
			PreserveFrames: true,
		}
	}

	return opts
}
