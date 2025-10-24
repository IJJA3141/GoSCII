package filters

import "image"

func Range(_in image.Gray, _min, _max uint8) *image.Gray {
	out := image.NewGray(_in.Rect)

	Split(_in.Rect.Dy(), func(_start, _end int) {
		for y := _start; y < _end && y < _in.Rect.Dy(); y++ {
			for x := range _in.Rect.Dx() {
				pix := _in.Pix[y*_in.Stride+x]
				pix = max(_min, pix)
				pix = min(pix, _max)
				out.Pix[y*_in.Stride+x] = pix
			}
		}
	})

	return out
}
