package handlers

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maulanashalihin/laju-go/app/services"
	"github.com/maulanashalihin/laju-go/app/session"
)

type UploadHandler struct {
	store       *session.Store
	userService *services.UserService
}

func NewUploadHandler(store *session.Store, userService *services.UserService) *UploadHandler {
	return &UploadHandler{
		store:       store,
		userService: userService,
	}
}

// Upload handles file uploads
func (h *UploadHandler) Upload(c *fiber.Ctx) error {
	sess, _ := h.store.Get(c)
	userID := sess.Get("user_id")

	if userID == nil {
		log.Printf("[Upload] User not authenticated\n")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	log.Printf("[Upload] User ID: %v\n", userID)

	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("[Upload] Failed to parse form: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form",
		})
	}

	// Get the file from the form
	files := form.File["file"]
	if len(files) == 0 {
		log.Printf("[Upload] No file uploaded\n")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	file := files[0]
	log.Printf("[Upload] File: %s, Size: %d, Type: %s\n", file.Filename, file.Size, file.Header.Get("Content-Type"))

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
		log.Printf("[Upload] Invalid file type: %s\n", contentType)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Allowed: JPEG, PNG, GIF, WEBP",
		})
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		log.Printf("[Upload] File too large: %d bytes\n", file.Size)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File too large. Max size: 5MB",
		})
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", userID.(int64), time.Now().UnixNano(), ext)

	// Save the file
	uploadPath := filepath.Join("storage", "avatars", filename)
	log.Printf("[Upload] Saving to: %s\n", uploadPath)
	
	if err := c.SaveFile(file, uploadPath); err != nil {
		log.Printf("[Upload] Failed to save file: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Update user avatar in database
	avatarURL := "/storage/avatars/" + filename
	log.Printf("[Upload] Updating avatar to: %s\n", avatarURL)
	
	if err := h.userService.UpdateAvatar(userID.(int64), avatarURL); err != nil {
		log.Printf("[Upload] Failed to update avatar in DB: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update avatar",
		})
	}

	log.Printf("[Upload] Success: %s\n", avatarURL)

	// Return the file URL
	return c.JSON(fiber.Map{
		"success": true,
		"url":     avatarURL,
		"message": "File uploaded successfully",
	})
}
