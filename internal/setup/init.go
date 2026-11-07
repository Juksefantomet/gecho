package setup

import (
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"strings"
)

var requiredDirs = []string{
	"app/models",
	"app/routes",
	"app/services/database",
	"db/migrations",
}

func RunInit() {
	ensureGoModExists()
	for _, dir := range requiredDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Failed to create %s: %v\n", dir, err)
			os.Exit(1)
		}
	}
	writeDotEnvIfMissing()
	writeMainGoIfMissing(getModuleName("go.mod"))
	writeDatabaseGoIfMissing()
	if WriteMigrationsGoIfMissing() {
		fmt.Println(MigrationsGoPath + " created.")
	}
	writeHelloWorldRouteIfMissing()
	fmt.Println("Gecho-compatible project structure initialized.")
	fmt.Println("Ready for scaffolding.")

	fmt.Println(`
Project ready!

Next steps:
  go install github.com/swaggo/swag/cmd/swag@v1.16.4
                                              # Install Swagger binary
  gecho scaffold user                          # Generate model, route, query, migration
  gecho migrate                                # Apply database migrations
  swag init                                    # Generate Swagger docs (before go mod tidy)
  go mod tidy                                  # Install dependencies
  go run main.go                               # Start the server

Optional:
  go get -u ./...                              # Upgrade all packages

Swagger UI:
  http://localhost:3000/swagger/index.html`)

}

func ensureGoModExists() {
	const path = "go.mod"
	if _, err := os.Stat(path); err == nil {
		return
	}

	fmt.Print("No go.mod found. Enter module name: ")
	var name string
	_, err := fmt.Scanln(&name)
	if err != nil || strings.TrimSpace(name) == "" {
		fmt.Println("Invalid module name.")
		os.Exit(1)
	}

	cmd := exec.Command("go", "mod", "init", strings.TrimSpace(name))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to create go.mod: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("go.mod created.")
}

func writeDotEnvIfMissing() {
	const path = ".env"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		content := `# Gecho default database config
DATABASE_HOST=localhost
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=gecho_dev
DATABASE_PORT=5432
`
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			fmt.Printf("Failed to write .env: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(".env file created with default values.")
	}
}

func getModuleName(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read %s: %v\n", path, err)
		os.Exit(1)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}

	fmt.Println("No module declaration found in go.mod")
	os.Exit(1)
	return ""
}

func writeMainGoIfMissing(module string) {
	const path = "main.go"
	if _, err := os.Stat(path); err == nil {
		return
	}

	content := fmt.Sprintf(`// @title Gecho API
// @version 1.0
// @description Auto-generated Echo main.go
// @host localhost
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"

	"%s/app/routes"
	"%s/app/services/database"

	_ "%s/docs"
)

var debug = false

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}
	if os.Getenv("DEBUG") == "true" {
		debug = true
	}
}

// RequestTimeLogger middleware
func RequestTimeLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		duration := time.Since(start)
		log.Printf("Request %%s %%s took %%v", c.Request().Method, c.Request().URL.Path, duration)
		return err
	}
}

func main() {
	e := echo.New()

	e.Use(RequestTimeLogger)

	database.InitDB()

	// Apply pending db/migrations on start:
	// if err := database.AutoMigrate(); err != nil {
	// 	log.Fatalf("Auto-migrate failed: %%v", err)
	// }

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	e.GET("/", routes.HelloWorld)

	// Example route:
	// e.GET("/users", routes.GetUsers)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Server is healthy"})
	})

	fmt.Println("Server running on port 3000")
	e.Logger.Fatal(e.Start(":3000"))
}
`, module, module, module)

	err := os.WriteFile(path, formatGo(content), 0644)
	if err != nil {
		fmt.Printf("Failed to write main.go: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("main.go created with Echo setup.")
}

func writeDatabaseGoIfMissing() {
	const path = "app/services/database/database.go"
	if _, err := os.Stat(path); err == nil {
		return
	}

	content := `package database

import (
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB initializes a PostgreSQL connection using env vars.
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

// GetDB returns the active DB instance
func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database not initialized. Call InitDB first.")
	}
	return db
}

// GetPGVersion returns the PostgreSQL version
func GetPGVersion() string {
	var version string
	if err := GetDB().Raw("SELECT version()").Scan(&version).Error; err != nil {
		log.Println("Failed to get PostgreSQL version:", err)
		return "Unknown"
	}
	return version
}
`

	if err := os.WriteFile(path, formatGo(content), 0644); err != nil {
		fmt.Printf("Failed to write %s: %v\n", path, err)
		os.Exit(1)
	}

	fmt.Println("app/services/database/database.go created.")
}

const MigrationsGoPath = "app/services/database/migrations.go"

// WriteMigrationsGoIfMissing writes the AutoMigrate runner and reports whether it was created.
func WriteMigrationsGoIfMissing() bool {
	path := MigrationsGoPath
	if _, err := os.Stat(path); err == nil {
		return false
	}
	if err := os.MkdirAll("app/services/database", 0755); err != nil {
		fmt.Printf("Failed to create app/services/database: %v\n", err)
		os.Exit(1)
	}

	content := `package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// migration is a row in the migrations table, shared with "gecho migrate".
type migration struct {
	ID        uint   ` + "`gorm:\"primaryKey\"`" + `
	Name      string ` + "`gorm:\"unique\"`" + `
	AppliedAt time.Time
}

// AutoMigrate applies all pending .up.sql files from db/migrations.
func AutoMigrate() error {
	db := GetDB()

	if err := db.AutoMigrate(&migration{}); err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	var applied []migration
	if err := db.Find(&applied).Error; err != nil {
		return fmt.Errorf("read applied migrations: %w", err)
	}
	done := make(map[string]bool)
	for _, m := range applied {
		done[m.Name] = true
	}

	files, err := filepath.Glob(filepath.Join("db/migrations", "*.up.sql"))
	if err != nil {
		return fmt.Errorf("list migration files: %w", err)
	}

	for _, file := range files {
		name := filepath.Base(file)
		if done[name] {
			continue
		}

		sql, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read %s: %w", file, err)
		}

		log.Printf("Applying migration: %s", name)
		if err := db.Exec(string(sql)).Error; err != nil {
			return fmt.Errorf("migration %s: %w", name, err)
		}
		if err := db.Create(&migration{Name: name, AppliedAt: time.Now()}).Error; err != nil {
			return fmt.Errorf("migration %s ran but could not be recorded: %w", name, err)
		}
	}

	return nil
}
`

	if err := os.WriteFile(path, formatGo(content), 0644); err != nil {
		fmt.Printf("Failed to write %s: %v\n", path, err)
		os.Exit(1)
	}
	return true
}

func writeHelloWorldRouteIfMissing() {
	const path = "app/routes/helloWorld.go"
	if _, err := os.Stat(path); err == nil {
		return
	}

	content := `package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HelloWorld returns a simple status message
func HelloWorld(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello from Gecho"})
}
`
	err := os.WriteFile(path, formatGo(content), 0644)
	if err != nil {
		fmt.Printf("Failed to write %s: %v\n", path, err)
		os.Exit(1)
	}

	fmt.Println("app/routes/helloWorld.go created.")
}

// formatGo gofmt-formats generated Go source, falling back to the raw text if it doesn't parse.
func formatGo(content string) []byte {
	if out, err := format.Source([]byte(content)); err == nil {
		return out
	}
	return []byte(content)
}
