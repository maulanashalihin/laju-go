package services

import (
	"encoding/json"
	"html/template"

	"github.com/gofiber/fiber/v2"
)

// InertiaService provides Inertia.js response helpers
type InertiaService struct {
	template     string        // template name for initial page load
	assetService *AssetService // Asset service for production builds
}

// NewInertiaService creates a new InertiaService
func NewInertiaService(assetService *AssetService) *InertiaService {
	return &InertiaService{
		template:     "inertia",
		assetService: assetService,
	}
}

// Render renders an Inertia response (auto-detect HTML vs JSON)
func (s *InertiaService) Render(c *fiber.Ctx, component string, props fiber.Map) error {
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
	// Marshal page data to JSON string for template
	pageData, _ := json.Marshal(fiber.Map{
		"component": component,
		"props":     props,
		"url":       c.OriginalURL(),
	})

	templateData := fiber.Map{
		"Title":     props["Title"],
		"Component": component,
		"Page":      template.JS(string(pageData)),
	}

	// Merge asset data (Vite dev server or production assets)
	for k, v := range s.assetService.GetAssetData() {
		templateData[k] = v
	}

	return c.Render(s.template, templateData)
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

	return c.Render(s.template, fiber.Map{
		"Title":     props["Title"],
		"Component": component,
		"Page":      template.JS(string(pageData)),
	})
}
