type Edge int

const (
	NONE Edge = iota
	DIM
	COLOR
)

type State struct {
	// Geometry & Sampling
	width, height int // in rune
	keepRatio     bool

	// Lanczos
	a int

	// Tone & Color
	monochrome bool
	//   - choose monochrome colors
	//   - brightness
	//   - contrast
	//   - gamma correction
	//   - invert color
	//   - invert luminance
	//   - color palette
	//   - palette size limit
	//   - quantization algorithm

	// Dithering
	dither bool
	//   - dithering algorithm
	ditherSize int
	//   - dither intensity
	//   - dither color
	//   - luma-only dithering
	//
	// Edges
	edgeType Edge

	//   - no edge
	//   - edge detection algorithm
	edgeThreshold float64
	//   - edge thickness

	edgeColor int
	edgeDim float64
	edgePalette string[]

	// Character Systems
	asciiPalette string[]
	
	// Braille
	braille bool
	//   - braille threshold
	//   - braille dot weighting
	//
	// Output
	//   - background character
	//   - transparent background
	//   - character spacing
}
