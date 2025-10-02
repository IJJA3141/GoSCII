package filters

import (
	"image"
	"image/color"
)

func InvertLuminosity(_img image.Image) image.Image {
	const A = 65535

	bounds := _img.Bounds()
	out := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := _img.At(x, y).RGBA()

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			r8 = ^r8
			g8 = ^g8
			b8 = ^b8

			out.SetRGBA(x, y, color.RGBA{r8, g8, b8, a8})
		}
	}

	return out
}
