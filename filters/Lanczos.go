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

func l2(_x, _y float64, _a int) float64 {
	// return l(math.Sqrt(_x*_x+_y*_y), _a)
	return l(_x, _a) * l(_y, _a)
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

type jop struct {
	y   int
	row []uint8
}

func woker(
	_in []uint8, _widthIn, _heightIn, _strideIn int, // in
	_widthOut, _yOut, // out
	_a int, _xStep, _yStep float64, _channel chan jop) {

	yIn := (float64(_yOut)+0.5)*_yStep - 0.5
	out := jop{
		y:   _yOut,
		row: make([]uint8, _widthOut*4),
	}

	for xOut := range _widthOut {
		xIn := (float64(xOut)+0.5)*_xStep - 0.5

		var σ, R, G, B, A float64

		for i := -_a + 1; i <= _a; i++ {
			ly := float64(i) - yIn + math.Floor(yIn)
			yKer := int(math.Floor(yIn)) - i

			if 0 <= yKer && yKer < _heightIn {
				for j := -_a + 1; j <= _a; j++ {
					lx := float64(j) - xIn + math.Floor(xIn)

					xKer := int(math.Floor(xIn)) - j

					if 0 <= xKer && xKer < _widthIn {
						l := l2(lx, ly, _a)
						σ += l

						index := int(yKer)*_strideIn + int(xKer)*4
						R += float64(_in[index]) * l
						G += float64(_in[index+1]) * l
						B += float64(_in[index+2]) * l
						A += float64(_in[index+3]) * l
					}
				}
			}
		}

		index := xOut * 4
		out.row[index] = clamp(R / σ)
		out.row[index+1] = clamp(G / σ)
		out.row[index+2] = clamp(B / σ)
		out.row[index+3] = clamp(A / σ)
	}

	_channel <- out
}

func Resize(_in *image.NRGBA, _width, _height, _nWorkers, _a int) *image.NRGBA {
	out := image.NewNRGBA(image.Rect(0, 0, _width, _height))

	workers := 0
	channel := make(chan jop, _nWorkers)

	xStep := float64(_in.Rect.Dx()) / float64(_width)
	yStep := float64(_in.Rect.Dy()) / float64(_height)

	for y := range _height {
		if workers < _nWorkers {
			workers++
		} else {
			a := <-channel
			copy(out.Pix[a.y*out.Stride:], a.row)
		}

		go woker(_in.Pix, _in.Rect.Dx(), _in.Rect.Dy(), _in.Stride, _width, y, _a, xStep, yStep, channel)
	}

	for range workers {
		a := <-channel
		copy(out.Pix[a.y*out.Stride:], a.row)
	}

	return out
}
