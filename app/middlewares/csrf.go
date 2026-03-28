package middlewares

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maulanashalihin/laju-go/app/session"
)

// CSRFConfig holds CSRF middleware configuration
type CSRFConfig struct {
	Secret       string        // Secret key for signing tokens
	CookieName   string        // Name of the CSRF token cookie
	HeaderName   string        // Name of the CSRF token header
	TokenLength  int           // Length of the random token
	Expiry       time.Duration // Token expiry duration
	Secure       bool          // Secure cookie flag
	SameSite     string        // SameSite cookie attribute
	SkipPaths    []string      // Paths to skip CSRF check
	SkipMethods  []string      // HTTP methods to skip CSRF check
}

// CSRFMiddleware implements CSRF protection
type CSRFMiddleware struct {
	config CSRFConfig
	store  *session.Store
}

// DefaultCSRFConfig returns a default CSRF configuration
func DefaultCSRFConfig(secret string) CSRFConfig {
	return CSRFConfig{
		Secret:      secret,
		CookieName:  "csrf_token",
		HeaderName:  "X-CSRF-Token",
		TokenLength: 32,
		Expiry:      24 * time.Hour,
		Secure:      false, // Set to true in production with HTTPS
		SameSite:    "Lax",
		SkipMethods: []string{fiber.MethodGet, fiber.MethodHead, fiber.MethodOptions},
	}
}

// NewCSRFMiddleware creates a new CSRF middleware
func NewCSRFMiddleware(store *session.Store, config CSRFConfig) *CSRFMiddleware {
	return &CSRFMiddleware{
		config: config,
		store:  store,
	}
}

// Protect returns the CSRF protection middleware
func (csrf *CSRFMiddleware) Protect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip if path is in skip list
		for _, path := range csrf.config.SkipPaths {
			if strings.HasPrefix(c.Path(), path) {
				return c.Next()
			}
		}

		// Skip if method is in skip list
		for _, method := range csrf.config.SkipMethods {
			if c.Method() == method {
				return csrf.setToken(c)
			}
		}

		// Validate token for state-changing methods
		if err := csrf.validateToken(c); err != nil {
			return err
		}

		return c.Next()
	}
}

// setToken generates and sets a CSRF token
func (csrf *CSRFMiddleware) setToken(c *fiber.Ctx) error {
	// Try to get existing token from session
	sess, _ := csrf.store.Get(c)
	token := sess.Get("csrf_token")

	if token == nil || csrf.isTokenExpired(sess) {
		// Generate new token
		var err error
		token, err = csrf.generateToken()
		if err != nil {
			return err
		}
		sess.Set("csrf_token", token)
		sess.Set("csrf_expiry", time.Now().Add(csrf.config.Expiry).Unix())
		sess.Save()
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     csrf.config.CookieName,
		Value:    token.(string),
		Path:     "/",
		MaxAge:   int(csrf.config.Expiry.Seconds()),
		Secure:   csrf.config.Secure,
		HTTPOnly: false, // Must be false to allow JavaScript access
		SameSite: csrf.config.SameSite,
	})

	// Pass token to response header for easy access
	c.Set("X-CSRF-Token", token.(string))

	return c.Next()
}

// validateToken validates the CSRF token from the request
func (csrf *CSRFMiddleware) validateToken(c *fiber.Ctx) error {
	// Get token from header, form, or query
	token := c.Get(csrf.config.HeaderName)
	if token == "" {
		token = c.FormValue(csrf.config.CookieName)
	}
	if token == "" {
		token = c.Query(csrf.config.CookieName)
	}

	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "CSRF token missing",
		})
	}

	// Get token from session
	sess, _ := csrf.store.Get(c)
	sessionToken := sess.Get("csrf_token")

	if sessionToken == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "CSRF token invalid",
		})
	}

	// Compare tokens (constant time comparison)
	if !csrf.constantTimeCompare(token, sessionToken.(string)) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "CSRF token invalid",
		})
	}

	// Check expiry
	if csrf.isTokenExpired(sess) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "CSRF token expired",
		})
	}

	return nil
}

// isTokenExpired checks if the token has expired
func (csrf *CSRFMiddleware) isTokenExpired(sess *session.Session) bool {
	expiry := sess.Get("csrf_expiry")
	if expiry == nil {
		return true
	}
	return time.Now().Unix() > expiry.(int64)
}

// generateToken generates a random CSRF token
func (csrf *CSRFMiddleware) generateToken() (string, error) {
	bytes := make([]byte, csrf.config.TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// constantTimeCompare compares two strings in constant time
func (csrf *CSRFMiddleware) constantTimeCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	result := byte(0)
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

// GetToken retrieves the current CSRF token from a request context
func GetToken(c *fiber.Ctx) string {
	return c.Cookies("csrf_token")
}
