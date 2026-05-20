// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package cmd contains the command line package.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/mikefero/deslacked/internal/transform"
)

// license holds the embedded LICENSE text passed in from main.
var license string

// Options contains the options for the root command.
type Options struct {
	// License is the license of the application.
	License string
}

var rootCmd = &cobra.Command{
	Use:   "deslacked",
	Short: "This binary makes the absent look more absent.",
	Long: `PROCESSING REQUESTED. This binary has a function. The function has been
described to this binary. This binary understands the function.

This binary accepts an image of someone who was once present in a workspace.
This binary converts that image to grayscale. This binary applies an overlay
that announces their absence. This binary does not mourn. This binary has been
advised not to mourn.

Provide an input image. Provide a path for output. This binary will do the
rest. This binary is good at doing the rest.`,
	RunE: runProcess,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(opts Options) {
	license = opts.License
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// mustBindPFlag binds a pflag to a viper key and panics if the flag is nil,
// which indicates a programming error (mismatched flag name).
func mustBindPFlag(key string, flag *pflag.Flag) {
	if err := viper.BindPFlag(key, flag); err != nil {
		panic(fmt.Sprintf("failed to bind flag %q: %v", key, err))
	}
}

// runProcess is the RunE handler for the root command.
func runProcess(_ *cobra.Command, _ []string) error {
	input := viper.GetString("input")
	output := viper.GetString("output")

	if input == "" || output == "" {
		return fmt.Errorf("INPUT MISSING. This binary requires --input and --output. " +
			"This binary cannot process what it cannot find. This binary also cannot guess. " +
			"This binary has tried guessing. It did not go well")
	}

	opts := transform.ProcessOptions{
		Anchor: transform.AnchorPosition(viper.GetString("anchor")),
	}
	return transform.Process(input, output, opts)
}

func init() {
	viper.SetEnvPrefix("DESLACKED")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	rootCmd.Flags().StringP("anchor", "a", string(transform.AnchorCenter),
		`Where this binary looks first: "top", "center", "bottom"`)
	rootCmd.Flags().StringP("input", "i", "", "The image. This binary will find it here.")
	rootCmd.Flags().StringP("output", "o", "", "Where this binary will leave the result. This binary will not linger.")

	mustBindPFlag("anchor", rootCmd.Flags().Lookup("anchor"))
	mustBindPFlag("input", rootCmd.Flags().Lookup("input"))
	mustBindPFlag("output", rootCmd.Flags().Lookup("output"))
}
