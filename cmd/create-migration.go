package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Juksefantomet/gecho/internal/migrate"
)

var createMigrationCmd = &cobra.Command{
	Use:   "create-migration <name>",
	Short: "Create empty up/down SQL migration files",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := migrate.Create(name); err != nil {
			fmt.Fprintf(os.Stderr, "✗ Migration creation failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createMigrationCmd)
}
