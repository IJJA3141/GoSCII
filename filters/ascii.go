package filters

import (
	_ "errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Stampable interface {
	Width_() int
	Height_() int
	Buffer() []string
}

type Ascii interface {
	Get(x, y, width, height int) []string
	Width_() int
	Height_() int
}

// Colorize converts an AsciiPlane into an AsciiColorPlane using an RGBAPlane
// as the source of color information.
//
// Each ASCII character is prefixed with an ANSI 24-bit color escape sequence
// based on the corresponding RGB pixel in the colors plane. The alpha channel
// is ignored.
//
// If consecutive pixels share the same RGB color, the ANSI escape sequence
// is reused to reduce redundant escape codes.
//
// Parameters:
//   - colors: the RGBAPlane providing per-pixel color information
//
// Returns:
//   - A new AsciiColorPlane where each character is colorized according to
//     the input colors plane
//   - An error if the dimensions of the ASCII plane and colors plane differ
func (ascii *AsciiPlane) Colorize(colors *RGBAPlane) (*AsciiColorPlane, error) {
	if ascii.Height != colors.Height || ascii.Width != colors.Width {
		// return nil, errors.New("Colorize: dimensions of ASCII plane and color plane do not match")
		return nil, fmt.Errorf("Colorize: dimensions of ASCII plane and color plane do not match\nascii.Height %d != colors.Height %d\nascii.Width %d != colors.Width %d", ascii.Height, colors.Height, ascii.Width, colors.Width)
	}

	out := NewAsciiColorPlane(ascii.Width, ascii.Height)

	split(out.Height, func(_start, _end int) {
		// var prevR, prevG, prevB uint8 = uint8(colors.RGBA[_start*colors.Stride]) - 1, 0, 0

		for y := _start; y < _end && y < out.Height; y++ {
			for x := range out.Width {
				index := y*colors.Stride + x*4
				r := uint8(colors.RGBA[index])
				g := uint8(colors.RGBA[index+1])
				b := uint8(colors.RGBA[index+2])
				// a:=uint8(_colors.RGBA[index+3]) // might take the alpha in consideration

				// if r == prevR && g == prevG && b == prevB {
				// 	out.Chars[y*out.Stride+x] = string(ascii.Chars[y*ascii.Stride+x])
				// } else {
				out.Chars[y*out.Stride+x] =
					"\x1B[38;2;" + strconv.FormatUint(uint64(r), 10) +
						";" + strconv.FormatUint(uint64(g), 10) +
						";" + strconv.FormatUint(uint64(b), 10) + "m" + string(ascii.Chars[y*ascii.Stride+x])

				// prevR = r
				// prevG = g
				// prevB = b
				// }
			}
		}
	}).Wait()

	return out, nil
}

// Ascii converts a GrayScalePlane into an AsciiPlane using a specified character palette.
//
// Each pixel's intensity (0–255) is mapped linearly to a character in the palette.
// Darker pixels correspond to characters earlier in the palette; brighter pixels
// correspond to characters later in the palette.
//
// Parameters:
//   - palette: a slice of runes representing the ASCII characters to use for mapping
//
// Returns:
//   - A new AsciiPlane where each pixel is replaced by a corresponding character
func (img *GrayScalePlane) Ascii(palette []rune) *AsciiPlane {
	out := NewAsciiPlane(img.Width, img.Height)

	split(out.Height, func(_start, _end int) {
		for y := _start; y < _end && y < out.Height; y++ {
			for x := range out.Width {
				index := y*out.Stride + x
				lum := img.Shades[index]                              // 0 .. 255
				bucket := int((lum / 255.) * float64(len(palette)-1)) // index in palette
				out.Chars[index] = palette[bucket]
			}
		}
	}).Wait()

	return out
}

const π2 = math.Pi * 2

func (img *EdgePlane) Ascii(threshold float64, palette []rune) *AsciiPlane {
	out := NewAsciiPlane(img.Width, img.Height)

	split(out.Height, func(_start, _end int) {
		for y := _start; y < _end && y < out.Height; y++ {
			for x := range out.Width {
				index := y*img.Stride + x*2
				magnitude := img.Gradient[index]
				angle := img.Gradient[index+1] // 0 .. 2pi

				if magnitude < threshold {
					out.Chars[y*out.Stride+x] = ' '
				} else {
					bucket := int(angle / (π2) * float64(len(palette)-1)) // index in palette
					out.Chars[y*out.Stride+x] = palette[bucket]
				}
			}
		}
	}).Wait()

	return out
}

var dotMatrix = [][]uint16{
	{0x1, 0x8},
	{0x2, 0x10},
	{0x4, 0x20},
	{0x40, 0x80},
}

// Braille converts a GrayScalePlane into an AsciiPlane using Unicode Braille characters.
//
// Each Braille character represents a block of 2×4 pixels from the input image.
// Pixels with intensity greater than the specified threshold are considered "on"
// and mapped to the corresponding dot in the Braille character.
//
// Braille Unicode characters start at 0x2800. Each bit in a character represents
// a dot in the 2×4 block.
//
// Parameters:
//   - threshold: the intensity threshold for turning a dot "on" (0–255)
//
// Returns:
//   - A new AsciiPlane where each character is a Unicode Braille symbol representing
//     a 2×4 pixel block
//
// Notes:
//   - The output width is img.Width / 2, and the output height is img.Height / 4
//   - Pixels are grouped top-to-bottom, left-to-right in the 2×4 block
//   - The function is parallelized across rows for performance
func (img *GrayScalePlane) Braille(threshold float64) *AsciiPlane {
	out := NewAsciiPlane(img.Width/2, img.Height/4)

	split(out.Height, func(_start, _end int) {
		var char uint16

		for y := _start; y < _end && y < out.Height; y++ {
			for x := range out.Width {
				char = 0x2800 // block offset

				for j := range 4 {
					for i := range 2 {
						if img.Shades[(y*4+j)*img.Stride+(x*2+i)] >= threshold {
							char += dotMatrix[j][i]
						}
					}
				}

				out.Chars[y*out.Stride+x] = rune(char)
			}
		}
	}).Wait()

	return out
}

// idee de merde il faudrait pouvoir avoir des bords d'une couleur differancte.
// func CarveEdges(_img *AsciiPlane, _edges *EdgePlane, _palette []rune) *AsciiPlane {
// 	return &AsciiPlane{}
// }

func (this *AsciiPlane) Width_() int  { return this.Width }
func (this *AsciiPlane) Height_() int { return this.Height }
func (this *AsciiPlane) Buffer() []string {
	out := make([]string, this.Height)

	split(this.Height, func(_start, _end int) {
		for y := _start; y < _end && y < this.Height; y++ {
			out[y] = string(this.Chars[y*this.Stride : y*this.Stride+this.Width])
		}
	}).Wait()

	return out
}

func (img *AsciiColorPlane) Stamp() (int, int, [][]string) {
	out := make([][]string, img.Height)

	for y := range img.Height {
		out[y] = img.Chars[y*img.Stride : y*img.Stride+img.Width]
	}

	return img.Width, img.Height, out
}

func (this *AsciiColorPlane) Width_() int  { return this.Width }
func (this *AsciiColorPlane) Height_() int { return this.Height }
func (this *AsciiColorPlane) Buffer() []string {
	out := make([]string, this.Height)

	split(this.Height, func(_start, _end int) {
		for y := _start; y < _end && y < this.Height; y++ {
			out[y] = strings.Join(this.Chars[y*this.Stride:y*this.Stride+this.Width], "")
		}
	}).Wait()

	return out
}

func (this *AsciiPlane) Get(x, y, width, height int) []string {
	out := make([]string, height)

	for i := range height {
		index := (y+i)*this.Stride + x
		out[i] = string(this.Chars[index : index+width])
	}

	return out
}

func (this *AsciiColorPlane) Get(x, y, width, height int) []string {
	out := make([]string, height)

	for i := range height {
		index := (y+i)*this.Stride + x
		out[i] = strings.Join((this.Chars[index : index+width]), "")
	}

	return out
}
