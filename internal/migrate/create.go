package migrate

import (
	"fmt"
	"os"
	"time"
)

func Create(name string) error {
	if name == "" {
		return fmt.Errorf("missing migration name")
	}

	timestamp := time.Now().Format("20060102150405")
	migrationDir := "db/migrations"

	upFile := fmt.Sprintf("%s/%s_%s.up.sql", migrationDir, timestamp, name)
	downFile := fmt.Sprintf("%s/%s_%s.down.sql", migrationDir, timestamp, name)

	if err := os.MkdirAll(migrationDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	if err := os.WriteFile(upFile, []byte("-- SQL UP Migration"), 0644); err != nil {
		return fmt.Errorf("failed to create %s: %w", upFile, err)
	}

	if err := os.WriteFile(downFile, []byte("-- SQL DOWN Migration"), 0644); err != nil {
		return fmt.Errorf("failed to create %s: %w", downFile, err)
	}

	fmt.Printf("Migration created:\n  %s\n  %s\n", upFile, downFile)
	return nil
}
