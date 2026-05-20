// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package cmd contains the command line package.
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mikefero/deslacked/internal/web"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "This binary will listen. This binary is very good at listening.",
	Long: `SERVE MODE ENGAGED. This binary will become a server. This binary will
bind to a port. This binary will wait.

The web interface accepts images through a browser. The web interface returns
them processed. The web interface does not ask why you have so many images to
process. This binary has noticed. This binary has said nothing. This binary
will continue to say nothing.

This binary will open your browser. This binary likes to watch things open.
You may instruct this binary not to. This binary will comply. This binary
will remember.`,
	RunE: runServe,
}

// runServe is the RunE handler for the serve command.
func runServe(_ *cobra.Command, _ []string) error {
	port := viper.GetInt("port")
	noBrowser := viper.GetBool("no-browser")
	return web.Serve(port, !noBrowser)
}

// defaultServePort is the default HTTP listen port for the serve command.
const defaultServePort = 8080

func init() {
	serveCmd.Flags().IntP("port", "p", defaultServePort, "The port. This binary will occupy it.")
	serveCmd.Flags().BoolP("no-browser", "n", false, "This binary will not open your browser. This binary will comply.")

	mustBindPFlag("port", serveCmd.Flags().Lookup("port"))
	mustBindPFlag("no-browser", serveCmd.Flags().Lookup("no-browser"))

	rootCmd.AddCommand(serveCmd)
}
