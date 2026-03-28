package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maulanashalihin/laju-go/app/services"
)

type PublicHandler struct {
	authService    *services.AuthService
	userService    *services.UserService
	inertiaService *services.InertiaService
	assetService   *services.AssetService
}

func NewPublicHandler(authService *services.AuthService, userService *services.UserService, inertiaService *services.InertiaService, assetService *services.AssetService) *PublicHandler {
	return &PublicHandler{
		authService:    authService,
		userService:    userService,
		inertiaService: inertiaService,
		assetService:   assetService,
	}
}

// Index renders the home page
func (h *PublicHandler) Index(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Welcome to Laju",
	}

	// Merge asset data (Vite dev server or production assets)
	for k, v := range h.assetService.GetAssetData() {
		data[k] = v
	}

	return c.Render("index", data)
}

// About renders the about page
func (h *PublicHandler) About(c *fiber.Ctx) error {
	return h.inertiaService.Render(c, "About", fiber.Map{
		"Title": "About Laju",
	})
}
