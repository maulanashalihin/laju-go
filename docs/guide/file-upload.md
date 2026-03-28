# File Upload

This guide covers file upload handling, validation, and storage in Laju Go.

## Overview

Laju Go supports file uploads with the following features:

- **Avatar Upload** - User profile picture upload
- **File Validation** - Type and size validation
- **Secure Storage** - Files stored outside web root
- **CSRF Protection** - Upload endpoints protected

## Upload Handler

### Implementation

```go
// app/handlers/upload.go
package handlers

import (
    "fmt"
    "path/filepath"
    "strings"
    "time"
    
    "github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
    uploadPath string
    maxSize    int64
}

func NewUploadHandler() *UploadHandler {
    return &UploadHandler{
        uploadPath: "storage/avatars",
        maxSize:    5 * 1024 * 1024, // 5MB
    }
}

func (h *UploadHandler) Upload(c *fiber.Ctx) error {
    // Parse multipart form (max 5MB)
    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Failed to parse form",
        })
    }
    
    // Get uploaded file
    files := form.File["avatar"]
    if len(files) == 0 {
        return c.Status(400).JSON(fiber.Map{
            "error": "No file uploaded",
        })
    }
    
    file := files[0]
    
    // Validate file type
    allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
    contentType := file.Header.Get("Content-Type")
    if !contains(allowedTypes, contentType) {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid file type. Allowed: JPEG, PNG, GIF",
        })
    }
    
    // Validate file size
    if int64(file.Size) > h.maxSize {
        return c.Status(400).JSON(fiber.Map{
            "error": fmt.Sprintf("File too large. Max size: %dMB", h.maxSize/1024/1024),
        })
    }
    
    // Validate file extension
    ext := strings.ToLower(filepath.Ext(file.Filename))
    allowedExts := []string{".jpg", ".jpeg", ".png", ".gif"}
    if !contains(allowedExts, ext) {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid file extension",
        })
    }
    
    // Generate unique filename
    filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), generateRandomString(10), ext)
    filepath := fmt.Sprintf("%s/%s", h.uploadPath, filename)
    
    // Save file
    if err := c.SaveFile(file, filepath); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to save file",
        })
    }
    
    return c.JSON(fiber.Map{
        "message": "File uploaded successfully",
        "path": filepath,
        "filename": filename,
    })
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

func generateRandomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
    }
    return string(b)
}
```

## Frontend Upload Component

### Avatar Upload Component

```svelte
<!-- frontend/src/components/AvatarUpload.svelte -->
<script>
  import { router, page } from '@inertiajs/svelte';
  
  const user = $page.props.user;
  let avatar = $state(null);
  let preview = $state(user?.avatar || null);
  let uploading = $state(false);
  let error = $state('');
  
  function handleFileChange(event) {
    const file = event.target.files[0];
    if (file) {
      // Validate file size (5MB)
      if (file.size > 5 * 1024 * 1024) {
        error = 'File size must be less than 5MB';
        return;
      }
      
      // Validate file type
      if (!['image/jpeg', 'image/png', 'image/gif'].includes(file.type)) {
        error = 'File must be JPEG, PNG, or GIF';
        return;
      }
      
      error = '';
      avatar = file;
      preview = URL.createObjectURL(file);
    }
  }
  
  function upload() {
    if (!avatar) return;
    
    uploading = true;
    error = '';
    
    const formData = new FormData();
    formData.append('avatar', avatar);
    
    router.post('/upload', formData, {
      onSuccess: () => {
        uploading = false;
        avatar = null;
        // Reload page to show new avatar
        window.location.reload();
      },
      onError: (errors) => {
        uploading = false;
        error = errors.avatar || 'Upload failed';
      },
    });
  }
</script>

<div class="avatar-upload">
  <div class="avatar-preview">
    {#if preview}
      <img src={preview} alt="Avatar" />
    {:else}
      <div class="avatar-placeholder">
        <span>No avatar</span>
      </div>
    {/if}
  </div>
  
  <div class="avatar-controls">
    <input 
      type="file" 
      accept="image/*" 
      onchange={handleFileChange}
      id="avatar-input"
    />
    <label for="avatar-input" class="button button-secondary">
      Choose File
    </label>
    
    {#if avatar}
      <button onclick={upload} disabled={uploading} class="button button-primary">
        {uploading ? 'Uploading...' : 'Upload'}
      </button>
    {/if}
  </div>
  
  {#if error}
    <p class="error">{error}</p>
  {/if}
</div>

<style>
  .avatar-upload {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }
  
  .avatar-preview {
    width: 150px;
    height: 150px;
    border-radius: 50%;
    overflow: hidden;
    border: 2px solid #ddd;
  }
  
  .avatar-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .avatar-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f5f5f5;
    color: #999;
  }
  
  .avatar-controls {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }
  
  #avatar-input {
    display: none;
  }
  
  .button {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.875rem;
  }
  
  .button-secondary {
    background: #6c757d;
    color: white;
  }
  
  .button-primary {
    background: #007bff;
    color: white;
  }
  
  .button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .error {
    color: #dc3545;
    font-size: 0.875rem;
  }
</style>
```

## Storage Configuration

### Environment Variables

```bash
# .env
UPLOAD_PATH=storage/avatars
MAX_UPLOAD_SIZE=5242880
```

### Load Configuration

```go
// app/config/config.go
type Config struct {
    UploadPath   string
    MaxUploadSize int64
}

func Load() *Config {
    uploadPath := getEnv("UPLOAD_PATH", "storage/avatars")
    maxSize := getEnvAsInt("MAX_UPLOAD_SIZE", 5*1024*1024)
    
    return &Config{
        UploadPath: uploadPath,
        MaxUploadSize: maxSize,
    }
}
```

## File Validation

### Validate File Type

```go
func validateFileType(file *multipart.FileHeader) error {
    allowedTypes := map[string]bool{
        "image/jpeg": true,
        "image/png":  true,
        "image/gif":  true,
    }
    
    // Check Content-Type header
    if !allowedTypes[file.Header.Get("Content-Type")] {
        return fmt.Errorf("invalid file type")
    }
    
    return nil
}
```

### Validate File Size

```go
func validateFileSize(file *multipart.FileHeader, maxSize int64) error {
    if int64(file.Size) > maxSize {
        return fmt.Errorf("file too large: %d bytes", file.Size)
    }
    return nil
}
```

### Validate File Extension

```go
func validateFileExtension(filename string, allowedExts []string) error {
    ext := strings.ToLower(filepath.Ext(filename))
    for _, allowed := range allowedExts {
        if ext == allowed {
            return nil
        }
    }
    return fmt.Errorf("invalid file extension: %s", ext)
}
```

### Validate Image Dimensions

```go
import "image"

func validateImageDimensions(filepath string, maxWidth, maxHeight int) error {
    file, err := os.Open(filepath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    config, _, err := image.DecodeConfig(file)
    if err != nil {
        return err
    }
    
    if config.Width > maxWidth {
        return fmt.Errorf("image width exceeds maximum: %d > %d", config.Width, maxWidth)
    }
    
    if config.Height > maxHeight {
        return fmt.Errorf("image height exceeds maximum: %d > %d", config.Height, maxHeight)
    }
    
    return nil
}
```

## File Operations

### Generate Unique Filename

```go
import (
    "crypto/rand"
    "encoding/hex"
)

func generateUniqueFilename(originalFilename string) string {
    // Generate random string
    b := make([]byte, 16)
    rand.Read(b)
    random := hex.EncodeToString(b)
    
    // Get extension
    ext := filepath.Ext(originalFilename)
    
    // Combine timestamp + random + extension
    return fmt.Sprintf("%d_%s%s", time.Now().Unix(), random, ext)
}
```

### Resize Image

```go
import (
    "image"
    "image/jpeg"
    "image/png"
    "github.com/nfnt/resize"
)

func resizeImage(inputPath, outputPath string, width, height uint) error {
    // Open file
    file, err := os.Open(inputPath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    // Decode image
    img, _, err := image.Decode(file)
    if err != nil {
        return err
    }
    
    // Resize
    resized := resize.Resize(width, height, img, resize.Lanczos3)
    
    // Create output file
    out, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer out.Close()
    
    // Encode as JPEG
    return jpeg.Encode(out, resized, &jpeg.Options{Quality: 90})
}
```

### Delete File

```go
func deleteFile(filepath string) error {
    return os.Remove(filepath)
}
```

## Route Setup

```go
// routes/web.go
func SetupRoutes(app *fiber.App) {
    // ... other routes ...
    
    // File upload (protected)
    uploadHandler := handlers.NewUploadHandler()
    protected := app.Group("/", middlewares.AuthRequired(store))
    protected.Post("/upload", uploadHandler.Upload)
}
```

## Security Considerations

### 1. Validate Everything

```go
// ✅ Good: Multiple validation layers
func (h *UploadHandler) Upload(c *fiber.Ctx) error {
    // 1. Validate content type
    if !allowedTypes[file.Header.Get("Content-Type")] {
        return errors.New("invalid type")
    }
    
    // 2. Validate extension
    if !allowedExts[filepath.Ext(file.Filename)] {
        return errors.New("invalid extension")
    }
    
    // 3. Validate size
    if file.Size > h.maxSize {
        return errors.New("file too large")
    }
    
    // 4. Sanitize filename
    filename = sanitizeFilename(file.Filename)
    
    // 5. Store outside web root
    filepath = path.Join(h.uploadPath, filename)
}
```

### 2. Sanitize Filename

```go
import "regexp"

func sanitizeFilename(filename string) string {
    // Remove path separators
    filename = strings.ReplaceAll(filename, "/", "")
    filename = strings.ReplaceAll(filename, "\\", "")
    
    // Remove special characters
    reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
    filename = reg.ReplaceAllString(filename, "_")
    
    // Limit length
    if len(filename) > 100 {
        filename = filename[:100]
    }
    
    return filename
}
```

### 3. Store Outside Web Root

```go
// ✅ Good: Store in non-public directory
uploadPath := "storage/avatars"  // Not in public/

// ❌ Bad: Store in public directory
uploadPath := "public/uploads"  // Accessible via URL
```

### 4. Use Random Filenames

```go
// ✅ Good: Random filename prevents overwriting and guessing
filename := fmt.Sprintf("%d_%s", time.Now().Unix(), randomString(10))

// ❌ Bad: Original filename can be guessed
filename := file.Filename
```

## Multiple File Upload

### Frontend

```svelte
<script>
  let files = $state([]);
  
  function handleFiles(event) {
    files = Array.from(event.target.files);
  }
  
  function uploadAll() {
    const formData = new FormData();
    files.forEach(file => {
      formData.append('files[]', file);
    });
    
    router.post('/upload/multiple', formData, {
      onSuccess: () => {
        files = [];
      },
    });
  }
</script>

<input 
  type="file" 
  multiple 
  accept="image/*"
  onchange={handleFiles}
/>

<button onclick={uploadAll}>
  Upload {files.length} Files
</button>
```

### Backend

```go
func (h *UploadHandler) UploadMultiple(c *fiber.Ctx) error {
    form, err := c.MultipartForm()
    if err != nil {
        return err
    }
    
    files := form.File["files[]"]
    uploaded := []string{}
    
    for _, file := range files {
        // Validate and save each file
        filename := generateUniqueFilename(file.Filename)
        filepath := fmt.Sprintf("%s/%s", h.uploadPath, filename)
        
        if err := c.SaveFile(file, filepath); err != nil {
            return err
        }
        
        uploaded = append(uploaded, filename)
    }
    
    return c.JSON(fiber.Map{
        "message": fmt.Sprintf("Uploaded %d files", len(files)),
        "files": uploaded,
    })
}
```

## Best Practices

### 1. Set Reasonable Limits

```go
// Max file size: 5MB for images
maxSize := 5 * 1024 * 1024

// Max dimensions: 2000x2000 for avatars
maxWidth := 2000
maxHeight := 2000
```

### 2. Clean Up Old Files

```go
// Delete old avatar when uploading new one
func (h *AppHandler) UpdateAvatar(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    
    // Delete old avatar
    if user.Avatar != "" {
        os.Remove(user.Avatar)
    }
    
    // Save new avatar
    // ...
}
```

### 3. Use CDN for Production

```go
// Serve files from CDN in production
func getAvatarURL(filepath string) string {
    if os.Getenv("APP_ENV") == "production" {
        return "https://cdn.example.com/" + filepath
    }
    return "/" + filepath
}
```

### 4. Log Upload Activity

```go
log.Printf("File uploaded: %s (%d bytes) by user %d", 
    filename, file.Size, userID)
```

## Troubleshooting

### File Too Large

**Error**: `http: request body too large`

**Solution**: Increase Fiber's body limit

```go
app := fiber.New(fiber.Config{
    BodyLimit: 10 * 1024 * 1024, // 10MB
})
```

### Multipart Form Parse Error

**Error**: `failed to parse multipart form`

**Solution**: Check memory limit

```go
// Increase memory limit for multipart form
form, err := c.MultipartForm(64 << 20) // 64MB
```

### Permission Denied

**Error**: `permission denied` when saving file

**Solution**: Check directory permissions

```bash
mkdir -p storage/avatars
chmod 755 storage/avatars
chown www-data:www-data storage/avatars
```

## Next Steps

- [Storage Guide](storage.md) - File storage management
- [Security Guide](../reference/security.md) - Security best practices
- [Frontend Guide](frontend.md) - Frontend file upload component
