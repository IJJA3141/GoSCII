package filters

import (
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
