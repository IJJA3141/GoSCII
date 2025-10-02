package filters

import "image"

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
