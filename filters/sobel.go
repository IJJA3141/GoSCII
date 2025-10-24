package filters

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
)

type Sob struct{ Gradiant, Angle float64 }

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

type NamingIsHard struct {
	// Pairs holds the (gradiant, angle) pairs,
	// The pair at (x, y) starts at Pairs[y*Stride + x*2].
	Pairs []float64

	// Stride is the Pairs stride between vertically adjacent pairs.
	Height, Width, Stride int
}

func Sobel2(_in *image.Gray) NamingIsHard {
	Gx := [3 * 3]float64{
		-1, 0, 1,
		-2, 0, 2,
		-1, 0, 1,
	}

	Gy := [3 * 3]float64{
		-1, -2, -1,
		0, 0, 0,
		1, 2, 1,
	}

	out := NamingIsHard{
		Stride: _in.Rect.Dx() * 2,
		Pairs:  make([]float64, 2*_in.Rect.Dx()*_in.Rect.Dy()),
		Height: _in.Rect.Dy(),
		Width:  _in.Rect.Dx(),
	}

	var wg sync.WaitGroup

	cpus := runtime.GOMAXPROCS(0)
	cpus = 1
	wg.Add(cpus)

	split := int(_in.Rect.Dy() / cpus)

	for cpu := range cpus {
		go func() {
			defer wg.Done()

			for y := cpu * split; y < (cpu+1)*split && y < _in.Rect.Dy(); y++ {
				for x := range _in.Rect.Dx() {
					var px, py float64

					for j := range 3 {
						yj := y - 1 + j

						for i := range 3 {
							xi := x - 1 + i

							if 0 <= yj && yj < _in.Rect.Dy() && 0 <= xi && xi < _in.Rect.Dx() {
								pix := float64(_in.Pix[yj*_in.Stride+xi])

								px += pix * Gx[j * 3 + i]
								py += pix * Gy[j * 3 + i]
							}
						}
					}

					out.Pairs[y*out.Stride+2*x] = math.Sqrt(px*px + py*py)
					out.Pairs[y*out.Stride+2*x+1] = math.Atan2(py, px)
				}
			}
		}()
	}

	wg.Wait()

	return out
}

func (this *NamingIsHard) ToSob() [][]Sob {
	out := make([][]Sob, this.Width)

	for x := range this.Width {
		out[x] = make([]Sob, this.Height)

		for y := range this.Height {
			out[x][y] = Sob{this.Pairs[y*this.Stride+x*2], this.Pairs[y*this.Stride+x*2+1]}
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
			if s.Gradiant > 100 {
				// if s.gradiant > 10000 {
				// if s.gradiant > 5000 {
				// if s.gradiant > 3000 {
				// if s.gradiant > 1000 {
				// if s.gradiant > 50 {
				// if s.gradiant > 1 {
				var a uint8 = 255

				hue := math.Mod(s.Angle, 2*math.Pi) / (2 * math.Pi)

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

				if r > g && r > b {
					fmt.Println(_so[x][y].Angle)
				}

				out.SetRGBA(x, y, color.RGBA{uint8(255 * r), uint8(255 * g), uint8(255 * b), a})
			}
		}
	}

	return out
}
