package handlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/velostack/velostack-go/app/services"
	"github.com/velostack/velostack-go/app/session"
	"golang.org/x/crypto/bcrypt"
)

type PasswordResetHandler struct {
	mailerService  *services.MailerService
	userService    *services.UserService
	store          *session.Store
	inertiaService *services.InertiaService
	appURL         string
}

func NewPasswordResetHandler(
	mailerService *services.MailerService,
	userService *services.UserService,
	store *session.Store,
	inertiaService *services.InertiaService,
	appURL string,
) *PasswordResetHandler {
	return &PasswordResetHandler{
		mailerService:  mailerService,
		userService:    userService,
		store:          store,
		inertiaService: inertiaService,
		appURL:         appURL,
	}
}

// ShowForgotPasswordForm displays the forgot password page
func (h *PasswordResetHandler) ShowForgotPasswordForm(c *fiber.Ctx) error {
	return h.inertiaService.Render(c, "auth/ForgotPassword", fiber.Map{
		"Title": "Forgot Password",
	})
}

// SendResetLink handles sending password reset link
func (h *PasswordResetHandler) SendResetLink(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate email
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Please provide a valid email address",
		})
	}

	// Find user by email
	user, err := h.userService.GetProfileByEmail(req.Email)
	if err != nil {
		// Don't reveal if email exists or not (security best practice)
		return h.inertiaService.Render(c, "auth/ForgotPassword", fiber.Map{
			"success": "If an account exists with that email, we've sent a password reset link.",
		})
	}

	// Generate reset URL
	resetURL := fmt.Sprintf("%s/reset-password/TOKEN_PLACEHOLDER", h.appURL)

	// Send reset email
	err = h.mailerService.SendPasswordResetEmail(user.Email, user.ID, resetURL)
	if err != nil {
		// Log error in production
		return h.inertiaService.Render(c, "auth/ForgotPassword", fiber.Map{
			"success": "If an account exists with that email, we've sent a password reset link.",
		})
	}

	return h.inertiaService.Render(c, "auth/ForgotPassword", fiber.Map{
		"success": "If an account exists with that email, we've sent a password reset link.",
	})
}

// ShowResetPasswordForm displays the reset password page
func (h *PasswordResetHandler) ShowResetPasswordForm(c *fiber.Ctx) error {
	token := c.Params("token")

	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid reset link",
		})
	}

	// Validate token
	_, err := h.mailerService.ValidateResetToken(token)
	if err != nil {
		return h.inertiaService.Render(c, "auth/ResetPassword", fiber.Map{
			"error": "Invalid or expired reset link",
		})
	}

	return h.inertiaService.Render(c, "auth/ResetPassword", fiber.Map{
		"Title": "Reset Password",
		"token": token,
	})
}

// ResetPassword handles password reset
func (h *PasswordResetHandler) ResetPassword(c *fiber.Ctx) error {
	token := c.Params("token")

	var req struct {
		Password          string `json:"password"`
		PasswordConfirmed string `json:"password_confirmation"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate password
	if req.Password == "" || len(req.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 8 characters",
		})
	}

	if req.Password != req.PasswordConfirmed {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Passwords do not match",
		})
	}

	// Validate token
	tokenEntry, err := h.mailerService.ValidateResetToken(token)
	if err != nil {
		return h.inertiaService.Render(c, "auth/ResetPassword", fiber.Map{
			"token": token,
			"error": "Invalid or expired reset link",
		})
	}

	// Get user
	user, err := h.userService.GetProfile(tokenEntry.UserID)
	if err != nil {
		return h.inertiaService.Render(c, "auth/ResetPassword", fiber.Map{
			"token": token,
			"error": "User not found",
		})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return h.inertiaService.Render(c, "auth/ResetPassword", fiber.Map{
			"token": token,
			"error": "Failed to reset password",
		})
	}

	// Update user password
	err = h.userService.UpdatePassword(user.ID, string(hashedPassword))
	if err != nil {
		return h.inertiaService.Render(c, "auth/ResetPassword", fiber.Map{
			"token": token,
			"error": "Failed to reset password",
		})
	}

	// Invalidate token
	h.mailerService.InvalidateResetToken(token)

	return h.inertiaService.Render(c, "auth/ResetPassword", fiber.Map{
		"token":   token,
		"success": "Password reset successfully. You can now login with your new password.",
	})
}

// GetResetTokenEntry is a helper to get token entry (for testing)
func (h *PasswordResetHandler) GetResetTokenEntry(token string) (*services.ResetTokenEntry, error) {
	return h.mailerService.ValidateResetToken(token)
}
