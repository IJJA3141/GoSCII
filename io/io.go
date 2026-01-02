package io

import (
	"image"
	_ "image/jpeg"
	"image/png"
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

	return out, nil
}

func Write(_path string, _img *filters.RGBAPlane) error {
	outfile, err := os.Create(_path)
	if err != nil {
		return err
	}

	defer outfile.Close()

	pix := make([]uint8, _img.Width*_img.Height*4)
	for y := range _img.Height {
		for x := range _img.Width {
			index := y*_img.Stride + x*4

			pix[index] = uint8(_img.RGBA[index])
			pix[index+1] = uint8(_img.RGBA[index+1])
			pix[index+2] = uint8(_img.RGBA[index+2])
			pix[index+3] = uint8(_img.RGBA[index+3])
		}
	}

	png.Encode(outfile, &image.RGBA{
		Pix:    pix,
		Stride: _img.Stride,
		Rect: image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: _img.Width, Y: _img.Height},
		},
	},
	)

	return nil
}
