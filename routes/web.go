package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/velostack/velostack-go/app/handlers"
	"github.com/velostack/velostack-go/app/middlewares"
	"github.com/velostack/velostack-go/app/session"
)

type Handlers struct {
	Public *handlers.PublicHandler
	Auth   *handlers.AuthHandler
	App    *handlers.AppHandler
	Upload *handlers.UploadHandler
}

func SetupRoutes(app *fiber.App, handlers Handlers, store *session.Store) {
	// Setup public routes
	setupPublicRoutes(app, handlers.Public)

	// Setup auth routes
	setupAuthRoutes(app, handlers.Auth, store)

	// Setup app routes (protected)
	setupAppRoutes(app, handlers.App, handlers.Upload, store)
}

func setupPublicRoutes(app *fiber.App, handler *handlers.PublicHandler) {
	app.Get("/", handler.Index)
	app.Get("/about", handler.About)
}

func setupAuthRoutes(app *fiber.App, handler *handlers.AuthHandler, store *session.Store) {
	// Guest routes (redirect if already logged in)
	guest := app.Group("/login", middlewares.Guest(store))
	guest.Get("/", handler.ShowLoginForm)
	guest.Post("/login", handler.Login)
	guest.Post("/register", handler.Register)

	// OAuth routes
	app.Get("/auth/google", handler.GoogleLogin)
	app.Get("/auth/google/callback", handler.GoogleCallback)

	// Logout (requires auth)
	app.Post("/logout", middlewares.AuthRequired(store), handler.Logout)

	// API: Get current user
	app.Get("/api/me", middlewares.AuthRequired(store), handler.Me)
}

func setupAppRoutes(app *fiber.App, appHandler *handlers.AppHandler, uploadHandler *handlers.UploadHandler, store *session.Store) {
	// Protected app routes
	protected := app.Group("/app", middlewares.AuthRequired(store))

	// Dashboard
	protected.Get("/", appHandler.Dashboard)

	// Profile
	protected.Get("/profile", appHandler.Profile)
	protected.Put("/profile", appHandler.UpdateProfile)

	// Upload
	protected.Post("/upload", uploadHandler.Upload)

	// Admin-only routes
	admin := app.Group("/admin", middlewares.AdminRequired(store))
	admin.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Admin dashboard",
		})
	})
}
