package middlewares

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Regression tests for the laju-go cache leak incident: user-specific HTML
// (Inertia pages embed per-session csrf_token + flash props) was stored by
// a Cloudflare "cache all GET" rule and served to a different user,
// because only /app/* and /admin/* sent Cache-Control. NoCacheByDefault
// closes the class: any HTML/JSON response without an explicit
// Cache-Control gets private, no-store.
func setupNoCacheApp() *fiber.App {
	app := fiber.New()
	app.Use(NoCacheByDefault())
	return app
}

func TestNoCacheByDefault_SetsNoStoreOnBareHTML(t *testing.T) {
	app := setupNoCacheApp()
	app.Get("/about", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString("<html>user-specific</html>")
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/about", nil))
	require.NoError(t, err)
	assert.Equal(t, "no-store, no-cache, must-revalidate, private", resp.Header.Get("Cache-Control"))
	assert.Equal(t, "Cookie", resp.Header.Get("Vary"))
}

func TestNoCacheByDefault_SetsNoStoreOnBareJSON(t *testing.T) {
	app := setupNoCacheApp()
	app.Get("/xhr", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"component": "About"})
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/xhr", nil))
	require.NoError(t, err)
	assert.Equal(t, "no-store, no-cache, must-revalidate, private", resp.Header.Get("Cache-Control"))
}

func TestNoCacheByDefault_PreservesExplicitCacheControl(t *testing.T) {
	app := setupNoCacheApp()
	// Static-asset style: explicit MaxAge must win (edge caching intact).
	app.Get("/dist/app-abc123.js", func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "public, max-age=31536000")
		c.Set("Content-Type", "application/javascript")
		return c.SendString("console.log(1)")
	})
	// Explicit NoCache groups must also pass through untouched.
	app.Get("/app", NoCache(), func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString("<html>dashboard</html>")
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/dist/app-abc123.js", nil))
	require.NoError(t, err)
	assert.Equal(t, "public, max-age=31536000", resp.Header.Get("Cache-Control"))

	resp, err = app.Test(httptest.NewRequest("GET", "/app", nil))
	require.NoError(t, err)
	assert.Equal(t, "no-store, no-cache, must-revalidate, private", resp.Header.Get("Cache-Control"))
}

func TestNoCacheByDefault_IgnoresNonHTMLNonJSON(t *testing.T) {
	app := setupNoCacheApp()
	app.Get("/avatar.png", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "image/png")
		return c.Send([]byte{0x89, 0x50})
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/avatar.png", nil))
	require.NoError(t, err)
	assert.Empty(t, resp.Header.Get("Cache-Control"))
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Len(t, body, 2)
}

func TestNoCacheByDefault_PropagatesHandlerErrors(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusTeapot).SendString("err")
		},
	})
	app.Use(NoCacheByDefault())
	app.Get("/boom", func(c *fiber.Ctx) error {
		return fiber.ErrTeapot
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/boom", nil))
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusTeapot, resp.StatusCode)
}
