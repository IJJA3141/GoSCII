package filters

type RGBAPlane struct {
	// RGBA holds the image's pixel data in R, G, B, A order.
	RGBA []float64
	// Width and Height define the dimensions of the image in pixels.
	Width, Height int

	// Stride is the number of float64 values between the start of two
	// vertically adjacent pixels. For tightly packed images, this is Width * 4.
	Stride int
}

func NewRGBAPlane(width, height int) RGBAPlane {
	return RGBAPlane{
		RGBA:   make([]float64, height*width*4),
		Width:  width,
		Height: height,
		Stride: width * 4,
	}
}

type GrayScalePlane struct {
	// Shades holds the grayscale intensity values for each pixel.
	Shades []float64

	// Width and Height define the dimensions of the image in pixels.
	Width, Height int

	// Stride is the number of float64 values between vertically adjacent pixels.
	Stride int
}

func NewGrayScalePlane(width, height int) GrayScalePlane {
	return GrayScalePlane{
		Shades: make([]float64, width*height), Width: width,
		Height: height,
		Stride: width,
	}
}

type AsciiPlane struct {
	// Chars holds the rune associated with each pixel.
	Chars []rune

	// Width and Height define the dimensions of the image in characters.
	Width, Height int

	// Stride is the number of runes between vertically adjacent rows.
	Stride int
}

func NewAsciiPlane(width, height int) AsciiPlane {
	return AsciiPlane{
		Chars:  make([]rune, width*height),
		Height: height,
		Width:  width,
		Stride: width,
	}
}

type AsciiColorPlane struct {
	// Chars holds the string representation of each pixel, including color data.
	Chars []string

	// Width and Height define the dimensions of the image in characters.
	Width, Height int

	// Stride is the number of elements between vertically adjacent rows.
	Stride int
}

func NewAsciiColorPlane(width, height int) AsciiColorPlane {
	return AsciiColorPlane{
		Chars:  make([]string, width*height),
		Width:  width,
		Height: height,
		Stride: width,
	}
}

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

func NewEdgePlane(width, height int) EdgePlane {
	return EdgePlane{
		Gradient: make([]float64, height*width*2),
		Width:    width,
		Height:   height,
		Stride:   width * 2,
	}
}
