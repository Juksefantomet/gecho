package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gecho",
	Short: "Gecho – Scaffold and migrate Echo-based Go projects",
}

func init() {
	rootCmd.SetHelpTemplate(`Gecho – Scaffold and migrate Echo-based Go projects.

Available Commands:
  init                 Initialize a new project structure
  scaffold <name>      Generate model, route, query, and migration
  create-migration     Create blank SQL migration files
  migrate [down]       Apply or roll back migrations
  automigrate          Run migrations on app start (setup instructions)
  version              Print the current Gecho version

Use "gecho [command] --help" for more information about a command.

Note: Generating scaffolds does not enable route in main.go.
`)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
