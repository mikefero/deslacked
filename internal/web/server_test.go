// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package web

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikefero/deslacked/internal/transform"
)

// imageRequest builds a multipart POST request with an "image" file field containing data.
func imageRequest(t *testing.T, filename string, data []byte) *http.Request {
	t.Helper()

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	part, err := mw.CreateFormFile("image", filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err = part.Write(data); err != nil {
		t.Fatalf("write form file data: %v", err)
	}
	if err = mw.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/process", body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return req
}

// minValidPNG returns the bytes of a 128x128 solid-gray PNG, the smallest image the
// transform pipeline accepts without error.
func minValidPNG(t *testing.T) []byte {
	t.Helper()

	img := image.NewNRGBA(image.Rect(0, 0, 128, 128))
	for y := range 128 {
		for x := range 128 {
			img.SetNRGBA(x, y, color.NRGBA{R: 100, G: 100, B: 100, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode PNG: %v", err)
	}
	return buf.Bytes()
}

func TestHandleProcessWrongMethod(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/api/process", nil)
	w := httptest.NewRecorder()
	handleProcess(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestHandleProcessMissingImageField(t *testing.T) {
	t.Parallel()

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	if err := mw.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/process", body)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	w := httptest.NewRecorder()
	handleProcess(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleProcessInvalidImageData(t *testing.T) {
	t.Parallel()

	req := imageRequest(t, "bad.png", []byte("not an image"))
	w := httptest.NewRecorder()
	handleProcess(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusUnprocessableEntity)
	}
}

func TestHandleProcessBodyTooLarge(t *testing.T) {
	t.Parallel()

	// Send more bytes than MaxBytesReader allows to trigger the 413 path.
	oversized := make([]byte, maxUploadBytes+formParseBuffer+1)
	req := imageRequest(t, "huge.png", oversized)

	w := httptest.NewRecorder()
	handleProcess(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestHandleProcessValidPNG(t *testing.T) {
	t.Parallel()

	req := imageRequest(t, "valid.png", minValidPNG(t))
	w := httptest.NewRecorder()
	handleProcess(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "image/png" {
		t.Errorf("content-type: got %q, want %q", ct, "image/png")
	}
}

func TestHandleProcessWithOffset(t *testing.T) {
	t.Parallel()

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	part, err := mw.CreateFormFile("image", "valid.png")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err = part.Write(minValidPNG(t)); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	if err = mw.WriteField("offset", "0.5"); err != nil {
		t.Fatalf("write offset field: %v", err)
	}
	if err = mw.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/process", body)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	w := httptest.NewRecorder()
	handleProcess(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestBuildOptsValidOffset(t *testing.T) {
	t.Parallel()

	opts := buildOpts("0.5")
	if opts.CropOffset == nil {
		t.Fatal("expected CropOffset to be set, got nil")
	}
	if *opts.CropOffset != 0.5 {
		t.Errorf("CropOffset: got %f, want 0.5", *opts.CropOffset)
	}
}

func TestBuildOptsInvalidOffset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		offset string
	}{
		{name: "empty string", offset: ""},
		{name: "out of range high", offset: "1.5"},
		{name: "out of range low", offset: "-0.1"},
		{name: "non-numeric", offset: "abc"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := buildOpts(tc.offset)
			if opts.CropOffset != nil {
				t.Errorf("expected nil CropOffset for offset %q, got %f", tc.offset, *opts.CropOffset)
			}
			if opts.Anchor != transform.AnchorCenter {
				t.Errorf("expected AnchorCenter for invalid offset %q, got %q", tc.offset, opts.Anchor)
			}
		})
	}
}
