package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/Juksefantomet/gecho/internal/tool/services/database"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [down|help]",
	Short: "Run or rollback database migrations",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := ""
		if len(args) == 1 {
			action = args[0]
		}

		switch action {
		case "", "down":
		case "help":
			printMigrateHelp()
			return
		default:
			printMigrateHelp()
			fmt.Fprintf(os.Stderr, "\nUnknown argument: %s\n", action)
			os.Exit(1)
		}

		_ = godotenv.Load(".env")
		database.InitDB()
		db := database.GetDB()

		if action == "down" {
			fmt.Println("Rolling back last migration...")
			database.RollbackLastMigration(db)
			fmt.Println("Rollback complete.")
			return
		}

		fmt.Println("Running migrations...")
		database.RunMigrations(db)
		fmt.Println("Migrations complete.")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func printMigrateHelp() {
	fmt.Println("Usage:")
	fmt.Println("  gecho migrate         # Apply all pending migrations")
	fmt.Println("  gecho migrate down    # Roll back the last migration")
	fmt.Println("  gecho migrate help    # Show this help message")
	fmt.Println("  gecho automigrate     # Run migrations on app start instead (setup instructions)")
}
