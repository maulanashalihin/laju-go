# Installation

This guide walks you through setting up Laju Go on your local machine.

## Quick Start (Recommended)

The fastest way to get started is using the `create-laju-go` CLI:

```bash
# Create new project
npx create-laju-go my-app

# Navigate to project
cd my-app

# Start development
npm run dev:all
```

The CLI automatically:
- Checks for Go and Git installation
- Lets you choose a package manager (npm, yarn, bun)
- Clones the template from GitHub
- Installs all dependencies
- Sets up environment configuration

Visit `http://localhost:8080` to see your application running.

---

## Manual Installation

Follow these steps if you prefer to install manually or are contributing to the framework.

Before you begin, ensure you have the following installed:

| Software | Version | Purpose | Download |
|----------|---------|---------|----------|
| **Go** | 1.26+ | Backend runtime | [go.dev](https://go.dev/dl/) |
| **Node.js** | 18+ | Frontend build tools | [nodejs.org](https://nodejs.org/) |
| **Git** | Latest | Version control | [git-scm.com](https://git-scm.com/) |
| **SQLite3** | Latest | Database | Pre-installed on macOS/Linux |

### Verify Installations

```bash
# Check Go version
go version
# Output: go version go1.26.0 darwin/amd64

# Check Node.js version
node --version
# Output: v20.x.x

# Check npm version
npm --version
# Output: 10.x.x

# Check SQLite version
sqlite3 --version
# Output: 3.x.x
```

## Step 1: Clone the Repository

```bash
git clone https://github.com/maulanashalihin/laju-go.git
cd laju-go
```

## Step 2: Install Go Dependencies

```bash
go mod download
```

This command reads `go.mod` and downloads all required Go packages to your local cache.

### Verify Go Dependencies

```bash
go list -m all
```

You should see packages like:
- `github.com/gofiber/fiber/v2`
- `github.com/mattn/go-sqlite3`
- `github.com/Masterminds/squirrel`
- `github.com/pressly/goose/v3`

## Step 3: Install Node.js Dependencies

```bash
npm install
```

This installs all frontend dependencies including:
- Svelte 5
- Vite
- Inertia.js
- Tailwind CSS
- Development tools

### Verify Node Dependencies

```bash
npm list --depth=0
```

## Step 4: Configure Environment

Copy the example environment file:

```bash
cp .env.example .env
```

### Minimum Required Configuration

Edit `.env` and set at least these variables:

```bash
# Application
APP_ENV=development
APP_PORT=8080

# Session (generate a random 32+ character string)
SESSION_SECRET=your-32-character-secret-key-change-in-production

# Database
DB_PATH=data/app.db
```

> 💡 **Tip**: Use a password generator or run this command to create a secure secret:
> ```bash
> openssl rand -base64 32
> ```

## Step 5: Initialize the Database

Database migrations run automatically when you start the server. However, you can manually run them:

```bash
# Install goose (Go migration tool)
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir migrations sqlite3 data/app.db up
```

### Verify Database Setup

```bash
# Check if database file exists
ls -la data/app.db

# Open SQLite CLI
sqlite3 data/app.db ".tables"
# Output: users  sessions  goose_db_version
```

## Step 6: Start Development Servers

### Option A: Run Everything Together (Recommended)

```bash
npm run dev:all
```

This starts both Vite (frontend) and Air (Go hot reload) concurrently.

### Option B: Run Servers Separately

**Terminal 1** - Vite dev server:
```bash
npm run dev
```

**Terminal 2** - Go server with hot reload:
```bash
air
# Or via npm
npm run dev:go
```

### Option C: Manual Run

**Terminal 1** - Vite:
```bash
npm run dev
```

**Terminal 2** - Go (requires manual restart after changes):
```bash
go run main.go
```

## Step 7: Verify Installation

Open your browser and visit `http://localhost:8080`.

You should see:
- ✅ Landing page with navigation
- ✅ Login and Register links
- ✅ No console errors in browser DevTools

### Test Authentication

1. Click **Register** and create an account
2. Verify you're redirected to the dashboard
3. Check that the navigation shows your username
4. Try logging out and back in

### Test Hot Reload

**Frontend HMR:**
1. Edit `frontend/src/pages/app/Dashboard.svelte`
2. Add `<h1>Test</h1>` to the template
3. Save the file
4. Browser should update instantly (no refresh)

**Backend Hot Reload:**
1. Edit `app/handlers/public.go`
2. Add a log statement: `fmt.Println("Test reload")`
3. Save the file
4. Terminal should show Air rebuilding (~1-2 seconds)

## Troubleshooting

### Port Already in Use

**Error**: `listen tcp :8080: bind: address already in use`

**Solution**: Kill the process using port 8080

```bash
# macOS/Linux
lsof -ti:8080 | xargs kill -9

# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

### SQLite Database Locked

**Error**: `database is locked`

**Solution**: Remove WAL files and restart

```bash
rm data/app.db-shm data/app.db-wal
```

### Vite Port Detection Fails

**Error**: Go server can't connect to Vite dev server

**Solution**: Remove the port cache file and restart

```bash
rm .vite-port
npm run dev
```

### Go Module Errors

**Error**: `missing go.sum entry for module`

**Solution**: Clean and re-download dependencies

```bash
go clean -modcache
go mod download
```

### Node Module Issues

**Error**: `Cannot find module` or version conflicts

**Solution**: Clean reinstall

```bash
rm -rf node_modules package-lock.json
npm install
```

### Air Not Installing

**Error**: `command not found: air`

**Solution**: Ensure Go bin directory is in PATH

```bash
# Add to ~/.zshrc or ~/.bashrc
export PATH=$PATH:$(go env GOPATH)/bin

# Reload shell
source ~/.zshrc

# Reinstall Air
go install github.com/air-verse/air@latest
```

## Next Steps

- [Configuration](configuration.md) - Complete environment variable reference
- [Development Workflow](../deployment/development.md) - Hot reload, scripts, and best practices
- [Architecture Guide](../guide/architecture.md) - Understanding the codebase structure

## Quick Reference

```bash
# Development commands
npm run dev          # Start Vite dev server
npm run dev:go       # Start Go server with Air
npm run dev:all      # Run both servers

# Build commands
npm run build        # Build for production
npm run serve        # Run production binary

# Database commands
goose -dir migrations sqlite3 data/app.db up      # Run migrations
goose -dir migrations sqlite3 data/app.db status  # Check status
goose -dir migrations sqlite3 data/app.db down    # Rollback

# Test commands
go test ./...        # Run Go tests
npm run test:run     # Run frontend tests
```
