package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// NoCache sets HTTP headers that prevent CDN/proxy caching of authenticated
// responses. Without this, a reverse proxy or Cloudflare can cache one user's
// Inertia page (HTML + JSON props containing user_id, email, role) and serve
// it to a different user — causing session mixing where user A sees user B's
// data.
//
// Applied to /app/* and /admin/* route groups.
func NoCache() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
		c.Set("Pragma", "no-cache")
		c.Set("Expires", "0")
		c.Set("Vary", "Cookie")
		return c.Next()
	}
}

// NoCacheByDefault is the safety net behind NoCache: it guarantees that no
// dynamic response is edge-cacheable unless it explicitly opted in via its
// own Cache-Control header (static assets do this through MaxAge).
//
// Background — the laju-go cache leak incident: every Inertia HTML page
// embeds per-session bytes (csrf_token <meta>, shared flash props), but
// only /app/* and /admin/* carried a Cache-Control header. The remaining
// pages (/about, /login, /register, password-reset, OAuth, …) sent none —
// so a single Cloudflare "cache all GET" rule stored user A's HTML and
// served it to user B. Tagging route groups one by one cannot close this
// class of bug; a future route without NoCache reopens it.
//
// This middleware runs after the handler (post-next) and fills the gap:
// any text/html or JSON response that still has no Cache-Control gets
// private, no-store. Responses that already set one (static files with
// MaxAge, explicit NoCache groups) are left untouched.
func NoCacheByDefault() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			return err
		}
		// Explicit opt-in wins (static assets via MaxAge, NoCache groups).
		if c.GetRespHeader(fiber.HeaderCacheControl) != "" {
			return nil
		}
		contentType := c.GetRespHeader(fiber.HeaderContentType)
		if strings.HasPrefix(contentType, "text/html") ||
			strings.HasPrefix(contentType, "application/json") {
			c.Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
			c.Set("Pragma", "no-cache")
			c.Set("Expires", "0")
			c.Set("Vary", "Cookie")
		}
		return nil
	}
}
