package transform

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// writePNG encodes a synthetic image as a PNG temp file and returns its path.
func writePNG(t *testing.T, dir string, img image.Image) string {
	t.Helper()

	path := filepath.Join(dir, "input.png")
	//nolint:gosec
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create temp PNG: %v", err)
	}

	if err = png.Encode(f, img); err != nil {
		_ = f.Close()
		t.Fatalf("encode temp PNG: %v", err)
	}
	if err = f.Close(); err != nil {
		t.Fatalf("close temp PNG: %v", err)
	}
	return path
}

// solidImage creates a w×h NRGBA image filled with c.
func solidImage(w, h int, c color.NRGBA) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.SetNRGBA(x, y, c)
		}
	}
	return img
}

// writeJPEG encodes a synthetic image as a JPEG temp file and returns its path.
func writeJPEG(t *testing.T, dir string, img image.Image) string {
	t.Helper()

	path := filepath.Join(dir, "input.jpg")
	//nolint:gosec
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create temp JPEG: %v", err)
	}
	if err = jpeg.Encode(f, img, nil); err != nil {
		_ = f.Close()
		t.Fatalf("encode temp JPEG: %v", err)
	}
	if err = f.Close(); err != nil {
		t.Fatalf("close temp JPEG: %v", err)
	}
	return path
}

// writeGIF encodes a synthetic paletted image as a GIF temp file and returns its path.
func writeGIF(t *testing.T, dir string, w, h int) string {
	t.Helper()

	palette := color.Palette{
		color.RGBA{R: 100, G: 150, B: 200, A: 255},
		color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
	img := image.NewPaletted(image.Rect(0, 0, w, h), palette)

	path := filepath.Join(dir, "input.gif")
	//nolint:gosec
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create temp GIF: %v", err)
	}
	if err = gif.Encode(f, img, nil); err != nil {
		_ = f.Close()
		t.Fatalf("encode temp GIF: %v", err)
	}
	if err = f.Close(); err != nil {
		t.Fatalf("close temp GIF: %v", err)
	}
	return path
}

// outputDimensions reads the PNG at path and returns its width and height.
func outputDimensions(t *testing.T, path string) (int, int) {
	t.Helper()

	//nolint:gosec
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode output: %v", err)
	}
	b := img.Bounds()
	return b.Dx(), b.Dy()
}

func TestProcess(t *testing.T) {
	t.Parallel()

	offset0 := 0.0

	tests := []struct {
		name       string
		img        image.Image
		anchor     AnchorPosition
		cropOffset *float64
		wantErr    bool
		wantW      int
		wantH      int
	}{
		{
			name:   "512x512 input stays 512x512",
			img:    solidImage(512, 512, color.NRGBA{R: 100, G: 150, B: 200, A: 255}),
			anchor: AnchorCenter,
			wantW:  512, wantH: 512,
		},
		{
			name:   "128x128 input stays 128x128",
			img:    solidImage(128, 128, color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			anchor: AnchorCenter,
			wantW:  128, wantH: 128,
		},
		{
			name:   "200x200 input stays 200x200 (no upscale)",
			img:    solidImage(200, 200, color.NRGBA{R: 80, G: 80, B: 80, A: 255}),
			anchor: AnchorCenter,
			wantW:  200, wantH: 200,
		},
		{
			name:   "1024x1024 scaled to 512x512",
			img:    solidImage(1024, 1024, color.NRGBA{R: 200, G: 200, B: 200, A: 255}),
			anchor: AnchorCenter,
			wantW:  512, wantH: 512,
		},
		{
			name:   "non-square 800x600 center crop to 512x512",
			img:    solidImage(800, 600, color.NRGBA{R: 50, G: 100, B: 150, A: 255}),
			anchor: AnchorCenter,
			wantW:  512, wantH: 512,
		},
		{
			name:       "crop offset selects top of portrait image",
			img:        solidImage(300, 400, color.NRGBA{R: 50, G: 100, B: 150, A: 255}),
			cropOffset: &offset0,
			wantW:      300, wantH: 300,
		},
		{
			name:   "HD 1920x1080 cropped and scaled to 512x512",
			img:    solidImage(1920, 1080, color.NRGBA{R: 50, G: 100, B: 150, A: 255}),
			anchor: AnchorCenter,
			wantW:  512, wantH: 512,
		},
		{
			name:   "irregular 333x444 portrait cropped without scaling",
			img:    solidImage(333, 444, color.NRGBA{R: 80, G: 120, B: 160, A: 255}),
			anchor: AnchorCenter,
			wantW:  333, wantH: 333,
		},
		{
			name:    "64x64 input rejected",
			img:     solidImage(64, 64, color.NRGBA{R: 0, G: 0, B: 0, A: 255}),
			anchor:  AnchorCenter,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			inputPath := writePNG(t, dir, tc.img)
			outputPath := filepath.Join(dir, "output.png")

			opts := ProcessOptions{Anchor: tc.anchor, CropOffset: tc.cropOffset}
			err := Process(inputPath, outputPath, opts)

			if tc.wantErr {
				if err == nil {
					t.Error("expected error, got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			gotW, gotH := outputDimensions(t, outputPath)
			if gotW != tc.wantW || gotH != tc.wantH {
				t.Errorf("output size: got %dx%d, want %dx%d", gotW, gotH, tc.wantW, tc.wantH)
			}
		})
	}
}

func TestProcessMissingInput(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	err := Process(filepath.Join(dir, "nonexistent.png"), filepath.Join(dir, "out.png"), DefaultProcessOptions())
	if err == nil {
		t.Error("expected error for missing input, got none")
	}
}

func TestProcessInvalidFormat(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	inputPath := filepath.Join(dir, "bad.png")
	if err := os.WriteFile(inputPath, []byte("not an image"), 0o600); err != nil {
		t.Fatalf("write bad file: %v", err)
	}

	err := Process(inputPath, filepath.Join(dir, "out.png"), DefaultProcessOptions())
	if err == nil {
		t.Error("expected error for invalid image format, got none")
	}
}

func TestProcessOutputPathFailure(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	inputPath := writePNG(t, dir, solidImage(200, 200, color.NRGBA{R: 128, G: 128, B: 128, A: 255}))
	err := Process(inputPath, filepath.Join(dir, "nonexistent", "output.png"), DefaultProcessOptions())
	if err == nil {
		t.Error("expected error for unwritable output path, got none")
	}
}

func TestProcessJPEGInput(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	inputPath := writeJPEG(t, dir, solidImage(200, 200, color.NRGBA{R: 100, G: 150, B: 200, A: 255}))
	outputPath := filepath.Join(dir, "output.png")

	if err := Process(inputPath, outputPath, DefaultProcessOptions()); err != nil {
		t.Fatalf("unexpected error processing JPEG input: %v", err)
	}

	w, h := outputDimensions(t, outputPath)
	if w != 200 || h != 200 {
		t.Errorf("output size: got %dx%d, want 200x200", w, h)
	}
}

func TestProcessGIFInput(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	inputPath := writeGIF(t, dir, 200, 200)
	outputPath := filepath.Join(dir, "output.png")

	if err := Process(inputPath, outputPath, DefaultProcessOptions()); err != nil {
		t.Fatalf("unexpected error processing GIF input: %v", err)
	}

	w, h := outputDimensions(t, outputPath)
	if w != 200 || h != 200 {
		t.Errorf("output size: got %dx%d, want 200x200", w, h)
	}
}

func TestProcessImage(t *testing.T) {
	t.Parallel()

	img := solidImage(512, 512, color.NRGBA{R: 100, G: 150, B: 200, A: 255})
	result, err := ProcessImage(img, DefaultProcessOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b := result.Bounds()
	if b.Dx() != 512 || b.Dy() != 512 {
		t.Errorf("output size: got %dx%d, want 512x512", b.Dx(), b.Dy())
	}
}

func TestProcessImageWithCropOffset(t *testing.T) {
	t.Parallel()

	offset := 0.0
	img := solidImage(400, 300, color.NRGBA{R: 100, G: 150, B: 200, A: 255})
	result, err := ProcessImage(img, ProcessOptions{CropOffset: &offset})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b := result.Bounds()
	if b.Dx() != 300 || b.Dy() != 300 {
		t.Errorf("output size: got %dx%d, want 300x300", b.Dx(), b.Dy())
	}
}
