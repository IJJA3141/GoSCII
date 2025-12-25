package filters

import (
	"math"
)

// ToRGBA converts a GrayScalePlane to a full RGBAPlane.
//
// Each grayscale pixel value is replicated across the red, green, and blue
// channels. The alpha channel is set to 255 (fully opaque).
//
// This operation does not perform any scaling or gamma correction; it is a
// straightforward replication of the grayscale intensity into RGB channels.
//
// The computation is parallelized across rows for performance.
func (img *GrayScalePlane) ToRGBA() *RGBAPlane {
	out := NewRGBAPlane(img.Width, img.Height)

	split(img.Height, func(_start, _end int) {
		for y := _start; y < _end && y < img.Height; y++ {
			for x := range img.Width {
				index := y*out.Stride + x*4
				out.RGBA[index] = img.Shades[y*img.Stride+x]
				out.RGBA[index+1] = img.Shades[y*img.Stride+x]
				out.RGBA[index+2] = img.Shades[y*img.Stride+x]
				out.RGBA[index+3] = 0xff
			}
		}
	}).Wait()

	return out
}

const r = 0.2126
const g = 0.7152
const b = 0.0722

// ToGrayScale converts an RGBAPlane to a GrayScalePlane.
//
// The grayscale value of each pixel is computed using the Rec. 709
// luminance formula:
//
//	Y = 0.2126*R + 0.7152*G + 0.0722*B
//
// The alpha channel is ignored. The result is clamped to [0, 255] to ensure
// valid pixel intensity values.
//
// The computation is parallelized across rows for performance.
func (img *RGBAPlane) ToGrayScale() *GrayScalePlane {
	out := NewGrayScalePlane(img.Width, img.Height)

	split(img.Height, func(_start, _end int) {
		for y := _start; y < _end && y < img.Height; y++ {
			for x := range img.Width {
				index := y*img.Stride + x*4
				shade := r*img.RGBA[index] + g*img.RGBA[index+1] + b*img.RGBA[index+2]
				out.Shades[y*out.Stride+x] = clamp(shade, 0, 255)
			}
		}
	}).Wait()

	return out
}

func (img *EdgePlane) ToRGBA(threshold float64) *RGBAPlane {
	out := NewRGBAPlane(img.Width, img.Height)

	split(img.Height, func(_start, _end int) {
		for y := _start; y < _end && y < img.Height; y++ {
			for x := range img.Width {
				srcIndex := y*img.Stride + x*2
				outIndex := y*out.Stride + x*4

				magnitude := img.Gradient[srcIndex]
				angle := img.Gradient[srcIndex+1]

				if magnitude >= threshold {
					out.RGBA[outIndex] = float64(uint8(angle / (2 * math.Pi) * 255))
					out.RGBA[outIndex+1] = math.Mod(((angle/(2*math.Pi))+1./3.)*255, 255)
					out.RGBA[outIndex+2] = math.Mod(((angle/(2*math.Pi))+2./3.)*255, 255)
					out.RGBA[outIndex+3] = 0xff
				} else {
					out.RGBA[outIndex] = 0
					out.RGBA[outIndex+1] = 0
					out.RGBA[outIndex+2] = 0
					out.RGBA[outIndex+3] = 0
				}
			}
		}
	}).Wait()

	return out
}
