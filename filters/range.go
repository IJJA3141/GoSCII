package filters

import "image"

func Range(_in image.Gray, _min, _max uint8) *image.Alpha {
	out := image.NewAlpha(_in.Rect)

	Split(_in.Rect.Dy(), func(_start, _end int) {
		for y := _start; y < _end && y < _in.Rect.Dy(); y++ {
			for x := range _in.Rect.Dx() {
				pix := _in.Pix[y*_in.Stride+x]

				if _min < pix && pix < _max {
					out.Pix[y*_in.Stride+x] = 255
				} else {
					out.Pix[y*_in.Stride+x] = 0
				}
			}
		}
	})

	return out
}
