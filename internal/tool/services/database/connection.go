package database

import (
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB reads env vars and initializes the GORM DB connection.
// Empty vars are left out so the driver falls back to its defaults.
func InitDB() {
	var err error

	parts := []string{"sslmode=disable"}
	for key, env := range map[string]string{
		"host":     "DATABASE_HOST",
		"user":     "DATABASE_USER",
		"password": "DATABASE_PASSWORD",
		"dbname":   "DATABASE_NAME",
		"port":     "DATABASE_PORT",
	} {
		if v := os.Getenv(env); v != "" {
			parts = append(parts, key+"="+v)
		}
	}

	db, err = gorm.Open(postgres.Open(strings.Join(parts, " ")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

// GetDB returns the active GORM DB connection
func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database is not initialized. Call InitDB() first.")
	}
	return db
}
