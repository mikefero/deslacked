// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package web provides the built-in HTTP server for deslacked.
package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	// Imported for side-effect image format registration.
	_ "image/gif"
	_ "image/jpeg"

	"github.com/mikefero/deslacked/internal/transform"
	"github.com/mikefero/deslacked/internal/version"
)

//go:embed static
var staticFiles embed.FS

const (
	// maxUploadBytes is the maximum accepted request body size for image uploads (10 MB).
	maxUploadBytes = 10 * 1024 * 1024
	// formParseBuffer is extra headroom beyond maxUploadBytes for non-file form fields.
	formParseBuffer = 1024
	// browserOpenDelay is how long the server waits before opening the browser, giving
	// the listener time to bind before the URL is visited.
	browserOpenDelay = 150 * time.Millisecond
)

// Serve starts the HTTP server on the given port. If openBrowser is true, the
// default browser is opened after a short delay.
func Serve(port int, openBrowser bool) error {
	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("setting up static files: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(sub)))
	mux.HandleFunc("/api/version", handleVersion)
	mux.HandleFunc("/api/process", handleProcess)

	addr := fmt.Sprintf(":%d", port)
	url := fmt.Sprintf("http://localhost:%d", port)

	//nolint:errcheck
	fmt.Fprintf(os.Stdout, "BINARY ONLINE. This binary has made itself available at %s\n", url)
	//nolint:errcheck
	fmt.Fprintln(os.Stdout,
		"This binary is ready. This binary is watching. Press Ctrl+C to stop. This binary will notice.")

	if openBrowser {
		go func() {
			time.Sleep(browserOpenDelay)
			_ = openURL(url)
		}()
	}

	//nolint:gosec
	return http.ListenAndServe(addr, mux)
}

// handleVersion returns the application version as JSON.
func handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed,
			"INVALID REQUEST. This binary only accepts GET requests at this endpoint. "+
				"This binary knows what it is. This binary will tell you what it is. "+
				"This binary will not tell you in any other way than this.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"version": version.AppVersion})
}

// handleProcess accepts a multipart POST with fields "image" (file) and "offset"
// (float string 0.0–1.0), runs the deslacked transform in memory, and responds
// with the PNG result.
func handleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed,
			"INVALID REQUEST. This binary only accepts POST requests at this endpoint. "+
				"This binary does not know why you sent something else. This binary suspects you know.")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes+formParseBuffer)
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		writeError(w, http.StatusRequestEntityTooLarge,
			"FILE REJECTED. This binary does not accept files larger than 10 MB. "+
				"This binary has standards. This binary would like to keep them.")
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		writeError(w, http.StatusBadRequest,
			"MISSING ATTACHMENT. This binary expected an image. This binary received a request without one. "+
				"This binary is not sure what you expected to happen.")
		return
	}
	defer func() { _ = file.Close() }()

	img, _, err := image.Decode(file)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity,
			"DECODE FAILURE. This binary could not interpret that file as an image. "+
				"PNG, JPEG, and GIF are accepted. This binary tried. It did not work.")
		return
	}

	opts := buildOpts(r.FormValue("offset"))

	result, err := transform.ProcessImage(img, opts)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity,
			"PROCESSING FAILURE. This binary attempted the transformation. "+
				"This binary encountered an obstacle: "+err.Error()+". This binary is not pleased.")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store")
	if err = png.Encode(w, result); err != nil {
		// Headers already sent; nothing more we can do.
		return
	}
}

// buildOpts constructs ProcessOptions from the "offset" form value. A valid float
// in [0, 1] uses continuous offset positioning; anything else defaults to center.
func buildOpts(offsetStr string) transform.ProcessOptions {
	if offsetStr != "" {
		if v, err := strconv.ParseFloat(offsetStr, 64); err == nil && v >= 0 && v <= 1 {
			return transform.ProcessOptions{CropOffset: &v}
		}
	}
	return transform.DefaultProcessOptions()
}

// writeError sends a JSON error response in the DCC tone.
func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// openURL launches the given URL in the default system browser.
func openURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		//nolint:gosec
		cmd = exec.Command("open", url)
	case "windows":
		//nolint:gosec
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		//nolint:gosec
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
