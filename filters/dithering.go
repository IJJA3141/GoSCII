package filters

import (
	"errors"
	"math"
)

// m generates an n-bit Bayer threshold matrix used for ordered dithering.
//
// The returned matrix has dimensions 2^n × 2^n and contains values in the range [0, 255].
//
// This algorithm follows the method described by Bisqwit:
// https://bisqwit.iki.fi/story/howto/dither/jy/

// Parameters:
//   - n: bit depth controlling the size of the Bayer matrix (matrix size = 2^n)
//
// Returns:
//   - A 2D slice of uint8 values representing the normalized Bayer threshold map.
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

// BayerDithering applies ordered Bayer dithering to a grayscale image.
//
// Each pixel in the output image is compared against a threshold from a
// Bayer matrix of size 2^n × 2^n. If the pixel intensity exceeds the threshold,
// it is set to 255 (white); otherwise, it remains 0 (black).
//
// Parameters:
//   - img: the input GrayscalePlane to be dithered
//   - n: bit depth controlling the Bayer matrix size (matrix size = 2^n)
//
// Returns:
//   - A new GrayscalePlane containing the dithered image
//   - An error if the bit depth n is less than 1
//
// The computation is parallelized across rows for improved performance.
func BayerDithering(img *GrayScalePlane, n int) (*GrayScalePlane, error) {
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

	return out, nil
}
