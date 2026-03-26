// Package cli defines all cobra commands for wsm.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// Version is set at build time via ldflags.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "wsm",
	Short: "Workspace manager for multi-repo projects",
	Long: `wsm manages a workspace of multiple git repositories defined in a
workspace.json manifest. It can clone repos, show their status, pull
updates, and manage the workspace configuration.

Use "wsm <command> --help" for details on any command.`,
	Version: Version,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "workspace.json", "path to workspace manifest file")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
