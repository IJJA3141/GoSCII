package filters

import (
	"image"
	"math"
)

func l(_x float64, _a int) float64 {
	a := float64(_a)

	if _x == 0 {
		return 1
	}

	if -a <= _x && _x < a {
		return (a * math.Sin(math.Pi*_x) * math.Sin(math.Pi*_x/a)) / (math.Pi * math.Pi * _x * _x)
	}

	return 0
}

func clamp(_x float64) uint8 {
	if uint32(_x) < 256 {
		return uint8(_x)
	}

	if _x > 255 {
		return 255
	}

	return 0
}

func coeffs(_λ float64, _a, _Δ int) []float64 {
	kern := make([]float64, _Δ*2*_a)

	Split(_Δ, func(_start, _end int) {
		for j := _start; j < _end && j < _Δ; j++ {
			x := (float64(j)+0.5)*_λ - 0.5
			σ := 0.

			for i := range 2 * _a {
				xd := math.Floor(x) + float64(i) - (float64(_a) - 1.)
				c := l(x-xd, _a)
				σ += c

				kern[j*_a*2+i] = c
			}

			for i := range 2 * _a {
				kern[j*_a*2+i] /= σ
			}
		}
	}).Wait()

	return kern
}

func Resize(_in *image.NRGBA, _width, _height, _a int) *image.NRGBA {
	assert(_a > 0)

	tmp := image.NewNRGBA(image.Rect(0, 0, _width, _in.Rect.Dy()))
	Δx := float64(_in.Rect.Dx()) / float64(_width)
	kern := coeffs(Δx, _a, _width)

	Split(_in.Rect.Dy(), func(_start, _end int) {
		for y := _start; y < _end && y < _in.Rect.Dy(); y++ {
			for j := range _width {

				x := (float64(j)+0.5)*Δx - 0.5
				var R, G, B, A float64

				for a := range 2 * _a {
					xRel := int(math.Floor(x) + float64(a) - (float64(_a) - 1.))
					xRel = max(0, xRel)
					xRel = min(xRel, _in.Rect.Dx()-1)

					index := y*_in.Stride + xRel*4
					coeff := kern[j*_a*2+a]

					R += float64(_in.Pix[index]) * coeff
					G += float64(_in.Pix[index+1]) * coeff
					B += float64(_in.Pix[index+2]) * coeff
					A += float64(_in.Pix[index+3]) * coeff
				}

				index := y*tmp.Stride + j*4
				tmp.Pix[index] = clamp(R)
				tmp.Pix[index+1] = clamp(G)
				tmp.Pix[index+2] = clamp(B)
				tmp.Pix[index+3] = clamp(A)
			}
		}
	}).Wait()

	out := image.NewNRGBA(image.Rect(0, 0, _width, _height))
	Δy := float64(_in.Rect.Dy()) / float64(_height)
	kern = coeffs(Δy, _a, _height)

	Split(_width, func(_start, _end int) {
		for x := _start; x < _end && x < _width; x++ {
			for j := range _height {

				y := (float64(j)+0.5)*Δy - 0.5
				var R, G, B, A float64

				for a := range 2 * _a {
					i := int(math.Floor(y) + float64(a) - (float64(_a) - 1.))
					i = max(0, i)
					i = min(i, tmp.Rect.Dy()-1)

					index := i*tmp.Stride + x*4
					coeff := kern[j*_a*2+a]

					R += float64(tmp.Pix[index]) * coeff
					G += float64(tmp.Pix[index+1]) * coeff
					B += float64(tmp.Pix[index+2]) * coeff
					A += float64(tmp.Pix[index+3]) * coeff
				}

				index := j*tmp.Stride + x*4
				out.Pix[index] = clamp(R)
				out.Pix[index+1] = clamp(G)
				out.Pix[index+2] = clamp(B)
				out.Pix[index+3] = clamp(A)
			}
		}
	}).Wait()

	return out
}
