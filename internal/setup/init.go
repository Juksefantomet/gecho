package setup

import (
	"fmt"
	"os"
	"strings"
)

var requiredDirs = []string{
	"app/models",
	"app/routes",
	"app/services/database",
	"db/migrations",
}

func RunInit() {
	for _, dir := range requiredDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("✗ Failed to create %s: %v\n", dir, err)
			os.Exit(1)
		}
	}
	fmt.Println("✓ Gecho-compatible project structure initialized.")
	fmt.Println("→ Ready for scaffolding.")
	writeDotEnvIfMissing()
	writeMainGoIfMissing(getModuleName("go.mod"))
	writeDatabaseGoIfMissing()
	writeHelloWorldRouteIfMissing(getModuleName("go.mod"))
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
			fmt.Printf("✗ Failed to write .env: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ .env file created with default values.")
	}
}

func getModuleName(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("✗ Failed to read %s: %v\n", path, err)
		os.Exit(1)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}

	fmt.Println("✗ No module declaration found in go.mod")
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
		log.Printf("Request %s %s took %v", c.Request().Method, c.Request().URL.Path, duration)
		return err
	}
}

func main() {
	e := echo.New()

	e.Use(RequestTimeLogger)

	database.InitDB()

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
`, module, module, module, module)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Printf("✗ Failed to write main.go: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ main.go created with Echo setup.")
}

func writeDatabaseGoIfMissing() {
	const path = "app/services/database/database.go"
	if _, err := os.Stat(path); err == nil {
		return
	}

	content := `package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB initializes a PostgreSQL connection using env vars
func InitDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PORT"),
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %%v", err)
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

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Printf("✗ Failed to write %s: %v\n", path, err)
		os.Exit(1)
	}

	fmt.Println("✓ app/services/database/database.go created.")
}

func writeHelloWorldRouteIfMissing(module string) {
	const path = "app/routes/helloWorld.go"
	if _, err := os.Stat(path); err == nil {
		return
	}

	content := fmt.Sprintf(`package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HelloWorld returns a simple status message
func HelloWorld(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello from Gecho"})
}
`)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Printf("✗ Failed to write %s: %v\n", path, err)
		os.Exit(1)
	}

	fmt.Println("✓ app/routes/helloWorld.go created.")
}
