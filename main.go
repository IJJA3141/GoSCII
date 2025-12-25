package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"

	"github.com/IJJA3141/GoSCII/io"
)

var in string
var out string

func init() {
	io.CreateStringFlag(&in, "in", "./example_images/test_uwu.png", "path to the input image")
	io.CreateStringFlag(&out, "out", "out.png", "path to the output image")
}

func main() {
	flag.Parse()

	img, err := io.Read(in)
	if err != nil {
		fmt.Println(err)
		return
	}

	outWidth := 80
	outHeight := int(80/3)
	// outHeight := float64(img.Height) / float64(img.Width) * float64(outWidth)

	img, err = img.LanczosResize(outWidth*2, int(outHeight*4), 3)
	if err != nil {
		fmt.Println(err)
		return
	}

	gray := img.ToGrayScale()

	// gray, err = filters.BayerDithering(gray, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// img, err = img.LanczosResize(outWidth, int(outHeight), 3)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// ascii := filters.Braille(gray, 1)
	// ascii := gray.Ascii([]rune(
	// 	" .:,`';^-_!~\"</>*+?\\v)x=cJY|Lil{}7T(1CetzVXnorsaujyUfI]23AFHZ5S[K#%4hw6&KOp9PbGmdq$08DERNQgMWB@",
	// ))
	edge := gray.SobelEdgeDetection()
	ascii := edge.Ascii(750, []rune(
		// "→↗↑↖←↙↓↘",
		// "↖↖",
		// "123455678",
		// "←↖↑↗→↘↓↙←",
		"|/-\\|/-\\|",
	))

	color, err := ascii.Colorize(edge.ToRGBA(750))
	if err != nil {
		fmt.Println(err)
		return
	}

	io.StampC(color)
	err = io.Write(out, edge.ToRGBA(750))
	if err != nil {
		fmt.Println(err)
		return
	}
}
