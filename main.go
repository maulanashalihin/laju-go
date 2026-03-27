package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/pressly/goose/v3"
	"github.com/velostack/velostack-go/app/config"
	"github.com/velostack/velostack-go/app/handlers"
	"github.com/velostack/velostack-go/app/repositories"
	"github.com/velostack/velostack-go/app/services"
	"github.com/velostack/velostack-go/app/session"
	"github.com/velostack/velostack-go/routes"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := initDatabase(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db, "./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize session store
	sessionStore := session.New(cfg.SessionSecret)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, services.AuthServiceConfig{
		SessionSecret:      cfg.SessionSecret,
		GoogleClientID:     cfg.GoogleClientID,
		GoogleClientSecret: cfg.GoogleClientSecret,
		GoogleRedirectURL:  cfg.GoogleRedirectURL,
	})
	userService := services.NewUserService(userRepo)

	// Initialize Asset service (for production builds with hashed filenames)
	assetService := services.NewAssetService("./dist/.vite/manifest.json", ".vite-port")

	// Initialize Inertia service (auto-detects Vite from .vite-port)
	inertiaService := services.NewInertiaService(assetService)

	// Initialize handlers
	routeHandlers := routes.Handlers{
		Public: handlers.NewPublicHandler(authService, userService, inertiaService, assetService),
		Auth:   handlers.NewAuthHandler(authService, userService, sessionStore, inertiaService),
		App:    handlers.NewAppHandler(userService, sessionStore, inertiaService),
		Upload: handlers.NewUploadHandler(sessionStore),
	}

	// Initialize Fiber app
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		AppName:      "VeloStack Go",
		ErrorHandler: customErrorHandler,
		Views:        engine,
	})

	// Global middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true // Allow all origins in development
		},
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Inertia, X-Inertia-Version, X-Requested-With",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Static files (production only, dev uses Vite dev server)
	app.Static("/storage", "./storage")
	app.Static("/", "./public")

	// Setup routes
	routes.SetupRoutes(app, routeHandlers, sessionStore)

	// Start server
	log.Printf("Starting server on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initDatabase initializes the SQLite database with optimized settings
func initDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Configure connection pooling
	db.SetMaxOpenConns(25)                // Maximum number of open connections
	db.SetMaxIdleConns(5)                 // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime for a connection

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	// Optimize SQLite for production (WAL mode for better concurrency)
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, err
	}

	// Balance between durability and performance
	if _, err := db.Exec("PRAGMA synchronous = NORMAL"); err != nil {
		return nil, err
	}

	// Set cache size to 64MB (negative value = KB)
	if _, err := db.Exec("PRAGMA cache_size = -64000"); err != nil {
		return nil, err
	}

	// Store temporary tables in memory for better performance
	if _, err := db.Exec("PRAGMA temp_store = MEMORY"); err != nil {
		return nil, err
	}

	// Set busy timeout to 5 seconds (wait for locks instead of failing immediately)
	if _, err := db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Log database optimization status
	logDatabaseOptimizations(db)

	return db, nil
}

// logDatabaseOptimizations logs the current SQLite optimization settings
func logDatabaseOptimizations(db *sql.DB) {
	var journalMode, synchronous string
	var cacheSize, busyTimeout int

	// Query current settings
	err := db.QueryRow("PRAGMA journal_mode").Scan(&journalMode)
	if err != nil {
		log.Printf("Warning: Could not verify journal_mode: %v", err)
	}

	err = db.QueryRow("PRAGMA synchronous").Scan(&synchronous)
	if err != nil {
		log.Printf("Warning: Could not verify synchronous: %v", err)
	}

	err = db.QueryRow("PRAGMA cache_size").Scan(&cacheSize)
	if err != nil {
		log.Printf("Warning: Could not verify cache_size: %v", err)
	}

	err = db.QueryRow("PRAGMA busy_timeout").Scan(&busyTimeout)
	if err != nil {
		log.Printf("Warning: Could not verify busy_timeout: %v", err)
	}

	log.Printf("SQLite optimizations: journal_mode=%s, synchronous=%s, cache_size=%dKB, busy_timeout=%dms",
		journalMode, synchronous, cacheSize, busyTimeout)
}

// runMigrations runs database migrations
func runMigrations(db *sql.DB, migrationsDir string) error {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
}

// customErrorHandler handles Fiber errors
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// For Inertia requests, return JSON
	if c.Get("X-Inertia") == "true" {
		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set Content-Type: application/json; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Return custom error page
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
