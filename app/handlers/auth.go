package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maulanashalihin/laju-go/app/models"
	"github.com/maulanashalihin/laju-go/app/services"
	"github.com/maulanashalihin/laju-go/app/session"
)

type AuthHandler struct {
	authService    *services.AuthService
	userService    *services.UserService
	store          *session.Store
	inertiaService *services.InertiaService
}

func NewAuthHandler(authService *services.AuthService, userService *services.UserService, store *session.Store, inertiaService *services.InertiaService) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		userService:    userService,
		store:          store,
		inertiaService: inertiaService,
	}
}

// ShowLoginForm displays the login page
func (h *AuthHandler) ShowLoginForm(c *fiber.Ctx) error {
	return h.inertiaService.Render(c, "auth/Login", fiber.Map{
		"Title": "Login",
	})
}

// ShowRegisterForm displays the register page
func (h *AuthHandler) ShowRegisterForm(c *fiber.Ctx) error {
	return h.inertiaService.Render(c, "auth/Register", fiber.Map{
		"Title": "Register",
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required",
		})
	}

	// Register user
	user, err := h.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already registered",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	// Create session
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session",
		})
	}
	sess.Set("user_id", user.ID)
	sess.Set("email", user.Email)
	sess.Set("role", string(user.Role))

	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save session",
		})
	}

	log.Printf("[Auth.Register] Session created for user %d, redirecting to /app\n", user.ID)

	// Inertia.js will automatically follow this redirect
	return c.Redirect("/app")
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// Authenticate user
	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to login",
		})
	}

	// Create session
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session",
		})
	}
	sess.Set("user_id", user.ID)
	sess.Set("email", user.Email)
	sess.Set("role", string(user.Role))

	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save session",
		})
	}

	log.Printf("[Auth.Login] Session created for user %d, redirecting to /app\n", user.ID)

	// Inertia.js will automatically follow this redirect
	return c.Redirect("/app")
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sess, _ := h.store.Get(c)
	sess.Destroy()

	log.Printf("[Auth.Logout] User logged out, redirecting to /login\n")

	// Inertia.js will automatically follow this redirect
	return c.Redirect("/login")
}

// GoogleLogin initiates Google OAuth login
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	state := generateState()
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   300, // 5 minutes
		HTTPOnly: true,
		SameSite: "Lax",
	})

	url := h.authService.GetOAuthURL(state)
	return c.Redirect(url)
}

// GoogleCallback handles Google OAuth callback
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")

	// Validate state
	storedState := c.Cookies("oauth_state")
	if state != storedState {
		log.Printf("State mismatch: got=%s, expected=%s\n", state, storedState)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid OAuth state",
		})
	}

	// Clear the state cookie
	c.ClearCookie("oauth_state")

	// Process the token
	user, err := h.authService.ProcessGoogleToken(c.Context(), code)
	if err != nil {
		log.Printf("Google token error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to authenticate with Google: " + err.Error(),
		})
	}

	// Create session
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session",
		})
	}
	sess.Set("user_id", user.ID)
	sess.Set("email", user.Email)
	sess.Set("role", string(user.Role))

	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save session",
		})
	}

	log.Printf("[Auth.GoogleCallback] Session created for user %d, redirecting to /app\n", user.ID)

	// Inertia.js will automatically follow this redirect
	return c.Redirect("/app")
}

// Me returns the current authenticated user
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	sess, _ := h.store.Get(c)
	userID := sess.Get("user_id")

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	user, err := h.authService.GetUserByID(userID.(int64))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	return c.JSON(fiber.Map{
		"user": user.ToResponse(),
	})
}

// GetAvatar proxies user avatar from external URL (e.g., Google) or serves local file
func (h *AuthHandler) GetAvatar(c *fiber.Ctx) error {
	userIDParam := c.Params("id")
	if userIDParam == "" {
		return c.Status(400).JSON(fiber.Map{"error": "User ID required"})
	}

	// Convert userID to int64
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		log.Printf("Invalid user ID: %s\n", userIDParam)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	log.Printf("[GetAvatar] Fetching user %d\n", userID)

	// Get user from database
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		log.Printf("[GetAvatar] User not found: %v\n", err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	log.Printf("[GetAvatar] User avatar URL: %s\n", user.Avatar)

	// Check if user has avatar
	if user.Avatar == "" {
		log.Printf("[GetAvatar] No avatar for user %d\n", userID)
		return c.Status(404).JSON(fiber.Map{"error": "No avatar"})
	}

	// Check if avatar is local file or external URL
	if strings.HasPrefix(user.Avatar, "/storage/") {
		// Local file - serve directly
		localPath := "." + user.Avatar
		log.Printf("[GetAvatar] Serving local file: %s\n", localPath)
		
		return c.SendFile(localPath)
	}

	// External URL - fetch and proxy
	log.Printf("[GetAvatar] Fetching from external URL: %s\n", user.Avatar)
	resp, err := http.Get(user.Avatar)
	if err != nil {
		log.Printf("[GetAvatar] Failed to fetch avatar: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch avatar"})
	}
	defer resp.Body.Close()

	log.Printf("[GetAvatar] Response status: %s, Content-Type: %s\n", resp.Status, resp.Header.Get("Content-Type"))

	// Set headers
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}
	c.Set("Content-Type", contentType)
	c.Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours

	// Read and send response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[GetAvatar] Failed to read body: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read avatar"})
	}

	log.Printf("[GetAvatar] Sending %d bytes\n", len(body))
	return c.Send(body)
}

// generateState generates a random state string for OAuth
func generateState() string {
	// Generate random bytes
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based
		return fmt.Sprintf("state_%d", time.Now().UnixNano())
	}
	// Convert to hex string
	return hex.EncodeToString(b)
}
