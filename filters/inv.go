package filters

import (
	"image"
	"image/color"
)

func InvertLuminosity(_img image.Image) image.Image {
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

// func job(in []uint8, channel chan bool) {
// 	for x := 0; x < len(in); x += 4 {
// 		in[x] = ^in[x]
// 		in[x+1] = ^in[x+1]
// 		in[x+2] = ^in[x+2]
// 	}
//
// 	channel <- true
// }
//
// func InvertLuminosity2(_img *image.NRGBA) *image.NRGBA {
// 	channel := make(chan bool)
//
// 	for y := _img.Bounds().Min.Y; y < _img.Bounds().Max.Y; y++ {
// 		start := (y - _img.Rect.Min.Y) * _img.Stride
// 		end := start + _img.Rect.Dx()*4 // 4 bytes per pixel (R,G,B,A)
// 		go job(_img.Pix[start:end], channel)
// 	}
//
// 	for range _img.Bounds().Dy() {
// 		<-channel
// 	}
//
// 	return _img
// }

func job(in []uint8) {
	for x := 0; x < len(in); x += 4 {
		in[x] = ^in[x]
		in[x+1] = ^in[x+1]
		in[x+2] = ^in[x+2]
	}
}

func InvertLuminosity2(_img *image.NRGBA) *image.NRGBA {
	for y := _img.Bounds().Min.Y; y < _img.Bounds().Max.Y; y++ {
		start := (y - _img.Rect.Min.Y) * _img.Stride
		end := start + _img.Rect.Dx()*4 // 4 bytes per pixel (R,G,B,A)
		go job(_img.Pix[start:end])
	}

	return _img
}
