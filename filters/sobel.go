package filters

import (
	"image"
	"image/color"
	"math"
)

type Sob struct{ gradiant, angle float64 }

func Sobel(_img *image.Gray) [][]Sob {
	Gx := [][]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	Gy := [][]float64{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	Mx := convolution(_img, Gx)
	My := convolution(_img, Gy)

	bounds := _img.Bounds()

	out := make([][]Sob, bounds.Dx())

	for x := bounds.Min.X; x < bounds.Max.X; x++ {

		out[x] = make([]Sob, bounds.Dy())

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			out[x][y] = Sob{math.Sqrt(Mx[x][y]*Mx[x][y] + My[x][y]*My[x][y]), math.Atan2(My[x][y], Mx[x][y])}
		}
	}

	return out
}

func SobelToImg(_so [][]Sob) image.Image {
	bounds := image.Rectangle{image.Pt(0, 0), image.Pt(len(_so), len(_so[0]))}
	out := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			s := _so[x][y]
			// if s.gradiant > 50000 {
			if s.gradiant > 30000 {
				// if s.gradiant > 10000 {
				// if s.gradiant > 5000 {
				// if s.gradiant > 3000 {
				// if s.gradiant > 1000 {
				// if s.gradiant > 50 {
				// if s.gradiant > 1 {
				var a uint8 = 255

				hue := math.Mod(s.angle, 2*math.Pi) / (2 * math.Pi)

				// HSV to RGB (with full saturation=1, value=1)
				h := hue * 6
				i := int(math.Floor(h))
				f := h - float64(i)
				p := 0.0
				q := 1 - f
				t := f

				var r, g, b float64
				switch i % 6 {
				case 0:
					r, g, b = 1, t, p
				case 1:
					r, g, b = q, 1, p
				case 2:
					r, g, b = p, 1, t
				case 3:
					r, g, b = p, q, 1
				case 4:
					r, g, b = t, p, 1
				case 5:
					r, g, b = 1, p, q
				}

				out.SetRGBA(x, y, color.RGBA{uint8(255 * r), uint8(255 * g), uint8(255 * b), a})
			}
		}
	}

	return out
}
