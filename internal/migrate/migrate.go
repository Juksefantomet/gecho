package migrate

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/Juksefantomet/gecho/internal/tool/services/database"
)

func Run(action string) error {
	_ = godotenv.Load(".env")

	database.InitDB()
	db := database.GetDB()

	switch action {
	case "":
		log.Println("Running migrations on database...")
		database.RunMigrations(db)
		log.Println("Migrations complete.")
	case "down":
		log.Println("Rolling back last migration...")
		database.RollbackLastMigration(db)
		log.Println("Rollback complete.")
	case "help":
		printHelp()
	default:
		printHelp()
		return fmt.Errorf("unknown action: %s", action)
	}

	return nil
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  gecho migrate         # Apply all pending migrations")
	fmt.Println("  gecho migrate down    # Roll back the last migration (if .down.sql exists)")
	fmt.Println("  gecho migrate help    # Show this help message")
}
