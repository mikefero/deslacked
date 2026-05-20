package transform

import (
	"image"
	"testing"
)

// makeImage creates a blank w x h NRGBA image for use in tests.
func makeImage(w, h int) image.Image {
	return image.NewNRGBA(image.Rect(0, 0, w, h))
}

func TestCropAndResize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		w, h    int
		anchor  AnchorPosition
		wantErr bool
		wantW   int
		wantH   int
	}{
		{
			name: "reject tiny square",
			w:    64, h: 64,
			anchor:  AnchorCenter,
			wantErr: true,
		},
		{
			name: "reject narrow image",
			w:    64, h: 512,
			anchor:  AnchorCenter,
			wantErr: true,
		},
		{
			name: "reject short image",
			w:    512, h: 64,
			anchor:  AnchorCenter,
			wantErr: true,
		},
		{
			name: "square small accepted unchanged",
			w:    200, h: 200,
			anchor: AnchorCenter,
			wantW:  200, wantH: 200,
		},
		{
			name: "square at maxSize accepted unchanged",
			w:    512, h: 512,
			anchor: AnchorCenter,
			wantW:  512, wantH: 512,
		},
		{
			name: "square large scaled to maxSize",
			w:    1024, h: 1024,
			anchor: AnchorCenter,
			wantW:  512, wantH: 512,
		},
		{
			name: "non-square landscape center crop then scale",
			w:    1024, h: 768,
			anchor: AnchorCenter,
			wantW:  512, wantH: 512,
		},
		{
			name: "non-square portrait center crop then scale",
			w:    768, h: 1024,
			anchor: AnchorCenter,
			wantW:  512, wantH: 512,
		},
		{
			name: "non-square top anchor produces square output",
			w:    1024, h: 768,
			anchor: AnchorTop,
			wantW:  512, wantH: 512,
		},
		{
			name: "non-square bottom anchor produces square output",
			w:    1024, h: 768,
			anchor: AnchorBottom,
			wantW:  512, wantH: 512,
		},
		{
			name: "non-square small fits without scaling",
			w:    300, h: 200,
			anchor: AnchorCenter,
			wantW:  200, wantH: 200,
		},
		{
			name: "square at minSize accepted unchanged",
			w:    128, h: 128,
			anchor: AnchorCenter,
			wantW:  128, wantH: 128,
		},
		{
			name: "HD landscape 1920x1080 cropped and scaled to 512x512",
			w:    1920, h: 1080,
			anchor: AnchorCenter,
			wantW:  512, wantH: 512,
		},
		{
			name: "irregular portrait 333x444 cropped without scaling",
			w:    333, h: 444,
			anchor: AnchorCenter,
			wantW:  333, wantH: 333,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			img := makeImage(tc.w, tc.h)
			got, err := cropAndResize(img, tc.anchor)

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for %dx%d input, got none", tc.w, tc.h)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			b := got.Bounds()
			if b.Dx() != tc.wantW || b.Dy() != tc.wantH {
				t.Errorf("output size: got %dx%d, want %dx%d", b.Dx(), b.Dy(), tc.wantW, tc.wantH)
			}
		})
	}
}

func TestCropAndResizeFromOffset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		w, h    int
		offset  float64
		wantErr bool
		wantW   int
		wantH   int
	}{
		{
			name:    "reject too-small image",
			w:       64,
			h:       64,
			offset:  0.5,
			wantErr: true,
		},
		{
			name:   "square at maxSize unchanged",
			w:      512,
			h:      512,
			offset: 0.5,
			wantW:  512,
			wantH:  512,
		},
		{
			name:   "portrait offset 0.0 produces square",
			w:      300,
			h:      400,
			offset: 0.0,
			wantW:  300,
			wantH:  300,
		},
		{
			name:   "portrait offset 1.0 produces square",
			w:      300,
			h:      400,
			offset: 1.0,
			wantW:  300,
			wantH:  300,
		},
		{
			name:   "landscape offset 0.5 produces square",
			w:      400,
			h:      300,
			offset: 0.5,
			wantW:  300,
			wantH:  300,
		},
		{
			name:   "offset below 0 clamped to 0",
			w:      300,
			h:      400,
			offset: -0.5,
			wantW:  300,
			wantH:  300,
		},
		{
			name:   "offset above 1 clamped to 1",
			w:      300,
			h:      400,
			offset: 1.5,
			wantW:  300,
			wantH:  300,
		},
		{
			name:   "large non-square scaled down",
			w:      1024,
			h:      800,
			offset: 0.5,
			wantW:  512,
			wantH:  512,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			img := makeImage(tc.w, tc.h)
			got, err := cropAndResizeFromOffset(img, tc.offset)

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for %dx%d input, got none", tc.w, tc.h)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			b := got.Bounds()
			if b.Dx() != tc.wantW || b.Dy() != tc.wantH {
				t.Errorf("output size: got %dx%d, want %dx%d", b.Dx(), b.Dy(), tc.wantW, tc.wantH)
			}
		})
	}
}

func TestSquareCropRectWithOffset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		imgW, imgH int
		offset     float64
		wantX0     int
		wantY0     int
		wantSide   int
	}{
		// Portrait 300x400: side=300, x0=0, maxY=100.
		{name: "portrait offset 0.0 (top)", imgW: 300, imgH: 400, offset: 0.0, wantX0: 0, wantY0: 0, wantSide: 300},
		{name: "portrait offset 0.5 (middle)", imgW: 300, imgH: 400, offset: 0.5, wantX0: 0, wantY0: 50, wantSide: 300},
		{name: "portrait offset 1.0 (bottom)", imgW: 300, imgH: 400, offset: 1.0, wantX0: 0, wantY0: 100, wantSide: 300},
		// Landscape 400x300: side=300, maxX=100, y0=0.
		{name: "landscape offset 0.0 (left)", imgW: 400, imgH: 300, offset: 0.0, wantX0: 0, wantY0: 0, wantSide: 300},
		{name: "landscape offset 0.5 (center)", imgW: 400, imgH: 300, offset: 0.5, wantX0: 50, wantY0: 0, wantSide: 300},
		{name: "landscape offset 1.0 (right)", imgW: 400, imgH: 300, offset: 1.0, wantX0: 100, wantY0: 0, wantSide: 300},
		// Clamping.
		{name: "offset clamped below 0", imgW: 300, imgH: 400, offset: -1.0, wantX0: 0, wantY0: 0, wantSide: 300},
		{name: "offset clamped above 1", imgW: 300, imgH: 400, offset: 2.0, wantX0: 0, wantY0: 100, wantSide: 300},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := squareCropRectWithOffset(tc.imgW, tc.imgH, tc.offset)
			if r.Min.X != tc.wantX0 {
				t.Errorf("x0: got %d, want %d", r.Min.X, tc.wantX0)
			}
			if r.Min.Y != tc.wantY0 {
				t.Errorf("y0: got %d, want %d", r.Min.Y, tc.wantY0)
			}
			if r.Dx() != tc.wantSide || r.Dy() != tc.wantSide {
				t.Errorf("size: got %dx%d, want %dx%d", r.Dx(), r.Dy(), tc.wantSide, tc.wantSide)
			}
		})
	}
}

func TestSquareCropRectAnchors(t *testing.T) {
	t.Parallel()

	// Landscape 1024×768: side = 768, x0 = (1024-768)/2 = 128.
	const w, h, side = 1024, 768, 768

	tests := []struct {
		anchor AnchorPosition
		wantY0 int
	}{
		{AnchorTop, 0},
		{AnchorCenter, (h - side) / 2},
		{AnchorBottom, h - side},
	}

	for _, tc := range tests {
		t.Run(string(tc.anchor), func(t *testing.T) {
			t.Parallel()

			r := squareCropRect(w, h, tc.anchor)
			if r.Min.Y != tc.wantY0 {
				t.Errorf("y0: got %d, want %d", r.Min.Y, tc.wantY0)
			}
			if r.Dx() != side || r.Dy() != side {
				t.Errorf("rect size: got %dx%d, want %dx%d", r.Dx(), r.Dy(), side, side)
			}
		})
	}
}

func TestSquareCropRectPortrait(t *testing.T) {
	t.Parallel()

	// Portrait 768×1024: side = 768, vertically centered x0 = 0 (already full width).
	r := squareCropRect(768, 1024, AnchorCenter)
	if r.Dx() != 768 || r.Dy() != 768 {
		t.Errorf("rect size: got %dx%d, want 768x768", r.Dx(), r.Dy())
	}
	// X should be centered: x0 = (768-768)/2 = 0.
	if r.Min.X != 0 {
		t.Errorf("x0: got %d, want 0", r.Min.X)
	}
}

func TestCropToRectOrigin(t *testing.T) {
	t.Parallel()

	src := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	rect := image.Rect(10, 20, 60, 70)
	got := cropToRect(src, rect)

	if got.Bounds().Min != (image.Point{}) {
		t.Errorf("Min: got %v, want (0,0)", got.Bounds().Min)
	}
	if got.Bounds().Dx() != 50 || got.Bounds().Dy() != 50 {
		t.Errorf("size: got %dx%d, want 50x50", got.Bounds().Dx(), got.Bounds().Dy())
	}
}
