package filters

import (
	"errors"
	"math"
)

// not sure if this is useful

// // ResizablePlane represents an data plane that can be resized using
// // Lanczos resampling.
// //
// // The type parameter T is the concrete plane type returned by the resize
// // operation. Implementations must return the same concrete type as the receiver.
// type ResizablePlane[T any] interface {
//
// 	// LanczosResize resizes the given plane to the specified width and height using
// 	// Lanczos resampling with window size a.
// 	//
// 	// The resize is performed in two separable passes:
// 	//  1. Horizontal resampling
// 	//  2. Vertical resampling
// 	//
// 	// The parameter a controls the size of the Lanczos window. Typical values are:
// 	//   - a = 2 (Lanczos-2): faster, slightly softer
// 	//   - a = 3 (Lanczos-3): higher quality, more expensive
// 	//
// 	// An error is returned if the target dimensions or window size are invalid.
// 	LanczosResize(width, height, a int) (T, error)
// }
//
// // define supported types to the compiler
// var _ ResizablePlane[*RGBAPlane] = (*RGBAPlane)(nil)
// var _ ResizablePlane[*GrayScalePlane] = (*GrayScalePlane)(nil)
// var _ ResizablePlane[*EdgePlane] = (*EdgePlane)(nil)

func l(x float64, a int) float64 {
	_a := float64(a)

	if x == 0 {
		return 1
	}

	if -_a <= x && x < _a {
		return (_a * math.Sin(math.Pi*x) * math.Sin(math.Pi*x/_a)) / (math.Pi * math.Pi * x * x)
	}

	return 0
}

func coeffs(ratio float64, a, dimension int) []float64 {
	kern := make([]float64, dimension*2*a)

	split(dimension, func(_start, _end int) {
		for j := _start; j < _end && j < dimension; j++ {
			x := (float64(j)+0.5)*ratio - 0.5
			sum := 0.

			for i := range 2 * a {
				xd := math.Floor(x) + float64(i) - (float64(a) - 1.)
				c := l(x-xd, a)
				sum += c

				kern[j*a*2+i] = c
			}

			for i := range 2 * a {
				kern[j*a*2+i] /= sum
			}
		}
	}).Wait()

	return kern
}

// LanczosResize resizes the given RGBAPlane to the specified width and height using
// Lanczos resampling with window size a.
//
// The resize is performed in two separable passes:
//  1. Horizontal resampling
//  2. Vertical resampling
//
// The parameter a controls the size of the Lanczos window. Typical values are:
//   - a = 2 (Lanczos-2): faster, slightly softer
//   - a = 3 (Lanczos-3): higher quality, more expensive
//
// An error is returned if the target dimensions or window size are invalid.
func (img *RGBAPlane) LanczosResize(width, height, a int) (*RGBAPlane, error) {

	// Input sanitisation
	if width < 0 {
		return nil, errors.New("width must be positive")
	}

	if height < 0 {
		return nil, errors.New("height must be positive")
	}

	if a < 0 {
		return nil, errors.New("Lanczos window size must be positive")
	}

	// Horizontal resampling
	tmp := NewRGBAPlane(width, img.Height)
	ratio := float64(img.Width) / float64(width)
	kern := coeffs(ratio, a, width)

	split(img.Height, func(_start, _end int) {
		for srcY := _start; srcY < _end && srcY < img.Height; srcY++ {
			for x := range width {

				relX := (float64(x)+0.5)*ratio - 0.5
				var R, G, B, A float64

				for i := range 2 * a {
					srcX := int(math.Floor(relX) + float64(i) - (float64(a) - 1.))
					srcX = clamp(srcX, 0, img.Width-1)

					index := srcY*img.Stride + srcX*4
					coeff := kern[x*a*2+i] // x === j

					R += float64(img.RGBA[index]) * coeff
					G += float64(img.RGBA[index+1]) * coeff
					B += float64(img.RGBA[index+2]) * coeff
					A += float64(img.RGBA[index+3]) * coeff
				}

				index := srcY*tmp.Stride + x*4
				tmp.RGBA[index] = clamp(R, 0, 255)
				tmp.RGBA[index+1] = clamp(G, 0, 255)
				tmp.RGBA[index+2] = clamp(B, 0, 255)
				tmp.RGBA[index+3] = clamp(A, 0, 255)
			}
		}
	}).Wait()

	// Vertical resampling
	out := NewRGBAPlane(width, height)
	ratio = float64(img.Height) / float64(height)
	kern = coeffs(ratio, a, height)

	split(width, func(_start, _end int) {
		for srcX := _start; srcX < _end && srcX < width; srcX++ {
			for y := range height {

				relY := (float64(y)+0.5)*ratio - 0.5
				var R, G, B, A float64

				for i := range 2 * a {
					srcY := int(math.Floor(relY) + float64(i) - (float64(a) - 1.))
					srcY = clamp(srcY, 0, tmp.Height-1)

					index := srcY*tmp.Stride + srcX*4
					coeff := kern[y*a*2+i] // y === j

					R += float64(tmp.RGBA[index]) * coeff
					G += float64(tmp.RGBA[index+1]) * coeff
					B += float64(tmp.RGBA[index+2]) * coeff
					A += float64(tmp.RGBA[index+3]) * coeff
				}

				index := y*tmp.Stride + srcX*4
				out.RGBA[index] = clamp(R, 0, 255)
				out.RGBA[index+1] = clamp(G, 0, 255)
				out.RGBA[index+2] = clamp(B, 0, 255)
				out.RGBA[index+3] = clamp(A, 0, 255)
			}
		}
	}).Wait()

	return out, nil
}

// LanczosResize resizes the given GrayScalePlane to the specified width and height using
// Lanczos resampling with window size a.
//
// The resize is performed in two separable passes:
//  1. Horizontal resampling
//  2. Vertical resampling
//
// The parameter a controls the size of the Lanczos window. Typical values are:
//   - a = 2 (Lanczos-2): faster, slightly softer
//   - a = 3 (Lanczos-3): higher quality, more expensive
//
// An error is returned if the target dimensions or window size are invalid.
func (img *GrayScalePlane) LanczosResize(width, height, a int) (*GrayScalePlane, error) {

	// Input sanitisation
	if width < 0 {
		return nil, errors.New("") // TODO add real text
	}

	if height < 0 {
		return nil, errors.New("") // TODO add real text
	}

	if a < 0 {
		return nil, errors.New("") // TODO add real text
	}

	// Horizontal resampling
	tmp := NewGrayScalePlane(width, img.Height)
	ratio := float64(img.Width) / float64(width)
	kern := coeffs(ratio, a, width)

	split(img.Height, func(_start, _end int) {
		for srcY := _start; srcY < _end && srcY < img.Height; srcY++ {
			for x := range width {

				relX := (float64(x)+0.5)*ratio - 0.5
				var shade float64

				for i := range 2 * a {
					srcX := int(math.Floor(relX) + float64(i) - (float64(a) - 1.))
					srcX = clamp(srcX, 0, img.Width-1)

					shade += float64(img.Shades[srcY*img.Stride+srcX]) * kern[x*a*2+i] // x === j
				}

				tmp.Shades[srcY*tmp.Stride+x*4] = clamp(shade, 0, 255)
			}
		}
	}).Wait()

	// Vertical resampling
	out := NewGrayScalePlane(width, height)
	ratio = float64(img.Height) / float64(height)
	kern = coeffs(ratio, a, height)

	split(width, func(_start, _end int) {
		for srcX := _start; srcX < _end && srcX < width; srcX++ {
			for y := range height {

				relY := (float64(y)+0.5)*ratio - 0.5
				var shade float64

				for i := range 2 * a {
					srcY := int(math.Floor(relY) + float64(i) - (float64(a) - 1.))
					srcY = clamp(srcY, 0, tmp.Height-1)

					shade += float64(tmp.Shades[srcY*tmp.Stride+srcX]) * kern[y*a*2+i] // y === j
				}

				out.Shades[y*tmp.Stride+srcX*4] = clamp(shade, 0, 255)
			}
		}
	}).Wait()

	return out, nil
}

func (img *EdgePlane) LanczosResize(width, height, a int) (*EdgePlane, error) {

	// Input sanitisation
	if width < 0 {
		return nil, errors.New("") // TODO add real text
	}

	if height < 0 {
		return nil, errors.New("") // TODO add real text
	}

	if a < 0 {
		return nil, errors.New("") // TODO add real text
	}

	// Horizontal resampling
	tmp := NewEdgePlane(width, img.Height)
	ratio := float64(img.Width) / float64(width)
	kern := coeffs(ratio, a, width)

	split(img.Height, func(_start, _end int) {
		for srcY := _start; srcY < _end && srcY < img.Height; srcY++ {
			for x := range width {

				relX := (float64(x)+0.5)*ratio - 0.5
				var gradiant, angle float64

				for i := range 2 * a {
					srcX := int(math.Floor(relX) + float64(i) - (float64(a) - 1.))
					srcX = clamp(srcX, 0, img.Width-1)

					index := srcY*img.Stride + srcX*2
					coeff := kern[x*a*2+i] // x === j

					gradiant += float64(img.Gradient[index]) * coeff
					angle += float64(img.Gradient[index+1]) * coeff
				}

				index := srcY*tmp.Stride + x*2
				tmp.Gradient[index] = gradiant
				tmp.Gradient[index+1] = clamp(angle, -float64(math.Pi), float64(math.Pi))

				// TODO those boundaries myght be bether
				// tmp.Pairs[index] = clamp(gradiant, ..., ...)
				// tmp.Pairs[index+1] = mod(angle, math.Pi)
			}
		}
	}).Wait()

	// Vertical resampling
	out := NewEdgePlane(width, height)
	ratio = float64(img.Height) / float64(height)
	kern = coeffs(ratio, a, height)

	split(width, func(_start, _end int) {
		for srcX := _start; srcX < _end && srcX < width; srcX++ {
			for y := range height {

				relY := (float64(y)+0.5)*ratio - 0.5
				var gradiant, angle float64

				for i := range 2 * a {
					srcY := int(math.Floor(relY) + float64(i) - (float64(a) - 1.))
					srcY = clamp(srcY, 0, tmp.Height-1)

					index := srcY*tmp.Stride + srcX*2
					coeff := kern[y*a*2+i] // y === j

					gradiant += float64(tmp.Gradient[index]) * coeff
					angle += float64(tmp.Gradient[index+1]) * coeff
				}

				index := y*tmp.Stride + srcX*4
				out.Gradient[index] = gradiant
				out.Gradient[index+1] = clamp(angle, -float64(math.Pi), float64(math.Pi))

				// TODO those boundaries might be better
				// out.Pairs[index] = clamp(Gradiant, ..., ...)
				// out.Pairs[index+1] = mod(angle, math.Pi)
			}
		}
	}).Wait()

	return out, nil
}
