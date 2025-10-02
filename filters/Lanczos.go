package filters

import (
	"image"
	"image/color"
	"math"
)

const Pi2 = math.Pi * math.Pi

func lanczosKernel(x, a float64) float64 {
	if x == 0 {
		return 1
	} else if -a <= x && x < a {
		return (a * math.Sin(math.Pi*x) * math.Sin(math.Pi*x/a) / (Pi2 * x * x))
	} else {
		return 0
	}
}

func l2(x, y, a float64) float64 {
	return lanczosKernel(x, a) * lanczosKernel(y, a)
}

func Resize(_img image.Image, _size image.Rectangle, a float64) image.Image {
	bounds := _img.Bounds()
	out := image.NewRGBA(_size)

	xStep := float64(bounds.Dx()) / float64(_size.Dx())
	yStep := float64(bounds.Dy()) / float64(_size.Dy())

	for x := _size.Min.X; x < _size.Max.X; x++ {
		xAt := xStep * float64(x+1)

		for y := _size.Min.Y; y < _size.Max.Y; y++ {
			yAt := yStep * float64(y+1)

			newColor := color.RGBA{}

			for i := int(math.Floor(xAt) - a + 1); i <= int(math.Floor(xAt)+a); i++ {
				for j := int(math.Floor(yAt) - a + 1); j <= int(math.Floor(yAt)+a); j++ {
					c := l2(xAt-float64(i), yAt-float64(j), a)
					r, g, b, aa := _img.At(i, j).RGBA()

					newColor.R += uint8(float64(uint8(r >> 8)) * c)
					newColor.G += uint8(float64(uint8(g >> 8)) * c)
					newColor.B += uint8(float64(uint8(b >> 8)) * c)
					newColor.A += uint8(float64(uint8(aa >> 8)) * c)
				}
			}

			out.SetRGBA(x, y, newColor)
		}
	}

	return out
}
