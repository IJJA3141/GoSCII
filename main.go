package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"reflect"
	"time"

	"github.com/IJJA3141/GoSCII/filters"
)

func getType(myvar any) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Pointer {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

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
	createStringFlag(&in, "in", "./example_images/test_uwu.png", "path to the input image")
	createStringFlag(&out, "out", "out.png", "path to the output image")
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

func toNrgba(_img image.Image) *image.NRGBA {
	nrgba, ok := _img.(*image.NRGBA)
	if ok {
		return nrgba
	}

	out := image.NewNRGBA(_img.Bounds())
	for x := range _img.Bounds().Dx() {
		for y := range _img.Bounds().Dy() {
			out.Set(x, y, _img.At(x, y))
		}
	}

	return out
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

func main() {
	flag.Parse()

	img, _, err := loadImage(in)
	if err != nil {
		fmt.Println(err)
		return
	}

	gray := grayScale(img)
	t1 := time.Now()
	_= filters.Sobel(gray)
	t2 := time.Now()

	t3 := time.Now()
	sob2 := filters.Sobel2(gray)
	t4 := time.Now()

	fmt.Println(t2.Sub(t1))
	fmt.Println(t4.Sub(t3))

	c := filters.SobelToImg(sob2.ToSob())

	// nrgba := toNrgba(img)
	// var c *image.NRGBA
	//
	// t1 := time.Now()
	// c = filters.Resize(nrgba, nrgba.Rect.Dx()*10, nrgba.Rect.Dy()*10, runtime.GOMAXPROCS(0), 3)
	// // c = filters.InvertLuminosity(nrgba)
	// t2 := time.Now()
	// fmt.Println(t2.Sub(t1))
	//

	outfile, _ := os.Create(out)
	defer outfile.Close()
	png.Encode(outfile, c)
}
