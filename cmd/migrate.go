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
		_ = godotenv.Load(".env")
		database.InitDB()
		db := database.GetDB()

		if len(args) == 0 {
			fmt.Println("Running migrations...")
			database.RunMigrations(db)
			fmt.Println("Migrations complete.")
			return
		}

		switch args[0] {
		case "down":
			fmt.Println("Rolling back last migration...")
			database.RollbackLastMigration(db)
			fmt.Println("Rollback complete.")
		case "help":
			printMigrateHelp()
		default:
			printMigrateHelp()
			fmt.Fprintf(os.Stderr, "\n✗ Unknown argument: %s\n", args[0])
			os.Exit(1)
		}
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
}
