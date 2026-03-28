# Configuration

This guide covers all environment variables and configuration options for Laju Go.

## Environment File

Laju Go uses a `.env` file in the project root for configuration. The application loads these variables at startup using `godotenv`.

### File Location

```
laju-go/
├── .env              # Your configuration (create from .env.example)
├── .env.example      # Template file (commit-safe)
└── .env.production   # Optional: Production-specific overrides
```

> ⚠️ **Never commit `.env` to version control!** It's gitignored by default.

## Required Variables

These variables must be set for the application to run:

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `APP_ENV` | string | `development` | Environment mode (`development`, `production`, `test`) |
| `APP_PORT` | int | `8080` | HTTP server port |
| `DB_PATH` | string | `data/app.db` | Path to SQLite database file |
| `SESSION_SECRET` | string | **required** | Secret key for session encryption (min 32 chars) |

### Example: Minimum Configuration

```bash
# Application
APP_ENV=development
APP_PORT=8080

# Database
DB_PATH=data/app.db

# Session (generate a secure random string)
SESSION_SECRET=super-secret-key-at-least-32-characters-long
```

## Optional Variables

### Google OAuth

Enable Google OAuth 2.0 authentication:

| Variable | Type | Description |
|----------|------|-------------|
| `GOOGLE_CLIENT_ID` | string | OAuth 2.0 Client ID from Google Cloud Console |
| `GOOGLE_CLIENT_SECRET` | string | OAuth 2.0 Client Secret |
| `GOOGLE_REDIRECT_URL` | string | OAuth callback URL (e.g., `http://localhost:8080/auth/google/callback`) |

### Email (SMTP)

Enable password reset via email:

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `SMTP_HOST` | string | `localhost` | SMTP server hostname |
| `SMTP_PORT` | int | `587` | SMTP server port |
| `SMTP_USER` | string | `` | SMTP username (email address) |
| `SMTP_PASS` | string | `` | SMTP password (or app password) |
| `FROM_EMAIL` | string | `noreply@example.com` | Sender email address |
| `FROM_NAME` | string | `Laju Go` | Sender name |

### Application URL

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `APP_URL` | string | `http://localhost:8080` | Base URL for email links and redirects |

### File Upload

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `UPLOAD_PATH` | string | `storage/avatars` | Directory for uploaded files |
| `MAX_UPLOAD_SIZE` | int | `5242880` | Maximum file size in bytes (default: 5MB) |

### Rate Limiting

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `RATE_LIMIT_AUTH` | int | `5` | Max auth requests per 15 minutes |
| `RATE_LIMIT_PASSWORD_RESET` | int | `3` | Max password reset requests per hour |

## Complete Example

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
SESSION_SECRET=your-32-character-secret-key-change-in-production

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

## Environment-Specific Configuration

### Development

```bash
APP_ENV=development
APP_PORT=8080
DB_PATH=data/app.db
SESSION_SECRET=dev-secret-key-not-for-production-use

# Enable verbose logging (if implemented)
# LOG_LEVEL=debug
```

### Production

```bash
APP_ENV=production
APP_PORT=8080
DB_PATH=/var/lib/laju/app.db
SESSION_SECRET=<strong-random-secret-generated-by-openssl>

# Production SMTP example (Gmail)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=noreply@yourdomain.com
SMTP_PASS=<app-password>
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your App Name

# Production OAuth
GOOGLE_CLIENT_ID=<production-client-id>
GOOGLE_CLIENT_SECRET=<production-secret>
GOOGLE_REDIRECT_URL=https://yourdomain.com/auth/google/callback
```

### Testing

```bash
APP_ENV=test
APP_PORT=3000
DB_PATH=data/test.db
SESSION_SECRET=test-secret-key
```

## Generating Secure Secrets

### Session Secret

```bash
# Using openssl
openssl rand -base64 32

# Using Go
go run -e 'package main; import ("crypto/rand"; "encoding/base64"; "fmt"); func main() { b := make([]byte, 32); rand.Read(b); fmt.Println(base64.StdEncoding.EncodeToString(b)) }'

# Using Python
python3 -c "import secrets; print(secrets.token_urlsafe(32))"
```

### Example Output

```
x7K9mP2vL5nQ8wR3tY6uI0oA4sD7fG1hJ9kM2xC5vB8n
```

## Google OAuth Setup

### Step 1: Create Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Click **Create Project**
3. Enter project name (e.g., "Laju Go App")
4. Click **Create**

### Step 2: Enable Google+ API

1. In the Google Cloud Console, go to **APIs & Services** → **Library**
2. Search for "Google+ API"
3. Click on it and press **Enable**

### Step 3: Create OAuth Credentials

1. Go to **APIs & Services** → **Credentials**
2. Click **Create Credentials** → **OAuth client ID**
3. Select **Web application**
4. Add authorized redirect URIs:
   - Development: `http://localhost:8080/auth/google/callback`
   - Production: `https://yourdomain.com/auth/google/callback`
5. Click **Create**

### Step 4: Copy Credentials

Copy the **Client ID** and **Client Secret** to your `.env` file.

## SMTP Configuration Examples

### Gmail

```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-16-character-app-password
FROM_EMAIL=noreply@example.com
FROM_NAME=Laju Go
```

> 🔐 **Important**: Use an [App Password](https://support.google.com/accounts/answer/185833), not your regular Gmail password.

### SendGrid

```bash
SMTP_HOST=smtp.sendgrid.com
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=your-sendgrid-api-key
FROM_EMAIL=verified-sender@yourdomain.com
FROM_NAME=Laju Go
```

### Mailgun

```bash
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USER=postmaster@yourdomain.mailgun.org
SMTP_PASS=your-mailgun-api-key
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Laju Go
```

### Office 365

```bash
SMTP_HOST=smtp.office365.com
SMTP_PORT=587
SMTP_USER=your-email@yourdomain.com
SMTP_PASS=your-password
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Laju Go
```

## Database Path Configuration

### Development

```bash
DB_PATH=data/app.db
```

### Production (Linux)

```bash
DB_PATH=/var/lib/laju/app.db
```

Ensure the directory exists and is writable:

```bash
sudo mkdir -p /var/lib/laju
sudo chown www-data:www-data /var/lib/laju
sudo chmod 755 /var/lib/laju
```

### Docker

```bash
DB_PATH=/root/data/app.db
```

Mount a volume to persist data:

```bash
docker run -v $(pwd)/data:/root/data laju-go
```

## File Upload Configuration

### Default Upload Path

```bash
UPLOAD_PATH=storage/avatars
```

### Maximum File Size

```bash
# 5MB (default)
MAX_UPLOAD_SIZE=5242880

# 10MB
MAX_UPLOAD_SIZE=10485760

# 1MB
MAX_UPLOAD_SIZE=1048576
```

## Rate Limiting Configuration

### Authentication Endpoints

```bash
# 5 requests per 15 minutes (default)
RATE_LIMIT_AUTH=5

# More restrictive: 3 requests per 15 minutes
RATE_LIMIT_AUTH=3

# Less restrictive: 10 requests per 15 minutes
RATE_LIMIT_AUTH=10
```

### Password Reset Endpoints

```bash
# 3 requests per hour (default)
RATE_LIMIT_PASSWORD_RESET=3

# More restrictive: 2 requests per hour
RATE_LIMIT_PASSWORD_RESET=2
```

## Configuration Validation

The application validates configuration at startup. If required variables are missing, you'll see an error like:

```
Error: SESSION_SECRET is required in production mode
```

## Best Practices

### 1. Use Strong Secrets

```bash
# ❌ Bad: Weak, predictable secret
SESSION_SECRET=secret123

# ✅ Good: Strong, random secret
SESSION_SECRET=x7K9mP2vL5nQ8wR3tY6uI0oA4sD7fG1hJ9kM2xC5vB8n
```

### 2. Separate Environment Files

Keep different `.env` files for different environments:

```bash
.env              # Local development (gitignored)
.env.example      # Template (commit this)
.env.production   # Production server (gitignored)
.env.test         # Testing (gitignored)
```

### 3. Use Environment Variables in Production

For containerized deployments, inject variables via environment:

```bash
# Docker Compose
environment:
  - SESSION_SECRET=${SESSION_SECRET}
  - DB_PATH=/var/lib/laju/app.db
```

### 4. Rotate Secrets Regularly

Change `SESSION_SECRET` periodically:

1. Generate new secret
2. Update `.env`
3. Restart application
4. All existing sessions will be invalidated (users must re-login)

### 5. Secure Production Files

```bash
# Set restrictive permissions on .env
chmod 600 .env
chown www-data:www-data .env
```

## Troubleshooting

### Session Not Persisting

**Problem**: Users are logged out immediately

**Solution**: Check `SESSION_SECRET` is set and consistent

```bash
# Verify .env is loaded
grep SESSION_SECRET .env
```

### OAuth Not Working

**Problem**: Google OAuth returns "redirect_uri_mismatch"

**Solution**: Ensure `GOOGLE_REDIRECT_URL` exactly matches Google Cloud Console settings

### Email Not Sending

**Problem**: Password reset emails not received

**Solution**: 
1. Verify SMTP credentials
2. Check firewall (port 587)
3. For Gmail, use App Password (not regular password)
4. Check spam folder

### Database Path Errors

**Problem**: `unable to open database file`

**Solution**: Ensure directory exists and is writable

```bash
mkdir -p data
chmod 755 data
```

## Next Steps

- [Architecture Guide](../guide/architecture.md) - Understanding the codebase
- [Development Workflow](../deployment/development.md) - Hot reload and scripts
- [Production Deployment](../deployment/production.md) - Deploying to production servers
