package cmd

import (
	"fmt"
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
  version			   Print the current Gecho version

Use "gecho [command] --help" for more information about a command.

Note: Generating scaffolds does not enable route in main.go.
`)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(scaffoldCmd)
	rootCmd.AddCommand(createMigrationCmd)
	rootCmd.AddCommand(migrateCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
