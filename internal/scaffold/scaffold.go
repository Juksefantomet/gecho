package scaffold

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/Juksefantomet/gecho/internal/migrate"
	"github.com/jinzhu/inflection"
)

func Run(rawName string) error {
	moduleName, err := getModuleName("go.mod")
	if err != nil {
		return err
	}

	input := toSnakeCase(strings.ReplaceAll(strings.TrimSpace(rawName), " ", "_"))
	if input == "" {
		return fmt.Errorf("missing model name")
	}

	singularName := inflection.Singular(input)
	tableName := inflection.Plural(singularName)
	if singularName != input {
		fmt.Printf("Plural version of the model detected, using singularized version: %s\n", singularName)
	}

	structName := toPascalCase(singularName)
	pluralStruct := toPascalCase(tableName)

	timestamp := time.Now().Format("20060102150405")
	migrationDir := "db/migrations"

	if existing, ok := migrate.FindExisting(migrationDir, input, singularName, tableName); ok {
		return fmt.Errorf("a migration for %s already exists: %s", tableName, existing)
	}
	modelsDir := "app/models"
	routesDir := "app/routes"
	queriesDir := "app/services/database"

	upFile := fmt.Sprintf("%s/%s_%s.up.sql", migrationDir, timestamp, input)
	downFile := fmt.Sprintf("%s/%s_%s.down.sql", migrationDir, timestamp, input)

	modelFile := fmt.Sprintf("%s/%s.go", modelsDir, toLowerCamelCase(structName))
	routeFile := fmt.Sprintf("%s/%s.go", routesDir, toLowerCamelCase(pluralStruct))
	queriesFile := fmt.Sprintf("%s/%sQueries.go", queriesDir, toLowerCamelCase(structName))

	for _, dir := range []string{migrationDir, modelsDir, routesDir, queriesDir} {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	if err := writeFile(upFile, fmt.Sprintf(`-- SQL UP Migration

CREATE TABLE %s
(
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`, tableName)); err != nil {
		return err
	}
	if err := writeFile(downFile, fmt.Sprintf(`-- SQL DOWN Migration

DROP TABLE IF EXISTS %s;
`, tableName)); err != nil {
		return err
	}

	if err := writeFileIfNotExists(modelFile, modelBoilerplate(structName)); err != nil {
		return err
	}
	if err := writeFileIfNotExists(routeFile, routeBoilerplate(structName, pluralStruct, tableName, moduleName)); err != nil {
		return err
	}
	if err := writeFileIfNotExists(queriesFile, queriesBoilerplate(structName, pluralStruct, tableName, moduleName)); err != nil {
		return err
	}

	fmt.Printf("Scaffold created:\n  %s\n  %s\n  %s\n  %s\n  %s\n",
		upFile, downFile, modelFile, routeFile, queriesFile)
	return nil
}

func getModuleName(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", path, err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}
	return "", fmt.Errorf("no module declaration found in %s", path)
}

func toLowerCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// toSnakeCase turns BlogPost / HTTPServer / blog_post into blog_post / http_server / blog_post.
func toSnakeCase(s string) string {
	runes := []rune(s)
	var result []rune
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) && runes[i-1] != '_' {
			prevLower := unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1])
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if prevLower || nextLower {
				result = append(result, '_')
			}
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

func writeFile(path string, content string) error {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}

func writeFileIfNotExists(path string, content string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return writeFile(path, content)
	}
	return nil
}

func modelBoilerplate(structName string) string {
	return fmt.Sprintf(`package models

import "time"

type %s struct {
	ID        int       `+"`json:\"id\"`"+`
	Name      string    `+"`json:\"name\"`"+`
	// populate the rest of the model here
	// Example  string    `+"`json:\"example\"`"+`

	CreatedAt time.Time `+"`gorm:\"autoCreateTime\"`"+`
	UpdatedAt time.Time `+"`gorm:\"autoUpdateTime\"`"+`
}
`, structName)
}

func queriesBoilerplate(structName, pluralStruct, tableName, moduleName string) string {
	funcName := "Get" + pluralStruct

	return fmt.Sprintf(`package database

import (
	"%[3]s/app/models"
)

// %[1]s returns all %[2]s
func %[1]s() ([]models.%[4]s, error) {
	var items []models.%[4]s
	err := GetDB().Raw("SELECT * FROM %[2]s").Scan(&items).Error
	return items, err
}
`, funcName, tableName, moduleName, structName)
}

func routeBoilerplate(structName, pluralStruct, tableName, moduleName string) string {
	funcName := "Get" + pluralStruct
	responseStruct := structName + "Response"

	return fmt.Sprintf(`package routes

import (
	"net/http"
	"%[4]s/app/models"
	"%[4]s/app/services/database"

	"github.com/labstack/echo/v4"
)

// %[2]s defines the response structure
type %[2]s struct {
	%[5]s []models.%[1]s `+"`json:\"%[6]s\"`"+`
}

// @Summary %[6]s
// @Description Returns all %[6]s
// @Tags %[6]s
// @Accept json
// @Produce json
// @Success 200 {object} %[2]s
// @Router /%[6]s [get]
func %[3]s(c echo.Context) error {
	items, err := database.%[3]s()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch %[6]s"})
	}

	return c.JSON(http.StatusOK, %[2]s{%[5]s: items})
}
`, structName, responseStruct, funcName, moduleName, pluralStruct, tableName)
}
