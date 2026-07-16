package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
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
	store          *session.Store
	inertiaService *services.InertiaService
}

func NewAuthHandler(authService *services.AuthService, store *session.Store, inertiaService *services.InertiaService) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		store:          store,
		inertiaService: inertiaService,
	}
}

func (h *AuthHandler) ShowLoginForm(c *fiber.Ctx) error {
	return h.inertiaService.Render(c, "auth/Login", fiber.Map{
		"Title": "Login",
	})
}

func (h *AuthHandler) ShowRegisterForm(c *fiber.Ctx) error {
	return h.inertiaService.Render(c, "auth/Register", fiber.Map{
		"Title": "Register",
	})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		h.store.Flash(c, "error", "All fields are required")
		return h.inertiaService.Redirect(c, "/register")
	}

	user, err := h.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			h.store.Flash(c, "error", "Email already registered")
			return h.inertiaService.Redirect(c, "/register")
		}
		h.store.Flash(c, "error", "Failed to register user. Please try again.")
		return h.inertiaService.Redirect(c, "/register")
	}

	if err := h.store.CreateAuthenticatedSession(c, user.ID, user.Name, user.Email, user.Avatar, string(user.Role), user.EmailVerified); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session",
		})
	}

	// Regenerate session ID to prevent session fixation
	if sess, err := h.store.Get(c); err == nil {
		sess.Regenerate()
	}

	slog.Info("session created", "handler", "Auth.Register", "user_id", user.ID, "redirect", "/app")
	return h.inertiaService.Redirect(c, "/app")
}

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

	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			h.store.Flash(c, "error", "Invalid email or password")
			return h.inertiaService.Redirect(c, "/login")
		}
		h.store.Flash(c, "error", "Failed to login. Please try again.")
		return h.inertiaService.Redirect(c, "/login")
	}

	if err := h.store.CreateAuthenticatedSession(c, user.ID, user.Name, user.Email, user.Avatar, string(user.Role), user.EmailVerified); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session",
		})
	}

	// Regenerate session ID to prevent session fixation
	if sess, err := h.store.Get(c); err == nil {
		sess.Regenerate()
	}

	slog.Info("session created", "handler", "Auth.Login", "user_id", user.ID, "redirect", "/app")
	return h.inertiaService.Redirect(c, "/app")
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sess, _ := h.store.Get(c)
	sess.Destroy()

	slog.Info("user logged out", "handler", "Auth.Logout", "redirect", "/login")

	return h.inertiaService.Redirect(c, "/login")
}

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
	// Use Location() so Inertia triggers a full window.location navigation
	// to Google's OAuth page (not an XHR follow).
	return h.inertiaService.Location(c, url)
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")

	storedState := c.Cookies("oauth_state")
	if state != storedState {
		slog.Warn("oauth state mismatch", "got", state, "expected", storedState)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid OAuth state",
		})
	}

	c.ClearCookie("oauth_state")

	user, err := h.authService.ProcessGoogleToken(c.Context(), code)
	if err != nil {
		slog.Error("google token error", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to authenticate with Google: " + err.Error(),
		})
	}

	// Create session
	if err := h.store.CreateAuthenticatedSession(c, user.ID, user.Name, user.Email, user.Avatar, string(user.Role), user.EmailVerified); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session",
		})
	}

	// Regenerate session ID to prevent session fixation
	if sess, err := h.store.Get(c); err == nil {
		sess.Regenerate()
	}

	slog.Info("session created", "handler", "Auth.GoogleCallback", "user_id", user.ID, "redirect", "/app")

	return h.inertiaService.Redirect(c, "/app")
}

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

	slog.Debug("fetching user avatar", "handler", "GetAvatar", "user_id", userID)

	// Get user from database
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		slog.Warn("avatar user not found", "handler", "GetAvatar", "error", err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	slog.Debug("user avatar URL", "handler", "GetAvatar", "avatar_url", user.Avatar)

	// Check if user has avatar
	if user.Avatar == "" {
		slog.Debug("no avatar for user", "handler", "GetAvatar", "user_id", userID)
		return c.Status(404).JSON(fiber.Map{"error": "No avatar"})
	}

	// Check if avatar is local file or external URL
	if strings.HasPrefix(user.Avatar, "/storage/") {
		// Local file - serve directly
		localPath := "." + user.Avatar
		slog.Debug("serving local avatar file", "handler", "GetAvatar", "path", localPath)

		return c.SendFile(localPath)
	}

	// External URL - fetch and proxy
	slog.Debug("fetching avatar from external URL", "handler", "GetAvatar", "url", user.Avatar)
	resp, err := http.Get(user.Avatar)
	if err != nil {
		slog.Error("failed to fetch avatar", "handler", "GetAvatar", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch avatar"})
	}
	defer resp.Body.Close()

	slog.Debug("avatar response", "handler", "GetAvatar", "status", resp.Status, "content_type", resp.Header.Get("Content-Type"))

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

	slog.Debug("sending avatar", "handler", "GetAvatar", "bytes", len(body))
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
