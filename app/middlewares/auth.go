package middlewares

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/maulanashalihin/laju-go/app/session"
)

// AuthRequired is a middleware that checks if the user is authenticated
func AuthRequired(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("[AuthRequired] Checking auth for path: %s\n", c.Path())
		
		sess, err := store.Get(c)
		if err != nil {
			log.Printf("[AuthRequired] Get session error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get session",
			})
		}

		userID := sess.Get("user_id")
		log.Printf("[AuthRequired] userID=%v\n", userID)
		
		if userID == nil {
			log.Printf("[AuthRequired] Not authenticated, redirecting to /login\n")
			// For Inertia requests, return redirect in JSON format
			if c.Get("X-Inertia") == "true" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"component": "Login",
					"props": fiber.Map{
						"error": "Please login to continue",
					},
				})
			}
			return c.Redirect("/login")
		}

		// Store user info in locals for handlers to use
		c.Locals("user_id", userID)
		c.Locals("email", sess.Get("email"))
		c.Locals("role", sess.Get("role"))
		log.Printf("[AuthRequired] Auth successful, user_id=%v\n", userID)

		return c.Next()
	}
}

// AdminRequired is a middleware that checks if the user is an admin
func AdminRequired(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get session",
			})
		}

		userID := sess.Get("user_id")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Not authenticated",
			})
		}

		role := sess.Get("role")
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		c.Locals("user_id", userID)
		c.Locals("email", sess.Get("email"))
		c.Locals("role", role)

		return c.Next()
	}
}

// Guest is a middleware that redirects authenticated users away from login/register pages
func Guest(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Next()
		}

		userID := sess.Get("user_id")
		log.Printf("[Guest] userID=%v\n", userID)
		
		if userID != nil {
			log.Printf("[Guest] Already authenticated, redirecting to /app\n")
			return c.Redirect("/app")
		}

		return c.Next()
	}
}

// CORS creates a CORS middleware handler
func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Inertia, X-Inertia-Version, X-Requested-With")
		c.Set("Access-Control-Allow-Credentials", "true")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}

// Logger creates a simple logger middleware
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Simple logging - in production, use a proper logger
		c.Next()
		return nil
	}
}
