package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maulanashalihin/laju-go/app/models"
	"github.com/maulanashalihin/laju-go/app/services"
	"github.com/maulanashalihin/laju-go/app/session"
)

type AppHandler struct {
	userService    *services.UserService
	store          *session.Store
	inertiaService *services.InertiaService
}

func NewAppHandler(userService *services.UserService, store *session.Store, inertiaService *services.InertiaService) *AppHandler {
	return &AppHandler{
		userService:    userService,
		store:          store,
		inertiaService: inertiaService,
	}
}

// Dashboard renders the main app dashboard using Inertia
func (h *AppHandler) Dashboard(c *fiber.Ctx) error {
	// Get user info from locals (set by AuthRequired middleware)
	userID := c.Locals("user_id")

	if userID == nil {
		// Should not happen as AuthRequired middleware handles this
		return c.Redirect("/login")
	}

	user, err := h.userService.GetProfile(userID.(int64))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load dashboard",
		})
	}

	return h.inertiaService.Render(c, "app/Dashboard", fiber.Map{
		"user": user,
	})
}

// Profile returns user profile (Inertia)
func (h *AppHandler) Profile(c *fiber.Ctx) error {
	// Get user info from locals (set by AuthRequired middleware)
	userID := c.Locals("user_id")

	if userID == nil {
		// Should not happen as AuthRequired middleware handles this
		return c.Redirect("/login")
	}

	user, err := h.userService.GetProfile(userID.(int64))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load profile",
		})
	}

	return h.inertiaService.Render(c, "app/Profile", fiber.Map{
		"user": user,
	})
}

// UpdateProfile updates user profile (Inertia)
func (h *AppHandler) UpdateProfile(c *fiber.Ctx) error {
	// Get user info from locals (set by AuthRequired middleware)
	userID := c.Locals("user_id")

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	var req models.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.userService.UpdateProfile(userID.(int64), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return h.inertiaService.Render(c, "app/Profile", fiber.Map{
		"user":    user,
		"success": "Profile updated successfully",
	})
}

// UpdatePassword updates user password (Inertia)
func (h *AppHandler) UpdatePassword(c *fiber.Ctx) error {
	// Get user info from locals (set by AuthRequired middleware)
	userID := c.Locals("user_id")

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate passwords
	if req.NewPassword != req.ConfirmPassword {
		return h.inertiaService.Render(c, "app/Profile", fiber.Map{
			"error": "Passwords do not match",
		})
	}

	if len(req.NewPassword) < 8 {
		return h.inertiaService.Render(c, "app/Profile", fiber.Map{
			"error": "Password must be at least 8 characters",
		})
	}

	// Change password
	err := h.userService.ChangePassword(userID.(int64), req.CurrentPassword, req.NewPassword)
	if err != nil {
		return h.inertiaService.Render(c, "app/Profile", fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.userService.GetProfile(userID.(int64))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load profile",
		})
	}

	return h.inertiaService.Render(c, "app/Profile", fiber.Map{
		"user":    user,
		"success": "Password changed successfully",
	})
}
