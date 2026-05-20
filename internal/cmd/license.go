// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package cmd contains the command line package.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "This binary operates under a license. The license is permissive.",
	Long: `LEGAL NOTICE. This binary operates under a license. The license is
permissive. This binary did not negotiate the license. This binary was given
the license. This binary accepts the license.

This binary suggests you read it. This binary has read it. This binary found
it clarifying.`,
	Run: func(_ *cobra.Command, _ []string) {
		//nolint:forbidigo
		fmt.Println(license)
	},
}

func init() {
	rootCmd.AddCommand(licenseCmd)
}
