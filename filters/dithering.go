package filters

import (
	"errors"
	"math"
)

func m(n int) [][]uint8 {
	dimension := uint(1 << n)
	normalisation := 1. / float64((dimension * dimension))
	out := make([][]uint8, dimension)

	for j := range dimension {
		out[j] = make([]uint8, dimension)

		for i := range dimension {
			mask := uint(n - 1)
			v := uint(0)

			for bit := 0; bit < 2*n; mask-- {
				v |= ((j >> mask) & 1) << bit
				bit++
				v |= (((i ^ j) >> mask) & 1) << bit
				bit++
			}

			// Normalize to 0-255
			out[j][i] = uint8(float64(v) * normalisation * 255.)
		}
	}

	return out
}

func (img *GrayScalePlane) BayerDithering(n int) (*GrayScalePlane, error) {
	if n < 1 {
		return nil, errors.New("BayerDithering: n must be >= 1")
	}

	// Allocate output image
	out := NewGrayScalePlane(img.Width, img.Height)

	// Generate the Bayer threshold map
	M := m(n)
	norm := float64(int(1) << n)

	split(img.Height, func(_start, _end int) {
		for y := _start; y < _end && y < img.Height; y++ {
			for x := range img.Width {

				lhs := uint8(img.Shades[y*img.Stride+x])
				rhs := M[int(math.Mod(float64(y), norm))][int(math.Mod(float64(x), norm))]

				if lhs > rhs {
					out.Shades[y*out.Stride+x] = 255
				}
			}
		}
	}).Wait()

	return &out, nil
}
