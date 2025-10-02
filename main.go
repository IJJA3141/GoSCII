package main

import (
	"flag"
	"fmt"
	"github.com/IJJA3141/GoSCII/filters"
	"image"
	_ "image/jpeg"
	"image/png"
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

func main() {
	flag.Parse()

	fmt.Println(in)
	fmt.Println(out)

	img, _, err := loadImage(in)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := filters.Resize(img,
		image.Rectangle{
			image.Pt(0, 0),
			image.Pt(int(float64(img.Bounds().Dx())*1.5), int(float64(img.Bounds().Dy())*1.5)),
		}, 1)
	// d := filters.Resize(img, image.Rectangle{image.Pt(0, 0), image.Pt(img.Bounds().Dx()/2, img.Bounds().Dy()/2)}, 5)
	// e := filters.Resize(img, image.Rectangle{image.Pt(0, 0), image.Pt(img.Bounds().Dx()*2, img.Bounds().Dy()/2)}, 5)

	// a := grayScale(img)
	// b := filters.Sobel(a)
	// c := filters.SobelToImg(b)

	// c := filters.InvertLuminosity(img)
	//
	// outfile, _ := os.Create("out1.pgn")
	// defer outfile.Close()
	// png.Encode(outfile, c)
	//
	// c = filters.InvertLuminosity(c)
	// c = filters.InvertLuminosity(c)

	outfile, _ := os.Create(out)
	defer outfile.Close()
	png.Encode(outfile, c)
	//
	// outfile, _ = os.Create("out2.png")
	// defer outfile.Close()
	// png.Encode(outfile, d)
	//
	// outfile, _ = os.Create("out3.png")
	// defer outfile.Close()
	// png.Encode(outfile, e)
}
