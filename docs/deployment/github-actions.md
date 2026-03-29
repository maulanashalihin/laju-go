# GitHub Actions CI/CD Setup Guide

This guide covers setting up automated deployment using GitHub Actions.

## Overview

```
┌─────────────────┐      ┌──────────────────┐      ┌─────────────────┐
│  Push to main   │      │  GitHub Actions  │      │  Production VPS │
│                 │─────▶│                  │─────▶│                 │
│  git push       │      │  - Build binary  │      │  - Binary only  │
│                 │      │  - Run tests     │      │  - No Go/Node   │
│                 │      │  - Deploy via SSH│      │  - Run service  │
└─────────────────┘      └──────────────────┘      └─────────────────┘
```

**Benefits:**
- ✅ No Go/Node.js needed on production server
- ✅ Automated testing on pull requests
- ✅ Consistent builds in clean environment
- ✅ One-command deployment (just `git push`)

---

## Step 1: Generate SSH Key

Create a dedicated SSH key for GitHub Actions:

```bash
# Generate ED25519 key (more secure than RSA)
ssh-keygen -t ed25519 -C "github-actions@laju-go" -f ~/.ssh/laju-go-deploy

# Set proper permissions
chmod 600 ~/.ssh/laju-go-deploy
```

**Copy the public key:**
```bash
cat ~/.ssh/laju-go-deploy.pub
# Output: ssh-ed25519 AAAA... github-actions@laju-go
```

---

## Step 2: Setup VPS

### Add SSH key to VPS:

```bash
# SSH to your VPS
ssh user@your-vps-ip

# Add GitHub Actions public key to authorized_keys
echo "ssh-ed25519 AAAA... github-actions@laju-go" >> ~/.ssh/authorized_keys

# Or for system-wide deployment (recommended):
sudo mkdir -p /opt/laju-go
sudo chown -R $USER:$USER /opt/laju-go

# Test connection
ssh -i ~/.ssh/laju-go-deploy user@your-vps-ip
```

### Create directories:

```bash
# On VPS
sudo mkdir -p /opt/laju-go
sudo mkdir -p /var/lib/laju
sudo mkdir -p /opt/laju-go/storage/avatars

# Set permissions
sudo chown -R www-data:www-data /var/lib/laju
sudo chown -R www-data:www-data /opt/laju-go/storage
```

### Create systemd service:

```bash
sudo nano /etc/systemd/system/laju-go.service
```

```ini
[Unit]
Description=Laju Go Application
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/laju-go
ExecStart=/opt/laju-go/laju-go
Restart=always
RestartSec=5
EnvironmentFile=/opt/laju-go/.env

[Install]
WantedBy=multi-user.target
```

### Create production `.env`:

```bash
sudo nano /opt/laju-go/.env
```

```bash
# Production Environment
APP_ENV=production
APP_PORT=8080
DB_PATH=/var/lib/laju/app.db

# Generate secure session secret
# Run: openssl rand -base64 32
SESSION_SECRET=<your-32-char-secret>

# Google OAuth
GOOGLE_CLIENT_ID=<your-client-id>
GOOGLE_CLIENT_SECRET=<your-client-secret>
GOOGLE_REDIRECT_URL=https://yourdomain.com/auth/google/callback

# Email/SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=<your-email>
SMTP_PASS=<your-app-password>
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your App Name

# Application URL
APP_URL=https://yourdomain.com
```

### Enable service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable laju-go
sudo systemctl start laju-go

# Check status
sudo systemctl status laju-go
```

---

## Step 3: Configure GitHub Secrets

Go to your GitHub repository → **Settings** → **Secrets and variables** → **Actions** → **New repository secret**

Add these secrets:

| Secret Name | Value | Description |
|-------------|-------|-------------|
| `VPS_HOST` | `192.168.1.100` or `example.com` | Your VPS IP address or domain |
| `VPS_USER` | `root` or `deploy` | SSH username |
| `VPS_SSH_KEY` | Full content of `~/.ssh/laju-go-deploy` | Private key (include BEGIN/END lines) |
| `VPS_PORT` | `22` (optional) | SSH port (default: 22) |

**To copy SSH private key:**
```bash
cat ~/.ssh/laju-go-deploy
# Copy entire output including:
# -----BEGIN OPENSSH PRIVATE KEY-----
# ...
# -----END OPENSSH PRIVATE KEY-----
```

---

## Step 4: Verify Workflows

### Workflow Files

Two workflows are created:

1. **`.github/workflows/deploy.yml`** - Build & Deploy on push to main
2. **`.github/workflows/test.yml`** - Run tests on pull requests

### Test the Setup

```bash
# 1. Commit workflow files
git add .github/workflows/
git commit -m "Setup GitHub Actions CI/CD"

# 2. Push to trigger workflow
git push origin main

# 3. Check Actions tab in GitHub
# Go to: https://github.com/yourusername/laju-go/actions
```

---

## Workflow Details

### Deploy Workflow (`deploy.yml`)

**Triggers:** Push to `main` branch

**Jobs:**
1. **Build** (ubuntu-latest):
   - Setup Go 1.26 + Node.js 18
   - Install dependencies
   - Build frontend (`npm run build`)
   - Build Go binary for Linux (`GOOS=linux GOARCH=amd64`)
   - Upload artifacts

2. **Deploy** (ubuntu-latest):
   - Download artifacts
   - SCP to VPS via SSH
   - Restart systemd service

### Test Workflow (`test.yml`)

**Triggers:** Pull request to `main` branch

**Jobs:**
1. **Test** (ubuntu-latest):
   - Setup Go 1.26 + Node.js 18
   - Install dependencies
   - Run Go tests (`go test ./...`)
   - Run frontend tests (`npm run test:run`)

---

## Deployment Process

### Automatic Deployment

```bash
# Just push to main branch
git add .
git commit -m "Fix login bug"
git push origin main

# GitHub Actions will:
# 1. Build binary (Linux x64)
# 2. Build frontend assets
# 3. Upload to VPS via SCP
# 4. Restart service
# Done! ✅
```

### Manual Deployment (if needed)

```bash
# Build locally
npm run build:all

# Upload to VPS
scp laju-go dist/ templates/ migrations/ public/ user@vps:/opt/laju-go/

# Restart service
ssh user@vps "sudo systemctl restart laju-go"
```

---

## Monitoring

### Check Deployment Status

**GitHub Actions:**
- Visit: `https://github.com/yourusername/laju-go/actions`
- See build/deploy logs in real-time

**VPS Service:**
```bash
# Check service status
ssh user@vps "sudo systemctl status laju-go"

# View logs
ssh user@vps "sudo journalctl -u laju-go -f"

# View recent errors
ssh user@vps "sudo journalctl -u laju-go -p err -n 50"
```

### Health Check

Add health check endpoint to monitor application:

```go
// routes/web.go
app.Get("/health", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "healthy",
        "timestamp": time.Now().UTC(),
        "version": "1.0.0",
    })
})
```

Test:
```bash
curl https://yourdomain.com/health
```

---

## Troubleshooting

### Deployment Fails

**Check GitHub Actions logs:**
```
GitHub → Repository → Actions → Latest workflow run → Check logs
```

**Common issues:**

| Issue | Solution |
|-------|----------|
| SSH connection failed | Verify VPS_HOST and SSH key |
| Permission denied | Check authorized_keys on VPS |
| Service won't start | Check systemd logs: `journalctl -u laju-go` |
| Binary not found | Verify SCP target path |

### Service Won't Start

```bash
# SSH to VPS
ssh user@vps

# Check service status
sudo systemctl status laju-go

# View logs
sudo journalctl -u laju-go -n 50

# Common fixes:
# 1. Check .env file exists
ls -la /opt/laju-go/.env

# 2. Check binary permissions
ls -la /opt/laju-go/laju-go
sudo chmod +x /opt/laju-go/laju-go

# 3. Check directory permissions
sudo chown -R www-data:www-data /opt/laju-go
sudo chown -R www-data:www-data /var/lib/laju

# 4. Test binary manually
sudo -u www-data /opt/laju-go/laju-go
```

### Database Errors

```bash
# Check database directory
ls -la /var/lib/laju/

# Fix permissions
sudo chown -R www-data:www-data /var/lib/laju

# Check SQLite files
sqlite3 /var/lib/laju/app.db ".tables"
```

---

## Security Best Practices

### 1. Use Dedicated Deploy User

```bash
# On VPS, create dedicated user
sudo useradd -r -s /bin/false laju-deploy

# Add SSH key for this user
sudo mkdir -p /home/laju-deploy/.ssh
sudo nano /home/laju-deploy/.ssh/authorized_keys
# Paste GitHub Actions public key

sudo chown -R laju-deploy:laju-deploy /home/laju-deploy/.ssh
sudo chmod 700 /home/laju-deploy/.ssh
sudo chmod 600 /home/laju-deploy/.ssh/authorized_keys
```

### 2. Limit SSH Key Permissions

In `authorized_keys` on VPS, add restrictions:

```
command="/usr/bin/systemctl restart laju-go",no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty ssh-ed25519 AAAA... github-actions@laju-go
```

### 3. Use Environment-Specific Secrets

```bash
# Production .env (on VPS only)
SESSION_SECRET=<production-secret>

# Development .env (local only)
SESSION_SECRET=<dev-secret>

# NEVER commit .env files!
```

### 4. Enable GitHub 2FA

Always enable Two-Factor Authentication for GitHub accounts with deployment access.

---

## Rollback Strategy

### Quick Rollback

```bash
# On VPS
cd /opt/laju-go

# List previous deployments (if you keep backups)
ls -la backups/

# Restore previous binary
sudo systemctl stop laju-go
cp backups/laju-go-previous ./laju-go
sudo systemctl start laju-go
```

### Git-Based Rollback

```bash
# Revert to previous commit
git revert HEAD
git push origin main

# GitHub Actions will deploy the reverted version
```

---

## Next Steps

- [Production Deployment Guide](production.md) - Bare-metal deployment details
- [Docker Deployment](docker.md) - Containerized deployment
- [Monitoring Guide](monitoring.md) - Application monitoring setup
- [Backup Strategy](backup.md) - Database backup automation
