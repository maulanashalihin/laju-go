package handlers

import (
	"fmt"
	"log/slog"
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
		slog.Warn("upload user not authenticated", "handler", "Upload")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	slog.Info("upload user ID", "handler", "Upload", "user_id", userID)

	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		slog.Error("upload failed to parse form", "handler", "Upload", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form",
		})
	}

	// Get the file from the form
	files := form.File["file"]
	if len(files) == 0 {
		slog.Info("upload no file uploaded", "handler", "Upload")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	file := files[0]
	slog.Info("upload file info", "handler", "Upload", "filename", file.Filename, "size", file.Size, "content_type", file.Header.Get("Content-Type"))

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
		slog.Warn("upload invalid file type", "handler", "Upload", "content_type", contentType)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Allowed: JPEG, PNG, GIF, WEBP",
		})
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		slog.Warn("upload file too large", "handler", "Upload", "size", file.Size)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File too large. Max size: 5MB",
		})
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", userID.(int64), time.Now().UnixNano(), ext)

	// Save the file
	uploadPath := filepath.Join("storage", "avatars", filename)
	slog.Info("upload saving file", "handler", "Upload", "path", uploadPath)
	
	if err := c.SaveFile(file, uploadPath); err != nil {
		slog.Error("upload failed to save file", "handler", "Upload", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Update user avatar in database
	avatarURL := "/storage/avatars/" + filename
	slog.Info("upload updating avatar", "handler", "Upload", "avatar_url", avatarURL)
	
	if err := h.userService.UpdateAvatar(userID.(int64), avatarURL); err != nil {
		slog.Error("upload failed to update avatar in DB", "handler", "Upload", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update avatar",
		})
	}

	// Sync session with new avatar URL
	sess.Set("avatar", avatarURL)
	sess.Save()

	slog.Info("upload success", "handler", "Upload", "avatar_url", avatarURL)

	// Return the file URL
	return c.JSON(fiber.Map{
		"success": true,
		"url":     avatarURL,
		"message": "File uploaded successfully",
	})
}
