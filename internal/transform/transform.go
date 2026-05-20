// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package transform converts a color profile image into a Slack-style deactivated account image.
package transform

import (
	"fmt"
	"image"
	"image/png"
	"os"

	// Imported for side-effect image format registration.
	_ "image/gif"
	_ "image/jpeg"

	"github.com/fogleman/gg"
)

// ProcessOptions controls how the image is transformed.
type ProcessOptions struct {
	// Anchor controls where the crop window sits vertically when the input is
	// non-square. Ignored when CropOffset is non-nil.
	Anchor AnchorPosition
	// CropOffset is the fractional position (0.0 = top/left, 1.0 = bottom/right)
	// of the crop box along the longer axis. nil means use Anchor instead.
	CropOffset *float64
}

// DefaultProcessOptions returns ProcessOptions with sensible defaults.
func DefaultProcessOptions() ProcessOptions {
	return ProcessOptions{Anchor: AnchorCenter}
}

// Process reads an image from inputPath, normalizes it to a square at or below 512px
// (cropping non-square inputs using opts.CropOffset when set, otherwise opts.Anchor,
// and scaling down oversized inputs), converts to grayscale, overlays the
// deactivated-account bar, and writes the result as a PNG to outputPath.
func Process(inputPath, outputPath string, opts ProcessOptions) error {
	src, err := loadImage(inputPath)
	if err != nil {
		return err
	}
	result, err := processImage(src, opts)
	if err != nil {
		return err
	}
	return savePNG(outputPath, result)
}

// ProcessImage transforms an already-decoded image entirely in memory, applying
// the same pipeline as Process: normalize to square, convert to grayscale, overlay
// the deactivated-account bar. The returned image.Image is the processed result.
func ProcessImage(src image.Image, opts ProcessOptions) (image.Image, error) {
	return processImage(src, opts)
}

// processImage is the shared implementation used by both Process and ProcessImage.
func processImage(src image.Image, opts ProcessOptions) (image.Image, error) {
	var (
		normalized image.Image
		err        error
	)
	if opts.CropOffset != nil {
		normalized, err = cropAndResizeFromOffset(src, *opts.CropOffset)
	} else {
		normalized, err = cropAndResize(src, opts.Anchor)
	}
	if err != nil {
		return nil, fmt.Errorf("normalizing image: %w", err)
	}

	gray := toGrayscale(normalized)
	bounds := gray.Bounds()
	dc := gg.NewContextForImage(gray)

	if err = drawOverlay(dc, bounds.Dx(), bounds.Dy()); err != nil {
		return nil, fmt.Errorf("drawing overlay: %w", err)
	}
	return dc.Image(), nil
}

// loadImage opens and decodes an image from the given path.
func loadImage(path string) (image.Image, error) {
	//nolint:gosec
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening input image: %w", err)
	}
	defer func() { _ = f.Close() }()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decoding input image: %w", err)
	}
	return img, nil
}

// savePNG encodes img as PNG to the given path.
func savePNG(path string, img image.Image) error {
	//nolint:gosec
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}

	if err = png.Encode(f, img); err != nil {
		_ = f.Close()
		return fmt.Errorf("encoding PNG: %w", err)
	}
	return f.Close()
}
