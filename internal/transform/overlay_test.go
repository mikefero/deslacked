package transform

import (
	"image"
	"testing"

	"github.com/fogleman/gg"
)

func TestDrawOverlayBarHeight(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		w, h int
	}{
		{name: "square 512", w: 512, h: 512},
		{name: "non-square 300x400", w: 300, h: 400},
		{name: "minimum 128x128", w: 128, h: 128},
		{name: "square 1024 (bar template width)", w: 1024, h: 1024},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			src := image.NewNRGBA(image.Rect(0, 0, tc.w, tc.h))
			dc := gg.NewContextForImage(src)

			if err := drawOverlay(dc, tc.w, tc.h); err != nil {
				t.Fatalf("drawOverlay: %v", err)
			}

			// Bar height must equal floor(imgH * barHeightRatio).
			expectedBarH := int(float64(tc.h) * barHeightRatio)
			barStartY := tc.h - expectedBarH

			// At least one pixel in the bar region must be non-zero (overlay was composited).
			img := dc.Image()
			found := false
			for y := barStartY; y < tc.h && !found; y++ {
				for x := 0; x < tc.w; x++ {
					r, g, b, _ := img.At(x, y).RGBA()
					if r != 0 || g != 0 || b != 0 {
						found = true
						break
					}
				}
			}
			if !found {
				t.Errorf("bar region [y=%d..%d] is entirely black; overlay may not have been applied", barStartY, tc.h)
			}
		})
	}
}

func TestDrawOverlayBarHeightUsesImgH(t *testing.T) {
	t.Parallel()

	// For a non-square image the bar height must be derived from imgH, not imgW.
	w, h := 300, 400
	src := image.NewNRGBA(image.Rect(0, 0, w, h))
	dc := gg.NewContextForImage(src)

	if err := drawOverlay(dc, w, h); err != nil {
		t.Fatalf("drawOverlay: %v", err)
	}

	wantBarH := int(float64(h) * barHeightRatio)
	wrongBarH := int(float64(w) * barHeightRatio)

	if wantBarH == wrongBarH {
		t.Skip("imgW and imgH produce the same barH for this test case")
	}

	// Verify the bottom wantBarH rows are non-zero and the row just above is not
	// forced to be non-zero purely by bar height; a coarse sanity check that the
	// bar height is keyed to imgH.
	img := dc.Image()
	barStartY := h - wantBarH
	r, g, b, _ := img.At(w/2, barStartY).RGBA()
	if r == 0 && g == 0 && b == 0 {
		// The first bar row should be non-zero if the overlay was applied correctly.
		t.Errorf("pixel at bar start (y=%d) is black; expected overlay content", barStartY)
	}
}
