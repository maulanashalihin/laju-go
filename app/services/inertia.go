package services

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/maulanashalihin/laju-go/app/session"
	"github.com/maulanashalihin/laju-go/templates"
)

// InertiaService provides Inertia.js response helpers
type InertiaService struct {
	assetService *AssetService  // Asset service for production builds
	store        *session.Store // Session store for flash messages
}

// NewInertiaService creates a new InertiaService
func NewInertiaService(assetService *AssetService, store *session.Store) *InertiaService {
	return &InertiaService{
		assetService: assetService,
		store:        store,
	}
}

// Render renders an Inertia response (auto-detect HTML vs JSON)
func (s *InertiaService) Render(c *fiber.Ctx, component string, props fiber.Map) error {
	// Read flash messages from cookies and add to props
	if s.store != nil {
		if flashError := s.store.GetFlash(c, "error"); flashError != "" {
			if props == nil {
				props = fiber.Map{}
			}
			props["flash"] = fiber.Map{
				"error": flashError,
			}
		}
		
		if flashSuccess := s.store.GetFlash(c, "success"); flashSuccess != "" {
			if props == nil {
				props = fiber.Map{}
			}
			if props["flash"] == nil {
				props["flash"] = fiber.Map{}
			}
			props["flash"].(fiber.Map)["success"] = flashSuccess
		}
	}

	// For Inertia requests, return JSON
	if c.Get("X-Inertia") == "true" {
		return s.renderJSON(c, component, props)
	}

	// For initial page load, render HTML template
	return s.renderHTML(c, component, props)
}

// renderJSON renders Inertia JSON response
func (s *InertiaService) renderJSON(c *fiber.Ctx, component string, props fiber.Map) error {
	c.Set("X-Inertia", "true")
	c.Set("X-Inertia-Version", "1.0")
	c.Set("Vary", "X-Inertia")
	c.Set("Content-Type", "application/json")

	return c.JSON(fiber.Map{
		"component": component,
		"props":     props,
		"url":       c.OriginalURL(),
	})
}

// renderHTML renders initial HTML page load
func (s *InertiaService) renderHTML(c *fiber.Ctx, component string, props fiber.Map) error {
	pageData, _ := json.Marshal(fiber.Map{
		"component": component,
		"props":     props,
		"url":       c.OriginalURL(),
	})

	title, _ := props["Title"].(string)
	isDev := s.assetService.IsDevelopment()
	mainJS := s.assetService.GetMainJS()
	mainCSS := s.assetService.GetMainCSS()

	c.Set("Content-Type", "text/html; charset=utf-8")
	return templates.InertiaPage(title, string(pageData), isDev, mainJS, mainCSS, nil).Render(c.Context(), c.Response().BodyWriter())
}

// RenderWithMeta renders an Inertia response with additional metadata
func (s *InertiaService) RenderWithMeta(c *fiber.Ctx, component string, props fiber.Map, meta fiber.Map) error {
	if c.Get("X-Inertia") == "true" {
		c.Set("X-Inertia", "true")
		c.Set("X-Inertia-Version", "1.0")
		c.Set("Vary", "X-Inertia")
		c.Set("Content-Type", "application/json")

		response := fiber.Map{
			"component": component,
			"props":     props,
			"url":       c.OriginalURL(),
		}

		if meta != nil {
			response["meta"] = meta
		}

		return c.JSON(response)
	}

	pageData, _ := json.Marshal(fiber.Map{
		"component": component,
		"props":     props,
		"url":       c.OriginalURL(),
		"meta":      meta,
	})

	title, _ := props["Title"].(string)

	c.Set("Content-Type", "text/html; charset=utf-8")
	return templates.InertiaPage(title, string(pageData), false, "", "", nil).Render(c.Context(), c.Response().BodyWriter())
}
