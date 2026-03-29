# One-Click Deployment

Laju Go includes a one-click deployment script that automates the entire deployment process to your production server via SSH.

## Overview

The deployment system consists of three scripts:

| Script | Purpose |
|--------|---------|
| `scripts/deploy.sh` | Main deployment script (auto-detects first vs update deploy) |
| `scripts/first-deploy.sh` | Initial deployment setup (creates directories, systemd service) |
| `scripts/update-deploy.sh` | Update existing deployment (pulls changes, restarts service) |

## Prerequisites

### Server Requirements

- **OS**: Ubuntu 20.04+ or Debian 11+
- **RAM**: Minimum 512MB (1GB recommended)
- **Storage**: 10GB+ (depends on database size)
- **CPU**: 1 core minimum (2+ recommended)
- **SSH Access**: Root or sudo user with SSH key authentication

### Local Requirements

- Go 1.26+ installed
- Node.js 18+ installed
- Git installed
- SSH access to your server

## Step 1: Configure Deployment

### Copy Deployment Configuration

```bash
cp .deploy.example .deploy
```

### Edit `.deploy` File

```bash
nano .deploy
```

```bash
# Deployment Configuration

# SSH Credentials
SERVER_USER=root
SERVER_HOST=your.server.com

# Remote server path where the application will be deployed
SERVER_PATH=/opt/your-app

# Git repository URL (for cloning on first deploy)
REPO_URL=https://github.com/yourusername/your-repo.git

# Systemd service name
SERVICE_NAME=your-app
```

### Configuration Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `SERVER_USER` | SSH username | `root` or `deploy` |
| `SERVER_HOST` | Server IP or domain | `192.168.1.100` or `example.com` |
| `SERVER_PATH` | Remote deployment path | `/opt/my-app` |
| `REPO_URL` | Git repository URL | `https://github.com/user/repo.git` |
| `SERVICE_NAME` | Systemd service name | `my-app` |

## Step 2: Setup SSH Access

### Generate SSH Key (if needed)

```bash
ssh-keygen -t ed25519 -C "your-email@example.com"
```

### Copy SSH Key to Server

```bash
ssh-copy-id root@your.server.com
```

### Test SSH Connection

```bash
ssh root@your.server.com
```

## Step 3: Deploy

### Run Deployment Script

```bash
npm run deploy
```

Or directly:

```bash
./scripts/deploy.sh
```

### What Happens During Deployment

**First Deploy:**
1. ✅ Creates remote directory
2. ✅ Clones Git repository
3. ✅ Sets up `.env` file (with interactive prompts)
4. ✅ Creates systemd service
5. ✅ Starts the application

**Update Deploy:**
1. ✅ Pulls latest changes
2. ✅ Builds frontend
3. ✅ Builds Go binary
4. ✅ Restarts systemd service

## First Deploy Interactive Prompts

During first deployment, you'll be prompted for:

### 1. Application Port

```
Application Port (default: 8080):
```

Enter the port your application will run on (default: 8080).

### 2. Application URL

```
Application URL (e.g., https://yourdomain.com):
```

Enter your production domain URL (e.g., `https://example.com`).

### Auto-Generated Configuration

The following are automatically configured:

- ✅ `SESSION_SECRET` - Random 64-character hex string
- ✅ `APP_ENV` - Set to `production`
- ✅ `APP_PORT` - From your input
- ✅ `APP_URL` - From your input

## Post-Deployment Configuration

### Required: Configure OAuth & SMTP

After deployment, you need to configure Google OAuth and SMTP for full functionality.

#### Edit .env on Server

```bash
ssh root@your.server.com 'nano /opt/your-app/.env'
```

#### Google OAuth Configuration

Get credentials from [Google Cloud Console](https://console.cloud.google.com):

```bash
# Google OAuth
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-secret-key
GOOGLE_REDIRECT_URL=https://yourdomain.com/auth/google/callback
```

#### SMTP Configuration (for password reset)

```bash
# SMTP (Gmail example)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your App Name
```

**Note**: For Gmail, use an [App Password](https://support.google.com/accounts/answer/185833), not your regular password.

## Deployment Scripts Explained

### deploy.sh (Main Script)

Auto-detects whether it's a first deploy or update:

```bash
#!/bin/bash
# 1. Loads .deploy configuration
# 2. Tests SSH connection
# 3. Builds assets locally
# 4. Commits and pushes changes
# 5. Calls first-deploy.sh or update-deploy.sh
```

### first-deploy.sh (Initial Setup)

Sets up the application from scratch:

```bash
#!/bin/bash
# 1. Creates remote directory
# 2. Clones Git repository
# 3. Creates .env from .env.example (with auto-configuration)
# 4. Creates and starts systemd service
```

### update-deploy.sh (Updates)

Updates existing deployment:

```bash
#!/bin/bash
# 1. Pulls latest changes
# 2. Builds frontend
# 3. Builds Go binary
# 4. Restarts systemd service
```

## Useful Commands

### View Deployment Status

```bash
ssh root@your.server.com 'systemctl status your-app'
```

### View Logs (Real-time)

```bash
ssh root@your.server.com 'journalctl -u your-app -f'
```

### View Recent Errors

```bash
ssh root@your.server.com 'journalctl -u your-app -p err -n 50'
```

### Restart Service

```bash
ssh root@your.server.com 'systemctl restart your-app'
```

### Stop Service

```bash
ssh root@your.server.com 'systemctl stop your-app'
```

### Start Service

```bash
ssh root@your.server.com 'systemctl start your-app'
```

## Troubleshooting

### Service Won't Start

**Check logs:**

```bash
ssh root@your.server.com 'journalctl -u your-app -n 50'
```

**Common issues:**
- Missing `.env` file
- Invalid `SESSION_SECRET`
- Database path not writable
- Port already in use

### Database Locked

**Solution:**

```bash
# Stop service
ssh root@your.server.com 'systemctl stop your-app'

# Remove WAL files
ssh root@your.server.com 'rm /opt/your-app/data/app.db-shm /opt/your-app/data/app.db-wal'

# Start service
ssh root@your.server.com 'systemctl start your-app'
```

### Deployment Script Fails

**Check SSH connection:**

```bash
ssh -v root@your.server.com
```

**Verify .deploy configuration:**

```bash
cat .deploy
```

**Run script with debug mode:**

```bash
bash -x scripts/deploy.sh
```

### Build Fails on Server

**Check Go version on server:**

```bash
ssh root@your.server.com 'go version'
```

**Check Node.js version:**

```bash
ssh root@your.server.com 'node --version'
```

**Manual rebuild:**

```bash
ssh root@your.server.com
cd /opt/your-app
npm install
npm run build
go build -o laju-go .
systemctl restart your-app
```

## Manual Deployment (Alternative)

If the automated script doesn't work, you can deploy manually:

### 1. Connect to Server

```bash
ssh root@your.server.com
```

### 2. Clone Repository

```bash
mkdir -p /opt/your-app
cd /opt/your-app
git clone https://github.com/yourusername/your-repo.git .
```

### 3. Install Dependencies

```bash
go mod download
npm install
```

### 4. Build Application

```bash
npm run build
go build -o laju-go .
```

### 5. Configure Environment

```bash
cp .env.example .env
nano .env
```

### 6. Create Systemd Service

```bash
nano /etc/systemd/system/your-app.service
```

See [Production Deployment](production.md) for service configuration.

### 7. Start Service

```bash
systemctl daemon-reload
systemctl enable your-app
systemctl start your-app
systemctl status your-app
```

## Next Steps

- [Production Deployment](production.md) - Detailed production setup guide
- [Docker Deployment](docker.md) - Containerized deployment
- [GitHub Actions CI/CD](github-actions.md) - Automated CI/CD pipeline
- [Optimization Guide](optimization.md) - Performance optimization
