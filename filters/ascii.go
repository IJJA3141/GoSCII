package filters

var dotMatrix = [][]uint16{
	{0x1, 0x8},
	{0x2, 0x10},
	{0x4, 0x20},
	{0x40, 0x80},
}

func (img *GrayScalePlane) Braille(threshold float64) *AsciiPlane {
	out := NewAsciiPlane(img.Width/2, img.Height/4)

	split(out.Height, func(_start, _end int) {
		var char uint16

		for y := _start; y < _end && y < out.Height; y++ {
			for x := range out.Width {
				char = 0x2800 // block offset

				for j := range 4 {
					for i := range 2 {
						if img.Shades[(y*4+j)*img.Stride+(x*2+i)] >= threshold {
							char += dotMatrix[j][i]
						}
					}
				}

				out.Chars[y*out.Stride+x] = rune(char)
			}
		}
	}).Wait()

	return &out
}

func (pln AsciiPlane) Dimensions() (int, int) { return pln.Width, pln.Height }
func (pln *AsciiPlane) Get(x, y, width, height int) []string {
	out := make([]string, height)

	for i := range height {
		index := (y+i)*pln.Stride + x
		out[i] = string(pln.Chars[index : index+width])
	}

	return out
}
