package handlers

import (
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/maulanashalihin/laju-go/app/session"
)

type UploadHandler struct {
	store *session.Store
}

func NewUploadHandler(store *session.Store) *UploadHandler {
	return &UploadHandler{
		store: store,
	}
}

// Upload handles file uploads
func (h *UploadHandler) Upload(c *fiber.Ctx) error {
	sess, _ := h.store.Get(c)
	userID := sess.Get("user_id")

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form",
		})
	}

	// Get the file from the form
	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	file := files[0]

	// Validate file type
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
	contentType := file.Header.Get("Content-Type")
	isAllowed := false
	for _, allowed := range allowedTypes {
		if contentType == allowed {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Allowed: JPEG, PNG, GIF, WEBP",
		})
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File too large. Max size: 5MB",
		})
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := generateFilename(userID.(int64), ext)

	// Save the file
	uploadPath := filepath.Join("storage", "avatars", filename)
	if err := c.SaveFile(file, uploadPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Return the file URL
	fileURL := "/storage/avatars/" + filename

	return c.JSON(fiber.Map{
		"message":  "File uploaded successfully",
		"file_url": fileURL,
	})
}

// generateFilename generates a unique filename for uploaded files
func generateFilename(userID int64, ext string) string {
	// In production, use a proper unique ID generator
	return strings.ReplaceAll(filepath.Base(filepath.Join("user", string(rune(userID)), ext)), "/", "_")
}
