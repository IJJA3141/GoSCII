package filters

import (
	"fmt"
	"image"
	"math"
	"runtime"
	"sync"
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

	for j := range _Δ {
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

	return kern
}

func Resize(_in *image.NRGBA, _width, _height, _nWorkers, _a int) *image.NRGBA {
	Δx := float64(_in.Rect.Dx()) / float64(_width)
	Δy := float64(_in.Rect.Dy()) / float64(_height)

	tmp := image.NewNRGBA(image.Rect(0, 0, _width, _in.Rect.Dy()))

	kern := coeffs(Δx, _a, _width)
	for j := range _width {
		for i := range 2 * _a {
			fmt.Print(kern[j*_a*2+i], ", ")
		}
		fmt.Println()
	}

	var wg sync.WaitGroup
	wg.Add(_in.Rect.Dy())

	runtime.GOMAXPROCS(0)

	for y := range _in.Rect.Dy() {

		go func() {
			defer wg.Done()

			for j := range _width {

				x := (float64(j)+0.5)*Δx - 0.5

				var R, G, B, A float64

				for a := range 2 * _a {
					i := math.Floor(x) + float64(a) - (float64(_a) - 1.)
					c := kern[j*_a*2+a]

					ii := int(i)

					if ii < 0 {
						ii = 0
					} else if _in.Rect.Dx() <= ii {
						ii = _in.Rect.Dx() - 1
					}

					index := y*_in.Stride + ii*4

					R += float64(_in.Pix[index]) * c
					G += float64(_in.Pix[index+1]) * c
					B += float64(_in.Pix[index+2]) * c
					A += float64(_in.Pix[index+3]) * c
				}

				index := y*tmp.Stride + j*4
				tmp.Pix[index] = clamp(R)
				tmp.Pix[index+1] = clamp(G)
				tmp.Pix[index+2] = clamp(B)
				tmp.Pix[index+3] = clamp(A)
			}
		}()
	}

	wg.Wait()

	out := image.NewNRGBA(image.Rect(0, 0, _width, _height))

	kern = coeffs(Δy, _a, _height)

	for j := range _height {
		for i := range 2 * _a {
			fmt.Print(kern[j*_a*2+i], ", ")
		}
		fmt.Println()
	}

	wg.Add(_width)

	for x := range _width {

		go func() {
			wg.Done()

			for j := range _height {
				y := (float64(j)+0.5)*Δy - 0.5

				var R, G, B, A float64

				for a := range 2 * _a {
					i := math.Floor(y) + float64(a) - (float64(_a) - 1.)
					c := kern[j*_a*2+a]

					ii := int(i)

					if ii < 0 {
						ii = 0
					} else if tmp.Rect.Dy() <= ii {
						ii = tmp.Rect.Dy() - 1
					}

					index := ii*tmp.Stride + x*4

					R += float64(tmp.Pix[index]) * c
					G += float64(tmp.Pix[index+1]) * c
					B += float64(tmp.Pix[index+2]) * c
					A += float64(tmp.Pix[index+3]) * c
				}

				index := j*tmp.Stride + x*4
				out.Pix[index] = clamp(R)
				out.Pix[index+1] = clamp(G)
				out.Pix[index+2] = clamp(B)
				out.Pix[index+3] = clamp(A)
			}
		}()
	}

	wg.Wait()

	return out
}
