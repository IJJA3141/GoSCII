package io

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/IJJA3141/GoSCII/filters"
)

func Read(_path string) (*filters.RGBAPlane, error) {
	file, err := os.Open(_path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	image, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	out := filters.NewRGBAPlane(image.Bounds().Dx(), image.Bounds().Dy())
	for y := range out.Height {
		for x := range out.Width {
			r, g, b, a := image.At(x, y).RGBA()

			index := y*out.Stride + x*4
			out.RGBA[index] = float64(r >> 8)
			out.RGBA[index+1] = float64(g >> 8)
			out.RGBA[index+2] = float64(b >> 8)
			out.RGBA[index+3] = float64(a >> 8)
		}
	}

	return &out, nil
}
