package filters

// The use of float64 for pixel storage is intentional: it allows image data
// to remain in a consistent, high precision representation while passing
// through multiple filters and transformations, thereby reducing cumulative
// rounding and conversion errors.

// RGBAPlane represents a two-dimensional RGBA image stored as a flat slice
// of float64 values.
//
// Pixels are stored in row major order, with four consecutive values per pixel
// representing the red, green, blue, and alpha channels respectively.
//
// The pixel at coordinates (x, y) starts at index:
//
//	(y * Stride) + (x * 4)
//
// where Stride is the number of float64 values between vertically adjacent pixels.
type RGBAPlane struct {
	// RGBA holds the image's pixel data in R, G, B, A order.
	RGBA []float64
	// Width and Height define the dimensions of the image in pixels.
	Width, Height int

	// Stride is the number of float64 values between the start of two
	// vertically adjacent pixels. For tightly packed images, this is Width * 4.
	Stride int
}

// NewRGBAPlane allocates and returns a new RGBAPlane with the given dimensions.
//
// The returned plane is tightly packed, with a stride of width * 4, and all
// channel values initialized to zero.
func NewRGBAPlane(width, height int) *RGBAPlane {
	return &RGBAPlane{
		RGBA:   make([]float64, height*width*4),
		Width:  width,
		Height: height,
		Stride: width * 4,
	}
}

// GrayscalePlane represents a two-dimensional grayscale image stored as a flat
// slice of float64 values.
//
// Each element in Shades corresponds to a single pixel intensity.
// Pixels are stored in row major order.
//
// The pixel at coordinates (x, y) starts at index:
//
//	(y * Stride) + (x)
//
// where Stride is the number of float64 values between vertically adjacent pixels.
type GrayScalePlane struct {
	// Shades holds the grayscale intensity values for each pixel.
	Shades []float64

	// Width and Height define the dimensions of the image in pixels.
	Width, Height int

	// Stride is the number of float64 values between vertically adjacent pixels.
	Stride int
}

// NewGrayscalePlane allocates and returns a new GrayscalePlane with the given dimensions.
//
// The returned plane is tightly packed, with a stride equal to the image width,
// and all values initialized to zero./
func NewGrayScalePlane(width, height int) *GrayScalePlane {
	return &GrayScalePlane{
		Shades: make([]float64, width*height), Width: width,
		Height: height,
		Stride: width,
	}
}

// AsciiPlane represents a two-dimensional image where each pixel is encoded
// as a single Unicode character.
//
// This type is typically used for ASCII art or text based visualizations.
// Characters are stored in row major order.
//
// The character at coordinates (x, y) starts at index:
//
//	(y * Stride) + (x)
//
// where Stride is the number of runes between vertically adjacent pixels.
type AsciiPlane struct {
	// Chars holds the rune associated with each pixel.
	Chars []rune

	// Width and Height define the dimensions of the image in characters.
	Width, Height int

	// Stride is the number of runes between vertically adjacent rows.
	Stride int
}

// NewAsciiPlane allocates and returns a new AsciiPlane with the given dimensions.
//
// The returned plane is tightly packed, with a stride equal to the image width.
func NewAsciiPlane(width, height int) *AsciiPlane {
	return &AsciiPlane{
		Chars:  make([]rune, width*height),
		Height: height,
		Width:  width,
		Stride: width,
	}
}

// AsciiColorPlane represents a two-dimensional image where each pixel is encoded
// as a colored ASCII element.
//
// Each entry in Chars typically represents a character along with its color
// information (for example, ANSI escape sequences).
//
// Elements are stored in row major order.
//
// The element at coordinates (x, y) starts at index:
//
//	(y * Stride) + (x)
//
// where Stride is the number of runes between vertically adjacent pixels.
type AsciiColorPlane struct {
	// Chars holds the string representation of each pixel, including color data.
	Chars []string

	// Width and Height define the dimensions of the image in characters.
	Width, Height int

	// Stride is the number of elements between vertically adjacent rows.
	Stride int
}

// NewAsciiColorPlane allocates and returns a new AsciiColorPlane with the given dimensions.
//
// The returned plane is tightly packed, with a stride equal to the image width.
func NewAsciiColorPlane(width, height int) *AsciiColorPlane {
	return &AsciiColorPlane{
		Chars:  make([]string, width*height),
		Width:  width,
		Height: height,
		Stride: width,
	}
}

// EdgePlane represents a two-dimensional image of edge information,
// where each pixel stores a (gradient magnitude, gradient angle) pair.
//
// The gradient magnitude describes the strength of the detected edge,
// while the gradient angle represents the edge orientation in radians
// and is constrained to the range [0, 2π].
//
// Data is stored as a flat slice of float64 values in row major order,
// with two consecutive values per pixel:
//
//	[ magnitude, angle ]
//
// The pair corresponding to the pixel at coordinates (x, y) starts at index:
//
//	(y * Stride) + (x * 2)
//
// where Stride is the number of float64 values between vertically adjacent pixels.
type EdgePlane struct {
	// Gradients holds the gradient magnitude and angle for each pixel,
	// stored as consecutive float64 values: [magnitude, angle].
	Gradient []float64

	// Width and Height define the dimensions of the image in pixels.
	Width, Height int

	// Stride is the number of float64 values between vertically adjacent pixels.
	// For tightly packed images, this is Width * 2.
	Stride int
}

// NewEdgePlane allocates and returns a new EdgePlane with the given dimensions.
//
// The returned plane is tightly packed, with a stride of width * 2, and all
// values initialized to zero.
func NewEdgePlane(width, height int) *EdgePlane {
	return &EdgePlane{
		Gradient: make([]float64, height*width*2),
		Width:    width,
		Height:   height,
		Stride:   width * 2,
	}
}
