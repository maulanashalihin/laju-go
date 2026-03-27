package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/velostack/velostack-go/app/models"
	"github.com/velostack/velostack-go/app/services"
	"github.com/velostack/velostack-go/app/session"
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
	sess, _ := h.store.Get(c)
	userID := sess.Get("user_id")

	if userID == nil {
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
	sess, _ := h.store.Get(c)
	userID := sess.Get("user_id")

	if userID == nil {
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
	sess, _ := h.store.Get(c)
	userID := sess.Get("user_id")

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

	return h.inertiaService.Render(c, "Profile", fiber.Map{
		"user":    user,
		"success": "Profile updated successfully",
	})
}
