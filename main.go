// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package main contains the main entry point for the application.
package main

import (
	_ "embed"

	"github.com/mikefero/deslacked/internal/cmd"
)

// license contains the embedded MIT license text.
//
//go:embed LICENSE
var license string

func main() {
	cmd.Execute(cmd.Options{
		License: license,
	})
}
