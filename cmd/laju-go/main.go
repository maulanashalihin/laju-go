package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pressly/goose/v3"
	"github.com/maulanashalihin/laju-go/app/cache"
	"github.com/maulanashalihin/laju-go/app/config"
	"github.com/maulanashalihin/laju-go/app/handlers"
	"github.com/maulanashalihin/laju-go/app/queries"
	"github.com/maulanashalihin/laju-go/app/services"
	"github.com/maulanashalihin/laju-go/app/session"
	"github.com/maulanashalihin/laju-go/routes"

	_ "modernc.org/sqlite"
)

var (
	Version = "dev"
	Commit  = "none"
)

func main() {
	showVersion := flag.Bool("version", false, "show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("laju-go %s (commit: %s)\n", Version, Commit)
		os.Exit(0)
	}

	// Load configuration
	cfg := config.Load()

	logLevel := slog.LevelInfo
	if cfg.AppEnv == "development" {
		logLevel = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{Level: logLevel}
	var handler slog.Handler
	if cfg.AppEnv == "development" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(handler))
	slog.Info("starting", "version", Version, "commit", Commit, "env", cfg.AppEnv)

	// Initialize database
	db, err := initDatabase(cfg.DBPath)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db, "./migrations"); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Initialize querier
	querier := queries.NewQuerier(db)

	// Initialize user profile cache
	userCache := cache.NewUserCache(cfg.UserCacheTTL)

	// Initialize session cache (in-memory, avoids DB lookup on every request)
	sessionCache := cache.NewSessionCache(cfg.SessionCacheTTL)

	// Initialize session store with database + in-memory cache
	sessionStore := session.New(querier, sessionCache)

	// Initialize services
	authService := services.NewAuthService(querier, services.AuthServiceConfig{
		SessionSecret:      cfg.SessionSecret,
		GoogleClientID:     cfg.GoogleClientID,
		GoogleClientSecret: cfg.GoogleClientSecret,
		GoogleRedirectURL:  cfg.GoogleRedirectURL,
		BcryptCost:         cfg.BcryptCost,
	})
	userService := services.NewUserService(querier, userCache)

	// Initialize Asset service (for production builds with hashed filenames)
	assetService := services.NewAssetService("./dist/.vite/manifest.json", ".vite-port", cfg.IsDevelopment())

	// Initialize Inertia service (auto-detects Vite from .vite-port)
	inertiaService := services.NewInertiaService(assetService, sessionStore)

	// Initialize handlers
	routeHandlers := routes.Handlers{
		Public: handlers.NewPublicHandler(authService, userService, inertiaService, assetService),
		Auth:   handlers.NewAuthHandler(authService, userService, sessionStore, inertiaService),
		App:    handlers.NewAppHandler(userService, sessionStore, inertiaService),
		Upload: handlers.NewUploadHandler(sessionStore, userService),
	}

	// Setup CSRF middleware
	csrfMiddleware := routes.SetupCSRFMiddleware(sessionStore, cfg.SessionSecret)

	// Setup mailer service
	mailerService := routes.SetupMailerService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
		cfg.FromEmail,
		cfg.FromName,
	)

	// Setup password reset handler
	appURL := routes.GetAppURL(cfg.AppPort, cfg.AppEnv)
	passwordResetHandler := routes.SetupPasswordResetHandler(
		mailerService,
		userService,
		sessionStore,
		inertiaService,
		appURL,
	)
	routeHandlers.PasswordReset = passwordResetHandler

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Laju",
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	// Logger: only in development (avoids string allocation per request in prod)
	if cfg.IsDevelopment() {
		app.Use(logger.New())
	}
	app.Use(recover.New())

	// Response compression (brotli > gzip, best speed for low CPU overhead)
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// CORS with explicit allowed origins (no AllowOriginsFunc in production)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.AllowedOrigins, ","),
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Inertia, X-Inertia-Version, X-Requested-With",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"version": Version,
		})
	})

	// Setup routes (includes static file serving)
	routes.SetupRoutes(app, routeHandlers, sessionStore, mailerService, csrfMiddleware)

	go func() {
		slog.Info("server listening", "port", cfg.AppPort)
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	slog.Info("shutting down", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}

// initDatabase initializes the SQLite database with optimized settings
func initDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Configure connection pooling (optimized for SQLite single-instance)
	db.SetMaxOpenConns(15)                  // Maximum number of open connections
	db.SetMaxIdleConns(10)                  // Keep more idle connections ready (avoid churn)
	db.SetConnMaxLifetime(5 * time.Minute)  // Maximum lifetime for a connection
	db.SetConnMaxIdleTime(30 * time.Second) // Recycle stale idle connections

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

	// Set cache size to 16MB (negative value = KB) - optimized for Vultr HF 1-2GB RAM
	if _, err := db.Exec("PRAGMA cache_size = -16000"); err != nil {
		return nil, err
	}

	// Enable memory-mapped I/O for NVMe performance (256MB)
	if _, err := db.Exec("PRAGMA mmap_size = 268435456"); err != nil {
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

	// Set WAL autocheckpoint to 1000 pages (default, but explicit is better)
	if _, err := db.Exec("PRAGMA wal_autocheckpoint = 1000"); err != nil {
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
	var cacheSize, busyTimeout, mmapSize, walAutocheckpoint int

	// Query current settings
	err := db.QueryRow("PRAGMA journal_mode").Scan(&journalMode)
	if err != nil {
		slog.Warn("could not verify journal_mode", "error", err)
	}

	err = db.QueryRow("PRAGMA synchronous").Scan(&synchronous)
	if err != nil {
		slog.Warn("could not verify synchronous", "error", err)
	}

	err = db.QueryRow("PRAGMA cache_size").Scan(&cacheSize)
	if err != nil {
		slog.Warn("could not verify cache_size", "error", err)
	}

	err = db.QueryRow("PRAGMA busy_timeout").Scan(&busyTimeout)
	if err != nil {
		slog.Warn("could not verify busy_timeout", "error", err)
	}

	err = db.QueryRow("PRAGMA mmap_size").Scan(&mmapSize)
	if err != nil {
		slog.Warn("could not verify mmap_size", "error", err)
	}

	err = db.QueryRow("PRAGMA wal_autocheckpoint").Scan(&walAutocheckpoint)
	if err != nil {
		slog.Warn("could not verify wal_autocheckpoint", "error", err)
	}

	slog.Info("sqlite optimizations",
		"journal_mode", journalMode,
		"synchronous", synchronous,
		"cache_size_kb", cacheSize,
		"mmap_size_kb", mmapSize,
		"wal_autocheckpoint", walAutocheckpoint,
		"busy_timeout_ms", busyTimeout,
	)
}

// runMigrations runs database migrations
func runMigrations(db *sql.DB, migrationsDir string) error {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
}

func Render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return component.Render(c.Context(), c.Response().BodyWriter())
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

