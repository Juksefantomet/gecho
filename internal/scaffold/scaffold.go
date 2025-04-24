package scaffold

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func Run(rawName string) error {
	moduleName := getModuleName("go.mod")
	if rawName == "" {
		return fmt.Errorf("missing model name")
	}

	timestamp := time.Now().Format("20060102150405")
	migrationDir := "db/migrations"
	modelsDir := "app/models"
	routesDir := "app/routes"
	queriesDir := "app/services/database"

	snakeName := strings.ToLower(strings.ReplaceAll(rawName, " ", "_"))
	structName := toPascalCase(singular(snakeName))
	pluralName := plural(snakeName)
	lowerCamel := toLowerCamelCase(structName)

	upFile := fmt.Sprintf("%s/%s_%s.up.sql", migrationDir, timestamp, snakeName)
	downFile := fmt.Sprintf("%s/%s_%s.down.sql", migrationDir, timestamp, snakeName)

	modelFile := fmt.Sprintf("%s/%s.go", modelsDir, lowerCamel)
	routeFile := fmt.Sprintf("%s/%ss.go", routesDir, lowerCamel)
	queriesFile := fmt.Sprintf("%s/%sQueries.go", queriesDir, lowerCamel)

	for _, dir := range []string{migrationDir, modelsDir, routesDir, queriesDir} {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create %s: %v", dir, err)
		}
	}

	writeFile(upFile, fmt.Sprintf(`-- SQL UP Migration

CREATE TABLE %s
(
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`, pluralName))
	writeFile(downFile, fmt.Sprintf(`-- SQL DOWN Migration

DROP TABLE IF EXISTS %s;
`, pluralName))

	writeFileIfNotExists(modelFile, modelBoilerplate(structName))
	writeFileIfNotExists(routeFile, routeBoilerplate(structName, moduleName))
	writeFileIfNotExists(queriesFile, queriesBoilerplate(structName, moduleName))

	fmt.Printf("Scaffold created:\n  %s\n  %s\n  %s\n  %s\n, %s\n",
		upFile, downFile, modelFile, routeFile, queriesFile)
return nil

}

func getModuleName(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", path, err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}
	log.Fatal("No module declaration found in go.mod")
	return ""
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

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, rune(strings.ToLower(string(r))[0]))
	}
	return string(result)
}

func writeFile(path string, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		log.Fatalf("Failed to write %s: %v", path, err)
	}
}

func writeFileIfNotExists(path string, content string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		writeFile(path, content)
	}
}

func plural(name string) string {
	if strings.HasSuffix(name, "y") && !strings.HasSuffix(name, "ay") && !strings.HasSuffix(name, "ey") {
		return strings.TrimSuffix(name, "y") + "ies"
	}
	if strings.HasSuffix(name, "s") {
		return name + "es"
	}
	return name + "s"
}

func singular(name string) string {
	if strings.HasSuffix(name, "ies") {
		return strings.TrimSuffix(name, "ies") + "y"
	}
	if strings.HasSuffix(name, "es") {
		return strings.TrimSuffix(name, "es")
	}
	if strings.HasSuffix(name, "s") {
		return strings.TrimSuffix(name, "s")
	}
	return name
}

func modelBoilerplate(structName string) string {
	return fmt.Sprintf(`package models

import "time"

type %s struct {
	ID        int       `+"`json:\"id\"`"+`
	// populate the rest of the model here
	// Example  string    `+"`json:\"example\"`"+`

	CreatedAt time.Time `+"`gorm:\"autoCreateTime\"`"+`
	UpdatedAt time.Time `+"`gorm:\"autoUpdateTime\"`"+`
}
`, structName)
}

func queriesBoilerplate(structName, moduleName string) string {
	funcName := "Get" + structName + "s"
	tableName := toSnakeCase(structName) + "s"

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

func routeBoilerplate(structName, moduleName string) string {
	jsonTag := toSnakeCase(structName) + "s"
	funcName := "Get" + structName + "s"
	responseStruct := structName + "Response"
	responseField := structName + "s"

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
`, structName, responseStruct, funcName, moduleName, responseField, jsonTag)
}
