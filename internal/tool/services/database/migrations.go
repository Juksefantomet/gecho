package database

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Migration represents a migration record in the database
type Migration struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique"`
	AppliedAt time.Time
}

// RunMigrations applies all pending .up.sql files from db/migrations
func RunMigrations(db *gorm.DB) {
	migrationDir := "db/migrations"

	if err := db.AutoMigrate(&Migration{}); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	var applied []Migration
	db.Find(&applied)

	appliedMap := make(map[string]bool)
	for _, m := range applied {
		appliedMap[m.Name] = true
	}

	files, err := filepath.Glob(filepath.Join(migrationDir, "*.up.sql"))
	if err != nil {
		log.Fatalf("Failed to read migration files: %v", err)
	}

	for _, file := range files {
		name := filepath.Base(file)
		if appliedMap[name] {
			continue
		}

		sql, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read %s: %v", file, err)
		}

		log.Printf("Applying migration: %s", name)
		if err := db.Exec(string(sql)).Error; err != nil {
			log.Fatalf("Migration %s failed: %v", name, err)
		}

		db.Create(&Migration{Name: name, AppliedAt: time.Now()})
	}

	log.Println("All migrations applied.")
}

// RollbackLastMigration undoes the most recently applied migration (if a .down.sql exists)
func RollbackLastMigration(db *gorm.DB) {
	var last Migration
	if err := db.Order("applied_at desc").First(&last).Error; err != nil {
		log.Println("No applied migrations to roll back.")
		return
	}

	downFile := filepath.Join("db/migrations", replaceSuffix(last.Name, ".up.sql", ".down.sql"))
	if _, err := os.Stat(downFile); os.IsNotExist(err) {
		log.Fatalf("Missing .down.sql for %s", last.Name)
	}

	sql, err := os.ReadFile(downFile)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", downFile, err)
	}

	log.Printf("Rolling back: %s", last.Name)
	if err := db.Exec(string(sql)).Error; err != nil {
		log.Fatalf("Rollback failed for %s: %v", downFile, err)
	}

	db.Delete(&last)
	log.Printf("Rollback complete for %s.", last.Name)
}

func replaceSuffix(name, old, new string) string {
	if !strings.HasSuffix(name, old) {
		log.Fatalf("Invalid filename format: %s", name)
	}
	return strings.TrimSuffix(name, old) + new
}
