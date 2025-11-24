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
	"github.com/fatih/color"
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

type lambda func()

func benchmark(_name string, _fc lambda) {
	start := time.Now()
	for _ = range 10 {
		_fc()
	}
	end := time.Now()

	fmt.Println(_name, end.Sub(start)/10)
}

func runBenchmarks() {
	img, _, _ := loadImage("./example_images/test_uwu.png")
	nrgba := toNrgba(img)
	benchmark("Lanczos uwu", func() {
		_ = filters.Resize(nrgba, nrgba.Rect.Dx()*10, nrgba.Rect.Dy()*10, 3)
	})

	img, _, _ = loadImage("./example_images/test_parrot.jpg")
	nrgba = toNrgba(img)
	benchmark("Lanczos parrot", func() {
		_ = filters.Resize(nrgba, nrgba.Rect.Dx()*10, nrgba.Rect.Dy()*10, 3)
	})

	img, _, _ = loadImage("./example_images/test_circle.jpg")
	nrgba = toNrgba(img)
	benchmark("Lanczos circle", func() {
		_ = filters.Resize(nrgba, nrgba.Rect.Dx()*10, nrgba.Rect.Dy()*10, 3)
	})
}

func printCol(ascii filters.AsciiImage, colors image.NRGBA) {
	for y := range ascii.Height {
		for x := range ascii.Width {
			r := int(colors.Pix[y*colors.Stride+x*4])
			g := int(colors.Pix[y*colors.Stride+x*4+1])
			b := int(colors.Pix[y*colors.Stride+x*4+2])
			color.RGB(r, g, b).Print(ascii.Runes[y*ascii.Width+x])
		}
		fmt.Println()
	}
}

func test_ascii(test [16]uint16) {
	var (
		offset     uint16 = 0
		y                 = 0
		x                 = 0
		_threshold uint16 = 1
	)

	_in := struct {
		Pix    [16]uint16
		Stride int
	}{
		Pix:    test,
		Stride: 4,
	}

	// [OO]
	// [OO]
	// [OO]
	// [XX] <-
	// offset <<= 1
	if _in.Pix[(y+3)*_in.Stride+(x+1)] < _threshold {
		offset += 1
	}

	offset <<= 1
	if _in.Pix[(y+3)*_in.Stride+(x)] < _threshold {
		offset += 1
	}

	// [OX]
	// [OX] ^
	// [OX] |
	// [OO]
	offset <<= 1
	if _in.Pix[(y+2)*_in.Stride+(x+1)] < _threshold {
		offset += 1
	}

	offset <<= 1
	if _in.Pix[(y+1)*_in.Stride+(x+1)] < _threshold {
		offset += 1
	}

	offset <<= 1
	if _in.Pix[(y)*_in.Stride+(x+1)] < _threshold {
		offset += 1
	}

	// [XO]
	// [XO] ^
	// [XO] |
	// [OO]
	offset <<= 1
	if _in.Pix[(y+2)*_in.Stride+(x)] < _threshold {
		offset += 1
	}

	offset <<= 1
	if _in.Pix[(y+1)*_in.Stride+(x)] < _threshold {
		offset += 1
	}

	offset <<= 1
	if _in.Pix[(y)*_in.Stride+(x)] < _threshold {
		offset += 1
	}

	fmt.Printf("%c\n", rune(0x2800+offset))

}

func main() {
	// test_ascii([16]uint16{0, 0, 0, 0,
	// 	0, 0, 0, 0,
	// 	0, 0, 0, 0,
	// 	0, 0, 0, 0})
	//
	// test_ascii([16]uint16{1, 1, 0, 0,
	// 	1, 1, 0, 0,
	// 	1, 1, 0, 0,
	// 	1, 1, 0, 0})
	//
	// test_ascii([16]uint16{0, 1, 0, 0,
	// 	1, 0, 0, 0,
	// 	0, 1, 0, 0,
	// 	1, 0, 0, 0})

	flag.Parse()

	img, _, err := loadImage(in)
	if err != nil {
		fmt.Println(err)
		return
	}

	// gray := grayScale(img)
	// t1 := time.Now()
	// sob := filters.Sobel(filters.In{Pix: gray.Pix, Height: gray.Rect.Dy(), Width: gray.Rect.Dx(), Stride: gray.Stride})
	// t2 := time.Now()
	//
	// fmt.Println(t2.Sub(t1))
	//
	// c := filters.SobelToImg(sob.ToSob())

	nrgba := toNrgba(img)
	var c *image.NRGBA

	width := 200
	height := int(float64(width) / 2)

	t1 := time.Now()
	// c = filters.Resize(nrgba, nrgba.Rect.Dx()/3, nrgba.Rect.Dy()/6, 3)
	c = filters.Resize(nrgba, width, height, 3)
	// c = filters.InvertLuminosity(nrgba)
	t2 := time.Now()
	fmt.Println(t2.Sub(t1))

	// runes := []string{" ", ".", "'", "`", "^", "\"", ",", ":", ";", "I", "l", "!", "i", ">", "<", "~", "+", "_", "-", "?", "]", "[", "}", "{", "1", ")", "(", "|", "\\", "/", "t", "f", "j", "r", "x", "n", "u", "v", "c", "z", "X", "Y", "U", "J", "C", "L", "Q", "0", "O", "Z", "m", "w", "q", "p", "d", "b", "k", "h", "a", "o", "*", "#", "M", "W", "&", "8", "%", "B", "@", "$"}
	// ascii := filters.Ascii(*grayScale(c), runes)
	ascii := filters.Braille(*grayScale(c), 155)

	for y := range ascii.Height {
		for x := range ascii.Width {
			fmt.Printf("%c", ascii.Runes[y*ascii.Width+x])
		}
		fmt.Println()
	}

	// printCol(ascii, *c)

	outfile, _ := os.Create(out)
	defer outfile.Close()
	// png.Encode(outfile, grayScale(c))
	png.Encode(outfile, c)
}
