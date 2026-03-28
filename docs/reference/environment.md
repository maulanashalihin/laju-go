# Environment Variables

Complete reference for all environment variables in Laju Go.

## Overview

Laju Go uses environment variables for configuration. The application loads these from a `.env` file in the project root or from system environment variables.

## Required Variables

These variables must be set for the application to run:

| Variable | Type | Default | Required | Description |
|----------|------|---------|----------|-------------|
| `APP_ENV` | string | `development` | No | Environment mode (`development`, `production`, `test`) |
| `APP_PORT` | int | `8080` | No | HTTP server port |
| `DB_PATH` | string | `data/app.db` | No | Path to SQLite database file |
| `SESSION_SECRET` | string | - | **Yes** | Secret key for session encryption (min 32 characters) |

### Example: Minimum Configuration

```bash
APP_ENV=development
APP_PORT=8080
DB_PATH=data/app.db
SESSION_SECRET=your-32-character-secret-key
```

## Application Variables

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `APP_ENV` | string | `development` | Environment mode |
| `APP_PORT` | int | `8080` | HTTP server port |
| `APP_URL` | string | `http://localhost:8080` | Base URL for links and redirects |

### APP_ENV

Environment mode affects application behavior:

- `development` - Verbose logging, hot reload, detailed errors
- `production` - Optimized, minimal logging, secure defaults
- `test` - Test mode for automated testing

```bash
# Development
APP_ENV=development

# Production
APP_ENV=production

# Testing
APP_ENV=test
```

### APP_PORT

Port for the HTTP server:

```bash
# Default
APP_PORT=8080

# Alternative
APP_PORT=3000
```

### APP_URL

Base URL for generating links:

```bash
# Development
APP_URL=http://localhost:8080

# Production
APP_URL=https://yourdomain.com
```

Used in:
- Password reset links
- Email templates
- OAuth redirect URLs

## Database Variables

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `DB_PATH` | string | `data/app.db` | Path to SQLite database file |

### DB_PATH

Location of the SQLite database file:

```bash
# Development
DB_PATH=data/app.db

# Production (Linux)
DB_PATH=/var/lib/laju/app.db

# Docker
DB_PATH=/root/data/app.db
```

## Session Variables

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `SESSION_SECRET` | string | - | **Required** - Secret key for session encryption |

### SESSION_SECRET

**Critical for security** - Used to encrypt session cookies.

**Requirements**:
- Minimum 32 characters
- Use cryptographically secure random generation
- Different value for each environment
- Rotate periodically

**Generate Secure Secret**:

```bash
# Using openssl
openssl rand -base64 32

# Output example:
# x7K9mP2vL5nQ8wR3tY6uI0oA4sD7fG1hJ9kM2xC5vB8n

# Using Go
go run -e 'package main; import ("crypto/rand"; "encoding/base64"; "fmt"); func main() { b := make([]byte, 32); rand.Read(b); fmt.Println(base64.StdEncoding.EncodeToString(b)) }'
```

**Usage**:

```bash
# ✅ Good: Strong random secret
SESSION_SECRET=x7K9mP2vL5nQ8wR3tY6uI0oA4sD7fG1hJ9kM2xC5vB8n

# ❌ Bad: Weak predictable secret
SESSION_SECRET=secret123
```

## Google OAuth Variables

Optional - Required only if using Google OAuth login.

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `GOOGLE_CLIENT_ID` | string | - | OAuth 2.0 Client ID from Google Cloud |
| `GOOGLE_CLIENT_SECRET` | string | - | OAuth 2.0 Client Secret |
| `GOOGLE_REDIRECT_URL` | string | `http://localhost:8080/auth/google/callback` | OAuth callback URL |

### Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add redirect URIs:
   - Development: `http://localhost:8080/auth/google/callback`
   - Production: `https://yourdomain.com/auth/google/callback`

### Example

```bash
GOOGLE_CLIENT_ID=123456789-abc123def456.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-abcdefghijklmnop
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
```

## Email/SMTP Variables

Optional - Required only if sending emails (password reset, welcome emails, etc.).

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `SMTP_HOST` | string | `localhost` | SMTP server hostname |
| `SMTP_PORT` | int | `587` | SMTP server port |
| `SMTP_USER` | string | - | SMTP username |
| `SMTP_PASS` | string | - | SMTP password |
| `FROM_EMAIL` | string | `noreply@example.com` | Sender email address |
| `FROM_NAME` | string | `Laju Go` | Sender name |

### Gmail Configuration

```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-16-character-app-password
FROM_EMAIL=noreply@example.com
FROM_NAME=Your App Name
```

> 🔐 **Important**: Use an [App Password](https://support.google.com/accounts/answer/185833), not your regular Gmail password.

### SendGrid Configuration

```bash
SMTP_HOST=smtp.sendgrid.com
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=your-sendgrid-api-key
FROM_EMAIL=verified-sender@yourdomain.com
FROM_NAME=Your App Name
```

### Mailgun Configuration

```bash
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USER=postmaster@yourdomain.mailgun.org
SMTP_PASS=your-mailgun-api-key
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your App Name
```

## File Upload Variables

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `UPLOAD_PATH` | string | `storage/avatars` | Directory for uploaded files |
| `MAX_UPLOAD_SIZE` | int | `5242880` | Maximum file size in bytes (5MB) |

### Example

```bash
UPLOAD_PATH=storage/avatars
MAX_UPLOAD_SIZE=5242880
```

## Rate Limiting Variables

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `RATE_LIMIT_AUTH` | int | `5` | Max auth requests per 15 minutes |
| `RATE_LIMIT_PASSWORD_RESET` | int | `3` | Max password reset requests per hour |

### Example

```bash
RATE_LIMIT_AUTH=5
RATE_LIMIT_PASSWORD_RESET=3
```

## Complete Example

### Development

```bash
# ===========================================
# Application Configuration
# ===========================================
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080

# ===========================================
# Database Configuration
# ===========================================
DB_PATH=data/app.db

# ===========================================
# Session Configuration
# ===========================================
SESSION_SECRET=dev-secret-key-not-for-production-use-only

# ===========================================
# Google OAuth Configuration (Optional)
# ===========================================
GOOGLE_CLIENT_ID=123456789-abc123def456.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-abcdefghijklmnop
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# ===========================================
# Email/SMTP Configuration (Optional)
# ===========================================
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-specific-password
FROM_EMAIL=noreply@example.com
FROM_NAME=Laju Go

# ===========================================
# File Upload Configuration
# ===========================================
UPLOAD_PATH=storage/avatars
MAX_UPLOAD_SIZE=5242880

# ===========================================
# Rate Limiting Configuration
# ===========================================
RATE_LIMIT_AUTH=5
RATE_LIMIT_PASSWORD_RESET=3
```

### Production

```bash
# ===========================================
# Application Configuration
# ===========================================
APP_ENV=production
APP_PORT=8080
APP_URL=https://yourdomain.com

# ===========================================
# Database Configuration
# ===========================================
DB_PATH=/var/lib/laju/app.db

# ===========================================
# Session Configuration
# ===========================================
SESSION_SECRET=<run: openssl rand -base64 32>

# ===========================================
# Google OAuth Configuration
# ===========================================
GOOGLE_CLIENT_ID=<production-client-id>
GOOGLE_CLIENT_SECRET=<production-secret>
GOOGLE_REDIRECT_URL=https://yourdomain.com/auth/google/callback

# ===========================================
# Email/SMTP Configuration
# ===========================================
SMTP_HOST=smtp.sendgrid.com
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=<sendgrid-api-key>
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your App Name

# ===========================================
# File Upload Configuration
# ===========================================
UPLOAD_PATH=/var/lib/laju/storage/avatars
MAX_UPLOAD_SIZE=5242880

# ===========================================
# Rate Limiting Configuration
# ===========================================
RATE_LIMIT_AUTH=5
RATE_LIMIT_PASSWORD_RESET=3
```

### Docker

```bash
# ===========================================
# Application Configuration
# ===========================================
APP_ENV=production
APP_PORT=8080
APP_URL=https://yourdomain.com

# ===========================================
# Database Configuration
# ===========================================
DB_PATH=/root/data/app.db

# ===========================================
# Session Configuration
# ===========================================
SESSION_SECRET=${SESSION_SECRET}

# ===========================================
# Email/SMTP Configuration
# ===========================================
SMTP_HOST=${SMTP_HOST}
SMTP_PORT=${SMTP_PORT}
SMTP_USER=${SMTP_USER}
SMTP_PASS=${SMTP_PASS}
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your App Name
```

## Loading Environment Variables

### From .env File

The application automatically loads `.env` from the project root:

```go
import "github.com/joho/godotenv"

func init() {
    godotenv.Load()
}
```

### From System Environment

Environment variables can also be set directly:

```bash
export SESSION_SECRET=your-secret-key
export APP_PORT=8080
./laju-go
```

### In Docker Compose

```yaml
services:
  app:
    environment:
      - SESSION_SECRET=${SESSION_SECRET}
      - APP_PORT=8080
      - DB_PATH=/root/data/app.db
    env_file:
      - .env
```

## Validation

The application validates required variables at startup:

```go
func validateConfig() {
    sessionSecret := os.Getenv("SESSION_SECRET")
    if sessionSecret == "" {
        log.Fatal("SESSION_SECRET is required")
    }
    
    if len(sessionSecret) < 32 {
        log.Fatal("SESSION_SECRET must be at least 32 characters")
    }
}
```

## Security Best Practices

### 1. Never Commit .env

```bash
# ✅ Good: .env is in .gitignore
cat .gitignore | grep .env
# Output: .env

# ❌ Bad: Committing .env with secrets
git add .env  # Never do this!
```

### 2. Use .env.example

Create a template file with placeholder values:

```bash
# .env.example
APP_ENV=development
APP_PORT=8080
DB_PATH=data/app.db
SESSION_SECRET=your-32-character-secret-key
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
```

### 3. Use Strong Secrets

```bash
# ✅ Good: Generate secure secret
openssl rand -base64 32

# ❌ Bad: Weak predictable secret
SESSION_SECRET=password123
```

### 4. Rotate Secrets Regularly

```bash
# Generate new secret
openssl rand -base64 32

# Update .env
nano .env

# Restart application
sudo systemctl restart laju-go
```

### 5. Use Different Secrets per Environment

```bash
# Development
SESSION_SECRET=dev-secret-key

# Production
SESSION_SECRET=<different-secure-key>

# Test
SESSION_SECRET=test-secret-key
```

## Troubleshooting

### Session Not Persisting

**Problem**: Users logged out after restart

**Solution**: Check `SESSION_SECRET` is set and consistent

```bash
grep SESSION_SECRET .env
```

### OAuth Not Working

**Problem**: Google OAuth fails

**Solution**: Verify all OAuth variables are set

```bash
grep GOOGLE_ .env
```

### Email Not Sending

**Problem**: Password reset emails not received

**Solution**: Check SMTP configuration

```bash
grep SMTP_ .env
```

### Port Already in Use

**Problem**: `listen tcp :8080: bind: address already in use`

**Solution**: Change port in `.env`

```bash
APP_PORT=8081
```

## Quick Reference

```bash
# Minimum required
SESSION_SECRET=<32+ characters>

# Common configuration
APP_ENV=development
APP_PORT=8080
DB_PATH=data/app.db

# Optional features
GOOGLE_CLIENT_ID=<client-id>
GOOGLE_CLIENT_SECRET=<secret>
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=<user>
SMTP_PASS=<password>
```

## Next Steps

- [Configuration Guide](../getting-started/configuration.md) - Setup guide
- [Security Guide](security.md) - Security best practices
- [Deployment Guide](../deployment/production.md) - Production configuration
