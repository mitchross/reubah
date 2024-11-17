package convert

// Format constants to avoid string literals
const (
	FormatJPEG = "jpeg"
	FormatJPG  = "jpg"
	FormatPNG  = "png"
	FormatWebP = "webp"
	FormatGIF  = "gif"
	FormatBMP  = "bmp"
)

// DefaultQuality values for different formats
const (
	DefaultJPEGQuality = 85
	DefaultPNGQuality  = 6  // compression level
	DefaultWebPQuality = 85
	DefaultGIFColors   = 256
)

// Quality ranges
const (
	MinQuality = 1
	MaxQuality = 100
) 