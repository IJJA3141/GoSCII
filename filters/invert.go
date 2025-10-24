package filters

import (
	"image"
)

func Invert(_in *image.NRGBA) *image.NRGBA {
	out := image.NewNRGBA(_in.Rect)

	Split(_in.Rect.Dy(),
		func(_start, _end int) {
			for y := _start; y < _end && y < _in.Rect.Dy(); y++ {
				for x := range _in.Rect.Dx() {
					index := y*out.Stride + x*4
					out.Pix[index] = ^_in.Pix[index]
					out.Pix[index+1] = ^_in.Pix[index+1]
					out.Pix[index+2] = ^_in.Pix[index+2]
				}
			}
		}).Wait()

	return out
}
