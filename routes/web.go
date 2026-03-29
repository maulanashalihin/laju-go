package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/maulanashalihin/laju-go/app/handlers"
	"github.com/maulanashalihin/laju-go/app/middlewares"
	"github.com/maulanashalihin/laju-go/app/services"
	"github.com/maulanashalihin/laju-go/app/session"
)

type Handlers struct {
	Public        *handlers.PublicHandler
	Auth          *handlers.AuthHandler
	App           *handlers.AppHandler
	Upload        *handlers.UploadHandler
	PasswordReset *handlers.PasswordResetHandler
}

func SetupRoutes(app *fiber.App, handlers Handlers, store *session.Store, mailerService *services.MailerService, csrfMiddleware *middlewares.CSRFMiddleware) {
	// Setup static file serving
	setupStaticRoutes(app)

	// Setup public routes
	setupPublicRoutes(app, handlers.Public)

	// Setup auth routes
	setupAuthRoutes(app, handlers.Auth, handlers.PasswordReset, store, mailerService)

	// Setup app routes (protected)
	setupAppRoutes(app, handlers.App, handlers.Upload, store, csrfMiddleware)
}

func setupStaticRoutes(app *fiber.App) {
	app.Static("/dist", "./dist")      // Serve built frontend assets
	app.Static("/public", "./public")  // Serve public assets at /public path
	app.Static("/storage", "./storage") // Serve uploaded files (avatars, etc.)
}

func setupPublicRoutes(app *fiber.App, handler *handlers.PublicHandler) {
	app.Get("/", handler.Index)
	app.Get("/about", handler.About)
}

func setupAuthRoutes(app *fiber.App, authHandler *handlers.AuthHandler, passwordResetHandler *handlers.PasswordResetHandler, store *session.Store, mailerService *services.MailerService) {
	// Login routes (with Guest middleware)
	app.Get("/login", middlewares.Guest(store), authHandler.ShowLoginForm)
	app.Post("/login", middlewares.Guest(store), authHandler.Login, middlewares.AuthRateLimit.Limit())

	// Register routes (with Guest middleware)
	app.Get("/register", middlewares.Guest(store), authHandler.ShowRegisterForm)
	app.Post("/register", middlewares.Guest(store), authHandler.Register, middlewares.AuthRateLimit.Limit())

	// OAuth routes
	app.Get("/auth/google", authHandler.GoogleLogin)
	app.Get("/auth/google/callback", authHandler.GoogleCallback)

	// Logout (requires auth)
	app.Post("/logout", middlewares.AuthRequired(store), authHandler.Logout)

	// API: Get current user
	app.Get("/api/me", middlewares.AuthRequired(store), authHandler.Me)

	// API: Get user avatar (proxied from external URL)
	app.Get("/api/avatar/:id", authHandler.GetAvatar)

	// Password reset routes
	app.Get("/forgot-password", passwordResetHandler.ShowForgotPasswordForm)
	app.Post("/forgot-password", passwordResetHandler.SendResetLink, middlewares.PasswordResetRateLimit.Limit())
	app.Get("/reset-password/:token", passwordResetHandler.ShowResetPasswordForm)
	app.Post("/reset-password/:token", passwordResetHandler.ResetPassword)
}

func setupAppRoutes(app *fiber.App, appHandler *handlers.AppHandler, uploadHandler *handlers.UploadHandler, store *session.Store, csrfMiddleware *middlewares.CSRFMiddleware) {
	// Protected app routes with CSRF protection
	protected := app.Group("/app", middlewares.AuthRequired(store))
	protected.Use(csrfMiddleware.Protect())

	// Dashboard
	protected.Get("/", appHandler.Dashboard)

	// Profile
	protected.Get("/profile", appHandler.Profile)
	protected.Put("/profile", appHandler.UpdateProfile)
	protected.Put("/profile/password", appHandler.UpdatePassword)

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

// SetupCSRFMiddleware sets up the CSRF middleware
func SetupCSRFMiddleware(store *session.Store, secret string) *middlewares.CSRFMiddleware {
	config := middlewares.DefaultCSRFConfig(secret)
	config.Secure = false // Set to true in production with HTTPS
	config.SameSite = "Lax"
	return middlewares.NewCSRFMiddleware(store, config)
}

// SetupMailerService sets up the mailer service
func SetupMailerService(smtpHost string, smtpPort int, smtpUser, smtpPass, fromEmail, fromName string) *services.MailerService {
	return services.NewMailerService(smtpHost, smtpPort, smtpUser, smtpPass, fromEmail, fromName)
}

// SetupPasswordResetHandler sets up the password reset handler
func SetupPasswordResetHandler(
	mailerService *services.MailerService,
	userService *services.UserService,
	store *session.Store,
	inertiaService *services.InertiaService,
	appURL string,
) *handlers.PasswordResetHandler {
	return handlers.NewPasswordResetHandler(
		mailerService,
		userService,
		store,
		inertiaService,
		appURL,
	)
}

// GetAppURL returns the application URL based on environment
func GetAppURL(appPort string, appEnv string) string {
	if appEnv == "production" {
		return "https://yourdomain.com"
	}
	return fmt.Sprintf("http://localhost:%s", appPort)
}
