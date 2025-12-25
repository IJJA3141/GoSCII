package filters

import (
	"math"
)

// SobelEdgeDetection computes the Sobel edge detection of the grayscale image and returns
// an EdgePlane containing the gradient magnitude and orientation for each pixel.
//
// Pixels near the image boundary are handled by ignoring samples that fall
// outside the image bounds.
//
// The computation is parallelized across image rows.
func (img *GrayScalePlane) SobelEdgeDetection() *EdgePlane {
	// For each pixel, the Sobel operator estimates the horizontal (Gx) and vertical
	// (Gy) intensity gradients using a 3×3 convolution kernel. The resulting edge
	//
	// The gradient magnitude is computed as:
	//
	//	sqrt(Gx² + Gy²)
	//
	// and the gradient angle is computed using atan2(Gy, Gx), yielding an angle
	// in radians in the range [0, 2π].
	//
	// https://en.wikipedia.org/wiki/Sobel_operator

	// Allocate the output edge plane with matching dimensions.
	out := NewEdgePlane(img.Width, img.Height)

	// Sobel kernel for horizontal gradient (Gx).
	// The vertical gradient (Gy) is obtained by transposing the kernel.
	kern := []float64{
		-1, 0, 1,
		-2, 0, 2,
		-1, 0, 1,
	}

	// Split the work across row ranges to enable parallel processing.
	split(img.Height, func(start, end int) {
		for y := start; y < end && y < img.Height; y++ {
			for x := range img.Width {
				var sumX, sumY float64

				// Apply the 3×3 Sobel kernel centered at (x, y).
				for j := range 3 {
					srcY := y - 1 + j
					for i := range 3 {
						srcX := x - 1 + i

						if 0 <= srcY && srcY < img.Height && 0 <= srcX && srcX < img.Width {
							// Fetch the source pixel intensity.
							pix := float64(img.Shades[srcY*img.Stride+srcX])

							// Accumulate horizontal and vertical gradients.
							sumX += pix * kern[j*3+i] // K
							sumY += pix * kern[j+i*3] // K_T
						}
					}
				}

				// Store the gradient magnitude and angle for this pixel.
				index := y*out.Stride + 2*x
				out.Gradient[index] = math.Sqrt(sumX*sumX + sumY*sumY)
				out.Gradient[index+1] = math.Mod(math.Atan2(sumY, sumX)+2*math.Pi, 2*math.Pi)
			}
		}
	}).Wait()

	return out
}
