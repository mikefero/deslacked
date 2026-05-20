// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package transform

import (
	"image"
	"image/color"
)

// bitsPerByte is used to shift 16-bit color channel values down to 8-bit.
const bitsPerByte = 8

// toGrayscale converts the source image to a fully desaturated NRGBA image.
func toGrayscale(src image.Image) *image.NRGBA {
	bounds := src.Bounds()
	dst := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.At(x, y)
			// Alpha is read from the original pixel before model conversion because
			// color.Gray has no alpha channel — its RGBA() always returns a=0xffff.
			_, _, _, a := c.RGBA()
			gray := color.GrayModel.Convert(c)
			r, _, _, _ := gray.RGBA()
			luma := uint8(r >> bitsPerByte)
			alpha := uint8(a >> bitsPerByte)
			dst.SetNRGBA(x, y, color.NRGBA{R: luma, G: luma, B: luma, A: alpha})
		}
	}
	return dst
}
