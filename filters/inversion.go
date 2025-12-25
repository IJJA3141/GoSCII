package filters

// Inverse returns a new RGBAPlane where the color channels (R, G, B) of each
// pixel are inverted, producing a photographic negative of the original image.
//
// The alpha channel (A) is preserved and not inverted, so transparency remains unchanged.
//
// Each color channel inversion is performed as:
//
//	newValue = 255 - oldValue
//
// This function is parallelized across rows to improve performance for large images.
func (img *RGBAPlane) Inverse() *RGBAPlane {
	out := NewRGBAPlane(img.Width, img.Height)

	split(img.Height, func(_start, _end int) {
		for y := _start; y < _end && y < img.Height; y++ {
			for x := range img.Width {
				index := y*out.Stride + x*4

				out.RGBA[index] = float64(^uint8(img.RGBA[index]))
				out.RGBA[index+1] = float64(^uint8(img.RGBA[index+1]))
				out.RGBA[index+2] = float64(^uint8(img.RGBA[index+2]))
				// out.RGBA[index+3] = float64(^uint8(img.RGBA[index+3])) // dont inverse alpha channel
			}
		}
	}).Wait()

	return out
}

// Inverse returns a new GrayScalePlane where the shade of each
// pixel is inverted, producing a photographic negative of the original image.
//
// This function is parallelized across rows to improve performance for large images.
func (img *GrayScalePlane) Inverse() *GrayScalePlane {
	out := NewGrayScalePlane(img.Width, img.Height)

	split(img.Height, func(_start, _end int) {
		for y := _start; y < _end && y < img.Height; y++ {
			for x := range img.Width {
				out.Shades[y*out.Stride+x] = float64(^uint8(img.Shades[y*img.Stride+x]))
			}
		}
	}).Wait()

	return out
}
