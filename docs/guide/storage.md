# Storage

This guide covers file storage in Laju Go.

## Overview

Laju Go stores uploaded files on the local filesystem under the `storage/` directory.

```
storage/
└── avatars/                   # User avatar uploads
```

## Directory Structure

| Path | Purpose | Served At |
|------|---------|-----------|
| `storage/` | All uploaded files | `/storage/*` |
| `storage/avatars/` | Profile pictures | `/storage/avatars/*` |
| `public/` | Static assets (shipped with code) | `/public/*` |

Static file serving is configured in `routes/web.go`:

```go
app.Static("/dist", "./dist")          # Built frontend assets
app.Static("/public", "./public")      # Static assets (images, favicon)
app.Static("/storage", "./storage")    # Uploaded files
```

## File Upload

### Handler

```go
// app/handlers/upload.go
func (h *UploadHandler) Upload(c *fiber.Ctx) error {
    // Get file from multipart form
    file, err := c.FormFile("avatar")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "No file provided"})
    }

    // Validate file type
    ext := strings.ToLower(filepath.Ext(file.Filename))
    if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid file type"})
    }

    // Generate unique filename
    filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
    path := filepath.Join("./storage/avatars", filename)

    // Save file
    if err := c.SaveFile(file, path); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to save file"})
    }

    // Update user avatar in DB
    avatarURL := "/storage/avatars/" + filename
    // ... update user record ...
}
```

### Avatar Proxy

For OAuth users (Google), avatars are proxied from external URLs:

```go
// GET /api/avatar/:id
func (h *AuthHandler) GetAvatar(c *fiber.Ctx) error {
    user, _ := h.authService.GetUserByID(userID)

    if strings.HasPrefix(user.Avatar, "/storage/") {
        return c.SendFile("." + user.Avatar)  // Local file
    }

    // Proxy from external URL (Google)
    resp, _ := http.Get(user.Avatar)
    // ...
}
```

## Backups

Backup `storage/` along with the database:

```bash
tar -czf backup.tar.gz data/ storage/
```

## Next Steps

- [File Upload Guide](file-upload.md) — Detailed upload handling
- [Data Protection Guide](../guide/data-protection.md) — Backup strategies
