package transform

import (
	"image"
	"image/color"
	"testing"
)

func TestToGrayscale(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		pixel    color.NRGBA
		wantLuma uint8
		wantA    uint8
	}{
		{
			name:     "solid red",
			pixel:    color.NRGBA{R: 255, G: 0, B: 0, A: 255},
			wantLuma: color.GrayModel.Convert(color.NRGBA{R: 255, G: 0, B: 0, A: 255}).(color.Gray).Y,
			wantA:    255,
		},
		{
			name:     "solid green",
			pixel:    color.NRGBA{R: 0, G: 255, B: 0, A: 255},
			wantLuma: color.GrayModel.Convert(color.NRGBA{R: 0, G: 255, B: 0, A: 255}).(color.Gray).Y,
			wantA:    255,
		},
		{
			name:     "solid blue",
			pixel:    color.NRGBA{R: 0, G: 0, B: 255, A: 255},
			wantLuma: color.GrayModel.Convert(color.NRGBA{R: 0, G: 0, B: 255, A: 255}).(color.Gray).Y,
			wantA:    255,
		},
		{
			name:     "fully transparent preserves alpha",
			pixel:    color.NRGBA{R: 200, G: 100, B: 50, A: 0},
			wantLuma: 0,
			wantA:    0,
		},
		{
			name:     "already gray is idempotent",
			pixel:    color.NRGBA{R: 128, G: 128, B: 128, A: 255},
			wantLuma: 128,
			wantA:    255,
		},
		{
			name:     "mixed RGBA produces correct luma",
			pixel:    color.NRGBA{R: 100, G: 150, B: 200, A: 255},
			wantLuma: color.GrayModel.Convert(color.NRGBA{R: 100, G: 150, B: 200, A: 255}).(color.Gray).Y,
			wantA:    255,
		},
		{
			name:     "white pixel produces max luma",
			pixel:    color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			wantLuma: 255,
			wantA:    255,
		},
		{
			name:     "black pixel produces zero luma",
			pixel:    color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			wantLuma: 0,
			wantA:    255,
		},
		{
			name:     "semi-transparent pixel preserves alpha",
			pixel:    color.NRGBA{R: 200, G: 100, B: 50, A: 128},
			wantLuma: color.GrayModel.Convert(color.NRGBA{R: 200, G: 100, B: 50, A: 128}).(color.Gray).Y,
			wantA:    128,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			src := image.NewNRGBA(image.Rect(0, 0, 1, 1))
			src.SetNRGBA(0, 0, tc.pixel)

			got := toGrayscale(src)

			p := got.NRGBAAt(0, 0)
			if p.R != tc.wantLuma || p.G != tc.wantLuma || p.B != tc.wantLuma {
				t.Errorf("luma: got (%d,%d,%d), want %d for all channels", p.R, p.G, p.B, tc.wantLuma)
			}
			if p.A != tc.wantA {
				t.Errorf("alpha: got %d, want %d", p.A, tc.wantA)
			}
		})
	}
}

func TestToGrayscalePreservesDimensions(t *testing.T) {
	t.Parallel()

	src := image.NewNRGBA(image.Rect(0, 0, 32, 64))
	got := toGrayscale(src)

	if got.Bounds() != src.Bounds() {
		t.Errorf("bounds changed: got %v, want %v", got.Bounds(), src.Bounds())
	}
}
