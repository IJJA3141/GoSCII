package filters

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

const Pi2 = math.Pi * math.Pi

func sinc(x float64) float64 {
	if x == 0 {
		return 1
	} else {
		return math.Sin(math.Pi*x) / (x * math.Pi)
	}
}

func lanczosKernel(x, a float64) float64 {
	if -a <= x && x < a {
		return sinc(x) * sinc(x/a)
	} else {
		return 0
	}
}

func l2(x, y, a float64) float64 {
	return lanczosKernel(math.Sqrt(x*x+y*y), a)
}

func clamp(x float64) uint8 {
	if uint32(x) < 256 {
		return uint8(x)
	}
	if x > 255 {
		return 255
	}
	return 0
}

func Resize2(_img *image.NRGBA, _size image.Rectangle, _a float64) *image.NRGBA {
	bounds := _img.Bounds()
	out := image.NewNRGBA(_size)

	xStep := float64(bounds.Dx()) / float64(_size.Dx())
	yStep := float64(bounds.Dy()) / float64(_size.Dy())

	for x := _size.Min.X; x < _size.Max.X; x++ {
		xSource := (float64(x)+0.5)*xStep - 0.5

		for y := _size.Min.Y; y < _size.Max.Y; y++ {
			ySource := (float64(y)+0.5)*yStep - 0.5

			var w, R, G, B, A float64

			for i := -_a + 1; i <= _a; i++ {
				for j := -_a + 1; j <= _a; j++ {

					c := l2(i-xSource+math.Floor(xSource), j-ySource+math.Floor(ySource), _a)

					pix := _img.NRGBAAt(int(math.Floor(xSource)-i), int(math.Floor(ySource)-j))

					w += c
					A += float64(pix.A) * c
					R += float64(pix.R) * c
					G += float64(pix.G) * c
					B += float64(pix.B) * c
				}
			}

			R /= w
			G /= w
			B /= w
			A /= w

			out.SetNRGBA(x, y, color.NRGBA{clamp(R), clamp(G), clamp(B), clamp(A)})
		}
	}

	return out
}

// if x axis -> _incr == 4
// if y axis -> _incr == stride
func stretch(_in []uint8, _startIn, _endInt int, _step float64, _out []uint8, _start, _end, _incr int, _a float64) {
	for xOut := _start; xOut < _end; xOut += _incr {
		xIn := (float64(xOut)/float64(_incr)+0.5)*_step - 0.5

		var σ, R, G, B, A float64

		for i := -_a + 1; i <= _a; i++ {
			l := lanczosKernel(i-xIn+math.Floor(xIn), _a)
			σ += l

			R += float64(_in[int(math.Floor(xIn)-i)]) * l
			G += float64(_in[int(math.Floor(xIn)-i)]+1) * l
			B += float64(_in[int(math.Floor(xIn)-i)]+2) * l
			A += float64(_in[int(math.Floor(xIn)-i)]+3) * l
		}

		_out[xOut] = clamp(R / σ)
		_out[xOut+1] = clamp(G / σ)
		_out[xOut+2] = clamp(B / σ)
		_out[xOut+3] = clamp(A / σ)
	}
}

func Resize3(_img *image.NRGBA, _size image.Rectangle, _a float64) *image.NRGBA {
	out := image.NewNRGBA(_size)

	return out
}

type jop struct {
	y   int
	out []uint8
}

func woker(_in []uint8, _rectIn, _rectOut image.Rectangle, y, _stride int, _a, _xStep, ySource float64, c chan jop) {

	out := jop{
		y:   y,
		out: make([]uint8, _rectOut.Dx()*4),
	}

	for x := range _rectOut.Dx() {
		xSource := (float64(x)+0.5)*_xStep - 0.5

		var σ, R, G, B, A float64

		for i := -_a + 1; i <= _a; i++ {
			ty := int(math.Floor(ySource) - i)
			ly := i - ySource + math.Floor(ySource)

			if int(ty) >= 0 && int(ty) < _rectIn.Dy() {
				for j := -_a + 1; j <= _a; j++ {
					tx := int(math.Floor(xSource) - j)
					lx := j - xSource + math.Floor(xSource)

					if int(tx) >= 0 && int(tx) < _rectIn.Dx() {
						l := l2(lx, ly, _a)
						σ += l

						at := int(ty)*_stride + int(tx)*4

						R += float64(_in[at]) * l
						G += float64(_in[at+1]) * l
						B += float64(_in[at+2]) * l
						A += float64(_in[at+3]) * l
					}
				}
			}
		}

		at := x * 4
		out.out[at] = clamp(R / σ)
		out.out[at+1] = clamp(G / σ)
		out.out[at+2] = clamp(B / σ)
		out.out[at+3] = clamp(A / σ)
	}

	c <- out
}

func Resize4(_in, _out []uint8, _rectIn, _rectOut image.Rectangle, _a int, _strideIn, _strideOut int) {
	workers := 0
	channel := make(chan jop, PARA)

	fmt.Println("")

	xStep := float64(_rectIn.Dx()) / float64(_rectOut.Dx())
	yStep := float64(_rectIn.Dy()) / float64(_rectOut.Dy())

	for y := range _rectOut.Dy() {
		ySource := (float64(y)+0.5)*yStep - 0.5

		if workers >= PARA {
			a := <-channel
			copy(_out[a.y*_strideOut:], a.out)
		} else {
			workers++
		}

		go woker(_in, _rectIn, _rectOut, y, _strideIn, float64(_a), xStep, ySource, channel)
	}

	for range workers {
		a := <-channel
		copy(_out[a.y*_strideOut:], a.out)
	}
}

func Patch(_img *image.NRGBA, _size image.Rectangle, _a float64) *image.NRGBA {
	out := image.NewNRGBA(_size)
	Resize4(_img.Pix, out.Pix, _img.Rect, out.Rect, int(_a), _img.Stride, out.Stride)
	return out
}
