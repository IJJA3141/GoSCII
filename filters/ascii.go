package filters

import (
	"image"
)

var (
	a   = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	b   = " .:-=+*#%@"
	run = [2]rune{'⠓', '@'}
)

// width x height, 11x23 ->

type AsciiImage struct {
	// The rune at (x, y) starts at Pairs[y*Stride + x].
	Runes         []rune
	Height, Width int
}

type Palette struct {
	Runes string
	Ratio float64
}

// Ascii converts a grayscale image into ASCII art using the given palette.
// The palette should be ordered from the darkest to the brightest characters.
func Ascii(_in image.Gray, palette []string) AsciiImage {
	out := AsciiImage{
		Runes:  make([]rune, _in.Rect.Dy()*_in.Rect.Dx()),
		Height: _in.Rect.Dy(),
		Width:  _in.Rect.Dx(),
	}

	// r := float64(len(palette)-1) / 255

	for y := range out.Height {
		for x := range out.Width {
			out.Runes[y*out.Width+x] = 0 // palette[int(float64(_in.Pix[y*_in.Stride+x])*r)]
		}
	}

	return out
}

func Braille(_in image.Gray, _threshold uint8) AsciiImage {
	out := AsciiImage{
		Runes:  make([]rune, (_in.Rect.Dy()/1)*(_in.Rect.Dx()/1)),
		Height: _in.Rect.Dy() / 1,
		Width:  _in.Rect.Dx() / 1,
	}

	Split(out.Height, func(_start, _end int) {
		for y := _start; y < _end && y < out.Height; y++ {
			for x := range out.Width {
				var offset uint16 = 0

				// [OO]
				// [OO]
				// [OO]
				// [XX] <-
				// offset <<= 1
				if _in.Pix[(y+3)*_in.Stride+(x+1)] < _threshold {
					offset += 1
				}

				offset <<= 1
				if _in.Pix[(y+3)*_in.Stride+(x)] < _threshold {
					offset += 1
				}

				// [OX]
				// [OX] ^
				// [OX] |
				// [OO]
				offset <<= 1
				if _in.Pix[(y+2)*_in.Stride+(x+1)] < _threshold {
					offset += 1
				}

				offset <<= 1
				if _in.Pix[(y+1)*_in.Stride+(x+1)] < _threshold {
					offset += 1
				}

				offset <<= 1
				if _in.Pix[(y)*_in.Stride+(x+1)] < _threshold {
					offset += 1
				}

				// [XO]
				// [XO] ^
				// [XO] |
				// [OO]
				offset <<= 1
				if _in.Pix[(y+2)*_in.Stride+(x)] < _threshold {
					offset += 1
				}

				offset <<= 1
				if _in.Pix[(y+1)*_in.Stride+(x)] < _threshold {
					offset += 1
				}

				offset <<= 1
				if _in.Pix[(y)*_in.Stride+(x)] < _threshold {
					offset += 1
				}

				out.Runes[y*out.Width+x] = rune(0x2800 + offset)
			}
		}
	}).Wait()

	return out
}
