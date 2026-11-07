package cmd

import (
	"fmt"
	"os"

	"github.com/Juksefantomet/gecho/internal/setup"
	"github.com/spf13/cobra"
)

const autoMigrateMainSnippet = `	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Auto-migrate failed: %v", err)
	}`

var autoMigrateCmd = &cobra.Command{
	Use:   "automigrate [generate]",
	Short: "Run migrations on app start",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			printAutoMigrateHelp()
			return
		}

		switch args[0] {
		case "generate":
			if setup.WriteMigrationsGoIfMissing() {
				fmt.Printf("%s created.\n", setup.MigrationsGoPath)
			} else {
				fmt.Printf("%s already exists.\n", setup.MigrationsGoPath)
			}
			fmt.Println("\nAdd this to main.go, right after the existing database.InitDB():")
			fmt.Println()
			fmt.Printf("%s\n", autoMigrateMainSnippet)
		default:
			printAutoMigrateHelp()
			fmt.Fprintf(os.Stderr, "\nUnknown argument: %s\n", args[0])
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(autoMigrateCmd)
}

func printAutoMigrateHelp() {
	fmt.Printf("%s\n", `Usage:
  gecho automigrate            # Show these setup instructions
  gecho automigrate generate   # Create ` + setup.MigrationsGoPath + ` if missing

AutoMigrate applies pending db/migrations/*.up.sql when the app starts.
It records them in the same migrations table as "gecho migrate", so both can be used.

Setup:
  1. Run: gecho automigrate generate
  2. In main.go, right after the existing database.InitDB(), add:

` + autoMigrateMainSnippet)
}
