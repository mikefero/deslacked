// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package transform

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"

	"github.com/fogleman/gg"
	xdraw "golang.org/x/image/draw"

	"github.com/mikefero/deslacked/assets"
)

const (
	// barTemplateWidth is the pixel width of bar-template.png, captured at 2x Retina resolution.
	barTemplateWidth = 1024
	// barTemplateHeight is the pixel height of bar-template.png, captured at 2x Retina resolution.
	barTemplateHeight = 144
	// barHeightRatio is the fraction of the output image height the overlay bar occupies,
	// derived from the bar template's aspect ratio.
	barHeightRatio = barTemplateHeight / float64(barTemplateWidth)
)

// drawOverlay scales the embedded bar template to fit the output image and composites
// it onto the bottom of dc.
func drawOverlay(dc *gg.Context, imgW, imgH int) error {
	tmpl, err := png.Decode(bytes.NewReader(assets.BarTemplateBytes))
	if err != nil {
		return fmt.Errorf("decoding bar template: %w", err)
	}

	barH := int(float64(imgH) * barHeightRatio)
	barY := imgH - barH

	dst := image.NewRGBA(image.Rect(0, 0, imgW, barH))
	xdraw.BiLinear.Scale(dst, dst.Bounds(), tmpl, tmpl.Bounds(), xdraw.Over, nil)

	// dc.Image() returns *image.NRGBA, which implements draw.Image.
	canvas, ok := dc.Image().(draw.Image)
	if !ok {
		return fmt.Errorf("gg context image does not implement draw.Image")
	}
	draw.Draw(canvas, image.Rect(0, barY, imgW, imgH), dst, image.Point{}, draw.Over)

	return nil
}
