// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package assets provides embedded resources for deslacked.
package assets

import _ "embed"

// BarTemplateBytes contains the embedded deactivated-account bar PNG; 1024×144.
// Extracted from a Slack profile screenshot with rounded corners repaired to a full rectangle.
//
//go:embed bar-template.png
var BarTemplateBytes []byte
