package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
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
		h.store.Flash(c, "error", "All fields are required")
		return c.Redirect("/register", fiber.StatusSeeOther)
	}

	// Register user
	user, err := h.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			h.store.Flash(c, "error", "Email already registered")
			return c.Redirect("/register", fiber.StatusSeeOther)
		}
		h.store.Flash(c, "error", "Failed to register user. Please try again.")
		return c.Redirect("/register", fiber.StatusSeeOther)
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

	slog.Info("session created", "handler", "Auth.Register", "user_id", user.ID, "redirect", "/app")

	// Inertia.js will automatically follow this redirect
	return c.Redirect("/app", fiber.StatusSeeOther)
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
			// Set flash error cookie and redirect back to login
			h.store.Flash(c, "error", "Invalid email or password")
			return c.Redirect("/login", fiber.StatusSeeOther)
		}
		// Set flash error cookie and redirect back to login
		h.store.Flash(c, "error", "Failed to login. Please try again.")
		return c.Redirect("/login", fiber.StatusSeeOther)
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

	slog.Info("session created", "handler", "Auth.Login", "user_id", user.ID, "redirect", "/app")

	// Inertia.js will automatically follow this redirect
	return c.Redirect("/app", fiber.StatusSeeOther)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sess, _ := h.store.Get(c)
	sess.Destroy()

	slog.Info("user logged out", "handler", "Auth.Logout", "redirect", "/login")

	// Inertia.js will automatically follow this redirect
	return c.Redirect("/login", fiber.StatusSeeOther)
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
		slog.Warn("oauth state mismatch", "got", state, "expected", storedState)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid OAuth state",
		})
	}

	// Clear the state cookie
	c.ClearCookie("oauth_state")

	// Process the token
	user, err := h.authService.ProcessGoogleToken(c.Context(), code)
	if err != nil {
		slog.Error("google token error", "error", err)
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

	slog.Info("session created", "handler", "Auth.GoogleCallback", "user_id", user.ID, "redirect", "/app")

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
		slog.Warn("invalid user ID", "user_id_param", userIDParam)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	slog.Info("fetching user avatar", "handler", "GetAvatar", "user_id", userID)

	// Get user from database
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		slog.Warn("avatar user not found", "handler", "GetAvatar", "error", err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	slog.Info("user avatar URL", "handler", "GetAvatar", "avatar_url", user.Avatar)

	// Check if user has avatar
	if user.Avatar == "" {
		slog.Info("no avatar for user", "handler", "GetAvatar", "user_id", userID)
		return c.Status(404).JSON(fiber.Map{"error": "No avatar"})
	}

	// Check if avatar is local file or external URL
	if strings.HasPrefix(user.Avatar, "/storage/") {
		// Local file - serve directly
		localPath := "." + user.Avatar
		slog.Info("serving local avatar file", "handler", "GetAvatar", "path", localPath)
		
		return c.SendFile(localPath)
	}

	// External URL - fetch and proxy
	slog.Info("fetching avatar from external URL", "handler", "GetAvatar", "url", user.Avatar)
	resp, err := http.Get(user.Avatar)
	if err != nil {
		slog.Error("failed to fetch avatar", "handler", "GetAvatar", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch avatar"})
	}
	defer resp.Body.Close()

	slog.Info("avatar response", "handler", "GetAvatar", "status", resp.Status, "content_type", resp.Header.Get("Content-Type"))

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
		slog.Error("failed to read avatar body", "handler", "GetAvatar", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read avatar"})
	}

	slog.Info("sending avatar", "handler", "GetAvatar", "bytes", len(body))
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
