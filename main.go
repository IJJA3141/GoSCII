package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
)

func createBoolFlag(_ptr *bool, _name string, _value bool, _desc string) {
	flag.BoolVar(_ptr, _name, _value, _desc)
	flag.BoolVar(_ptr, string(_name[0]), _value, _desc)
}

func createIntFlag(_ptr *int, _name string, _value int, _desc string) {
	flag.IntVar(_ptr, _name, _value, _desc)
	flag.IntVar(_ptr, string(_name[0]), _value, _desc)
}

func createStringFlag(_ptr *string, _name, _value, _desc string) {
	flag.StringVar(_ptr, _name, _value, _desc)
	flag.StringVar(_ptr, string(_name[0]), _value, _desc)
}

var in string
var out string

func init() {
	createStringFlag(&in, "in", "./test_uwu.png", "path to the input image")
	createStringFlag(&out, "out", "out.png", "path to the output image")
}

func grayScale(_img image.Image) *image.Gray {
	bounds := _img.Bounds()
	img := image.NewGray(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			img.Set(x, y, _img.At(x, y))
		}
	}

	return img
}

func loadImage(_path string) (image.Image, string, error) {
	file, err := os.Open(_path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	image, format, err := image.Decode(file)
	return image, format, err
}

func convolution(_img *image.Gray, _kern [][]float64) [][]float64 {
	kernBounds := image.Rectangle{image.Pt(0, 0), image.Pt(len(_kern), len(_kern[0]))}
	xOffset := len(_kern) / 2
	yOffset := len(_kern[0]) / 2
	bounds := _img.Bounds()
	out := make([][]float64, bounds.Dx())

	for x := bounds.Min.X; x < bounds.Max.X; x++ {

		out[x] = make([]float64, bounds.Dy())
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {

			for i := kernBounds.Min.X; i < kernBounds.Max.X; i++ {
				for j := kernBounds.Min.Y; j < kernBounds.Max.Y; j++ {
					r, _, _, _ := _img.At(x-xOffset+i, y-yOffset+j).RGBA()
					out[x][y] += float64(r) * _kern[i][j]
				}
			}

		}
	}

	return out
}

type sob struct{ gradiant, angle float64 }

func sobel(_img *image.Gray) [][]sob {
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

	out := make([][]sob, bounds.Dx())

	for x := bounds.Min.X; x < bounds.Max.X; x++ {

		out[x] = make([]sob, bounds.Dy())

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			out[x][y] = sob{math.Sqrt(Mx[x][y]*Mx[x][y] + My[x][y]*My[x][y]), math.Atan2(My[x][y], Mx[x][y])}
		}
	}

	return out
}

func sobelToImg(_so [][]sob) image.Image {
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

func main() {
	flag.Parse()

	fmt.Println(in)
	fmt.Println(out)

	img, _, err := loadImage(in)
	if err != nil {
		fmt.Println(err)
		return
	}

	a := grayScale(img)
	b := sobel(a)
	c := sobelToImg(b)

	outfile, _ := os.Create(out)
	defer outfile.Close()
	png.Encode(outfile, c)
}
