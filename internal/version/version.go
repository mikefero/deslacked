// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package version holds build-time version information injected via ldflags.
package version

// AppVersion is the application version string.
var AppVersion = "dev"

// AppCommit is the git commit hash at build time.
var AppCommit = "unknown"

// BuildDate is the UTC timestamp at build time.
var BuildDate = "unknown"
