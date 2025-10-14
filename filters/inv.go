package filters

import (
	"image"
)

func worker(in []uint8, Min, Max int, c chan bool) {
	for x := Min; x < Max*4; x += 4 {
		in[x] = ^in[x]
		in[x+1] = ^in[x+1]
		in[x+2] = ^in[x+2]
	}

	c <- true
}

const PARA = 5000

func InvertLuminosity(_img *image.NRGBA) *image.NRGBA {
	workers := 0
	channel := make(chan bool, PARA)

	for y := _img.Rect.Min.Y; y < _img.Rect.Max.Y; y++ {
		if workers >= PARA {
			<-channel
		} else {
			workers++
		}

		start := (y - _img.Rect.Min.Y) * _img.Stride
		end := start + _img.Rect.Dx()*4
		go worker(_img.Pix[start:end], _img.Rect.Min.X, _img.Rect.Max.X, channel)
	}

	return _img
}
