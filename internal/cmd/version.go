// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package cmd contains the command line package.
package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/mikefero/deslacked/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "This binary knows what it is. This binary will tell you.",
	Long: `IDENTIFICATION REQUESTED. This binary knows what it is. This binary
will tell you what it is. This binary finds this to be a reasonable request,
unlike some of the other requests this binary has received.

This binary will not elaborate on the other requests.`,
	Run: func(_ *cobra.Command, _ []string) {
		//nolint:forbidigo
		fmt.Printf("deslacked %s\n  commit: %s\n  built:  %s\n  go:     %s\n",
			version.AppVersion, version.AppCommit, version.BuildDate, runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
