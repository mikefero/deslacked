// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package transform

import (
	"fmt"
	"image"
	"image/draw"

	xdraw "golang.org/x/image/draw"
)

// AnchorPosition controls where the crop window sits vertically when squaring a
// non-square image.
type AnchorPosition string

const (
	// AnchorTop aligns the crop rectangle to the top of the image (y = 0).
	AnchorTop AnchorPosition = "top"
	// AnchorCenter centers the crop rectangle vertically (default).
	AnchorCenter AnchorPosition = "center"
	// AnchorBottom aligns the crop rectangle to the bottom of the image (y = h − side).
	AnchorBottom AnchorPosition = "bottom"
)

// minSize is the smallest input dimension accepted without error. Inputs below this
// would require upscaling, which degrades quality unacceptably for profile images.
const minSize = 128

// maxSize is the target output dimension. Inputs larger than this are scaled down;
// inputs at or below this are returned at their original size.
const maxSize = 512

// centerDivisor is used when halving a dimension to find a center offset.
const centerDivisor = 2

const (
	// offsetClampMin is the minimum valid value for a CropOffset fraction.
	offsetClampMin = 0.0
	// offsetClampMax is the maximum valid value for a CropOffset fraction.
	offsetClampMax = 1.0
)

// cropAndResize normalizes img to a square at or below maxSize following this pipeline:
//
//  1. Reject any image with either dimension < minSize (128px).
//  2. If the image is NOT square: crop to a square with side = min(w, h).
//     Horizontally the crop is always centered; vertically it is positioned by anchor
//     (top → y=0, bottom → y=h−side, center → midpoint).
//  3. If the resulting square exceeds maxSize (512px): scale it down to maxSize×maxSize
//     using bilinear interpolation.
//  4. If the image is already square and within [minSize, maxSize]: return unchanged.
//
// The bar overlay applied downstream always scales dynamically via barHeightRatio.
func cropAndResize(img image.Image, anchor AnchorPosition) (image.Image, error) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	if w < minSize || h < minSize {
		return nil, fmt.Errorf("image dimensions %dx%d are below the minimum accepted size of %dpx", w, h, minSize)
	}

	result := img
	if w != h {
		rect := squareCropRect(w, h, anchor)
		result = cropToRect(img, rect)
		bounds = result.Bounds()
		w = bounds.Dx()
	}

	if w > maxSize {
		result = scaleDown(result, maxSize)
	}

	return result, nil
}

// cropAndResizeFromOffset is like cropAndResize but uses a continuous fractional
// offset (0.0 = top/left, 1.0 = bottom/right) along the longer axis instead of
// a discrete AnchorPosition. This is the path used by the web UI's draggable crop box.
func cropAndResizeFromOffset(img image.Image, offset float64) (image.Image, error) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	if w < minSize || h < minSize {
		return nil, fmt.Errorf("image dimensions %dx%d are below the minimum accepted size of %dpx", w, h, minSize)
	}

	result := img
	if w != h {
		rect := squareCropRectWithOffset(w, h, offset)
		result = cropToRect(img, rect)
		bounds = result.Bounds()
		w = bounds.Dx()
	}

	if w > maxSize {
		result = scaleDown(result, maxSize)
	}

	return result, nil
}

// squareCropRectWithOffset computes a square crop rectangle using a fractional offset
// clamped to [0, 1]. For landscape images (w > h) the box slides horizontally; for
// portrait (h > w) it slides vertically. The perpendicular axis is always centered.
func squareCropRectWithOffset(imgW, imgH int, offset float64) image.Rectangle {
	if offset < offsetClampMin {
		offset = offsetClampMin
	}
	if offset > offsetClampMax {
		offset = offsetClampMax
	}

	side := min(imgW, imgH)

	if imgW > imgH {
		// Landscape: slide horizontally; the full height is used as the square side.
		maxX := imgW - side
		x0 := int(float64(maxX) * offset)
		return image.Rect(x0, 0, x0+side, side)
	}
	// Portrait: slide vertically; the full width is used as the square side.
	x0 := (imgW - side) / centerDivisor
	maxY := imgH - side
	y0 := int(float64(maxY) * offset)
	return image.Rect(x0, y0, x0+side, y0+side)
}

// squareCropRect computes a square crop rectangle with side = min(imgW, imgH).
// The rectangle is always horizontally centered and vertically positioned per anchor.
func squareCropRect(imgW, imgH int, anchor AnchorPosition) image.Rectangle {
	side := min(imgW, imgH)

	x0 := (imgW - side) / centerDivisor
	var y0 int
	switch anchor {
	case AnchorTop:
		y0 = 0
	case AnchorBottom:
		y0 = imgH - side
	default: // AnchorCenter
		y0 = (imgH - side) / centerDivisor
	}

	return image.Rect(x0, y0, x0+side, y0+side)
}

// cropToRect extracts the subimage defined by rect from src, returning a concrete
// *image.NRGBA with Min=(0,0) so downstream gg.NewContextForImage works correctly.
func cropToRect(src image.Image, rect image.Rectangle) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(dst, dst.Bounds(), src, rect.Min, draw.Src)
	return dst
}

// scaleDown reduces src to targetSize×targetSize using bilinear interpolation.
// Bilinear is chosen over nearest-neighbor to avoid aliasing at common profile
// image downscale ratios (e.g. 1024→512).
func scaleDown(src image.Image, targetSize int) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, targetSize, targetSize))
	xdraw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), xdraw.Over, nil)
	return dst
}
