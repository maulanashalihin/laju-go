# Laju Go - Project Context

## Project Overview

**Laju Go** is a high-performance SaaS boilerplate built with:
- **Backend**: Go Fiber (fasthttp-based web framework)
- **Frontend**: Svelte 5 with Inertia.js for SPA experience
- **Database**: SQLite with production optimizations (WAL mode, connection pooling)
- **Architecture**: Clean layered architecture (Routes → Middleware → Handler → Service → Repository → Database)

### Key Features
- User authentication (email/password + Google OAuth)
- Password reset via email (SMTP mailer service)
- Role-based access control (Admin/User roles)
- Database-backed sessions (persistent, SQLite)
- CSRF protection middleware
- Rate limiting middleware (auth, password reset)
- File upload support (avatars)
- Database migrations with Goose
- Hot module replacement (HMR) in development
- Docker-ready for production deployment

---

## Tech Stack

| Layer | Technology | Version |
|-------|------------|---------|
| Language | Go | 1.26+ |
| Web Framework | Go Fiber v2 | 2.52.5 |
| Database | SQLite3 | 1.14.22 |
| Query Builder | Squirrel | 1.5.4 |
| Migrations | Goose | 3.20.0 |
| Frontend | Svelte | 5.55.0 |
| Build Tool | Vite | 8.0.3 |
| Styling | Tailwind CSS | 4.2.2 |
| SPA Bridge | Inertia.js | 3.0.0 |
| Icons | Lucide Svelte | 1.0.1 |
| OAuth | golang.org/x/oauth2 | 0.18.0 |
| Session | Database-backed (SQLite) | - |
| Email | SMTP (MailerService) | - |
| Testing (Frontend) | Vitest + Happy-DOM | 4.1.2 |

---

## Project Structure

```
laju-go/
├── main.go                  # Entry point (server bootstrap, DB setup, migrations)
├── go.mod                   # Go dependencies
├── package.json             # Node.js dependencies
├── vite.config.js           # Vite configuration
│
├── app/                     # Go backend code
│   ├── config/              # Environment configuration (config.go)
│   ├── handlers/            # HTTP controllers
│   │   ├── app.go           # App dashboard handler
│   │   ├── auth.go          # Authentication handler
│   │   ├── password-reset.go # Password reset handler
│   │   ├── public.go        # Public pages handler
│   │   └── upload.go        # File upload handler
│   ├── middlewares/         # Auth guards (AuthRequired, AdminRequired, Guest)
│   ├── models/              # Data structures (User, DTOs)
│   ├── repositories/        # Database queries (Squirrel SQL builder)
│   │   ├── session.repository.go # Session repository
│   │   └── user.repository.go    # User repository
│   ├── services/            # Business logic
│   │   ├── asset.go         # Asset service (Vite manifest, hashed filenames)
│   │   ├── auth.go          # Authentication service
│   │   ├── inertia.go       # Inertia.js response builder
│   │   ├── mailer.go        # SMTP email service
│   │   └── user.go          # User management service
│   └── session/             # Session infrastructure (cookie encoding/decoding)
│
├── routes/
│   └── web.go               # Route definitions & middleware setup
│
├── frontend/                # Svelte 5 frontend
│   ├── src/
│   │   ├── components/      # Reusable UI components (Button, Input)
│   │   ├── layouts/         # Page layouts
│   │   ├── lib/             # Utility modules
│   │   ├── pages/           # Page components
│   │   │   ├── admin/       # Admin pages
│   │   │   ├── app/         # Protected app pages
│   │   │   └── auth/        # Authentication pages (login, register)
│   │   ├── main.ts          # Inertia.js entry point
│   │   └── app.css          # Global styles (Tailwind)
│
├── migrations/              # Database migrations (Goose)
│   ├── 0001_create_users_table.sql
│   └── 0002_create_sessions_table.sql
├── templates/               # HTML templates
│   ├── index.html           # Index page template
│   └── inertia.html         # Base Inertia.js template
├── data/                    # SQLite database (gitignored)
├── storage/                 # User uploads (gitignored)
├── dist/                    # Production build assets (Vite output)
└── public/                  # Static assets
```

### Layer Separation Note

The `app/session/` folder is **separate** from `app/services/`:
- **`session/`** = Infrastructure layer (generic cookie storage, reusable)
- **`services/`** = Business logic layer (domain-specific rules)

Services use session infrastructure but session knows nothing about business domain.

---

## Building and Running

### Prerequisites
- Go 1.26+
- Node.js 18+
- SQLite3

### Development Setup (First Time)

```bash
# 1. Install Go dependencies
go mod download

# 2. Install Node.js dependencies
npm install

# 3. Copy environment file
cp .env.example .env

# 4. Edit .env with your configuration
#    At minimum: SESSION_SECRET, APP_PORT, DB_PATH
```

### Development Workflow (Every Session)

```bash
# If you encounter errors, reset dependencies:
go mod download && npm install
cp .env.example .env  # Only if .env is missing
```

### Pre-Commit Checklist

Before committing changes, always verify:

```bash
# 1. Go build passes (no compile errors)
npm run build:go

# 2. Go tests pass
go test ./...

# 3. Frontend build passes
npm run build

# 4. Frontend tests pass (optional)
npm run test:run
```

For production deployment:

```bash
# Build everything (frontend + Go binary)
npm run build:all
```

### Available Commands

| Command | Description |
|---------|-------------|
| `npm run dev` | Start Vite dev server (frontend HMR) |
| `npm run dev:go` | Start Go server with Air (hot reload) |
| `npm run dev:all` | Run both Vite and Air concurrently |
| `npm run build` | Build frontend only (Vite → dist/) |
| `npm run build:go` | Build Go binary only |
| `npm run build:all` | Build frontend + Go binary |
| `npm run serve` | Run production binary (`./laju-go`) |
| `npm run db:refresh` | Remove database files (app.db, app.db-shm, app.db-wal) |
| `npm run db:migrate` | Run migrations via main.go |
| `go run main.go` | Run Go server directly |
| `go test ./...` | Run Go tests |
| `npm run test:run` | Run frontend tests (Vitest) |
| `npm run test:ui` | Run frontend tests with UI |

### Environment Configuration

Required `.env` variables:
```bash
# Server
APP_PORT=8080
APP_ENV=development

# Database
DB_PATH=./data/app.db

# Session
SESSION_SECRET=your-secret-key-change-this-in-production

# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# Frontend (development)
FRONTEND_URL=http://localhost:5173

# Email/SMTP (for password reset)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-smtp-username
SMTP_PASS=your-smtp-password
FROM_EMAIL=noreply@example.com
FROM_NAME=Laju
```

---

## Development Conventions

### Code Organization
- **Handlers**: Parse requests, call services, return responses
- **Services**: Business logic, authentication flows, user management, email sending
- **Repositories**: Database operations using Squirrel query builder
- **Models**: Domain models and DTOs
- **Middlewares**: Request interception (auth checks, rate limiting, guest checks)

### Inertia.js Pattern
- Initial load: Server renders HTML via `inertia.html` template
- Subsequent navigation: Inertia XHR requests with `X-Inertia: true` header
- Server returns JSON with component name and props
- Svelte dynamically loads components client-side

### Making Inertia Requests (Frontend)
```svelte
<script>
  import { router } from '@inertiajs/svelte'

  function handleSubmit() {
    router.post('/login', {
      email: formData.email,
      password: formData.password
    })
  }
</script>
```

### Protected Routes (Backend)
```go
// Requires authentication
app.Get("/app", middlewares.AuthRequired(store), AppHandler.Dashboard)

// Admin only
app.Get("/admin", middlewares.AdminRequired(store), AdminHandler.Dashboard)

// Guest only (redirect if authenticated)
app.Get("/login", middlewares.Guest(store), AuthHandler.ShowLoginForm)
```

### Database Migrations
```bash
# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run all migrations
goose -dir migrations sqlite3 ./data/app.db up

# Check status
goose -dir migrations sqlite3 ./data/app.db status

# Rollback last migration
goose -dir migrations sqlite3 ./data/app.db down

# Create new migration
goose -dir migrations create migration_name
```

### SQLite Production Optimizations
Applied automatically in `main.go`:
- `journal_mode=WAL` - Write-Ahead Logging
- `synchronous=NORMAL` - Balance speed/durability
- `cache_size=-16000` - 16MB cache (optimized for 1-2GB RAM)
- `mmap_size=268435456` - 256MB memory-mapped I/O
- `temp_store=MEMORY` - Memory temp tables
- `busy_timeout=5000` - 5s lock wait timeout
- `wal_autocheckpoint=1000` - WAL autocheckpoint pages
- `max_open_conns=15` - Connection pool max size
- `max_idle_conns=5` - Connection pool idle size

---

## Key Routes

### Public
- `GET /` - Home page
- `GET /about` - About page

### Authentication
- `GET /login` - Login form (Guest only)
- `POST /login` - User login (Guest only, rate-limited)
- `GET /register` - Registration form (Guest only)
- `POST /register` - User registration (Guest only, rate-limited)
- `GET /auth/google` - Google OAuth
- `GET /auth/google/callback` - OAuth callback
- `POST /logout` - Logout (requires auth)
- `GET /api/me` - Current user API (requires auth)

### Password Reset
- `GET /forgot-password` - Request reset form
- `POST /forgot-password` - Send reset email (rate-limited)
- `GET /reset-password/:token` - Reset password form
- `POST /reset-password/:token` - Complete password reset

### Protected
- `GET /app` - Dashboard (requires auth, CSRF protected)
- `GET /app/profile` - Profile page (requires auth, CSRF protected)
- `PUT /app/profile` - Update profile (requires auth, CSRF protected)
- `PUT /app/profile/password` - Update password (requires auth, CSRF protected)
- `POST /upload` - File upload (requires auth, CSRF protected)

### Admin
- `GET /admin` - Admin dashboard

---

## Testing Practices

- **Backend**: Go tests in `*_test.go` files, run with `go test ./...`
- **Frontend**: Vitest with Happy-DOM, run with `npm run test:run` or `npm run test:ui`
- **Test isolation**: Use separate test database if needed

---

## Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| Port 8080 in use | `lsof -ti:8080 | xargs kill -9` |
| Database locked | Remove `data/app.db-shm` and `data/app.db-wal` |
| Vite port detection fails | Delete `.vite-port` and restart Vite |
| Migration errors | Check status with `goose status`, reset if needed |
| No .env file | Copy `.env.example` to `.env` and configure |
| Google OAuth fails | Verify redirect URL matches exactly in Google Cloud Console |
| Email not sending | Use Gmail App Password, not regular password |

---

## Documentation

Full documentation is available in the `docs/` directory:

| Section | Path |
|---------|------|
| Getting Started | `docs/getting-started/` |
| Architecture Guide | `docs/guide/architecture.md` |
| Routing & Handlers | `docs/guide/routing.md`, `docs/guide/handlers.md` |
| Database | `docs/guide/database.md` |
| Frontend | `docs/guide/frontend.md` |
| File Upload | `docs/guide/file-upload.md` |
| Email | `docs/guide/email.md` |
| Data Protection | `docs/guide/data-protection.md` |
| Development | `docs/deployment/development.md` |
| Production | `docs/deployment/production.md` |
| Docker | `docs/deployment/docker.md` |
| GitHub Actions CI/CD | `docs/deployment/github-actions.md` |
| SQLite Config | `docs/deployment/sqlite-configuration.md` |
| Optimization | `docs/deployment/optimization.md` |
| API Reference | `docs/reference/api-reference.md` |
| Project Structure | `docs/reference/project-structure.md` |
| Environment | `docs/reference/environment.md` |
| Troubleshooting | `docs/reference/troubleshooting.md` |

---

## Deployment

### Option 1: CI/CD with GitHub Actions (Recommended)

Automated build and deployment on every push to `main`:

```bash
# Just push to main
git add .
git commit -m "Fix login bug"
git push origin main

# GitHub Actions will:
# 1. Build binary (Linux x64) in clean environment
# 2. Build frontend assets
# 3. Upload to VPS via SCP
# 4. Restart systemd service
# No Go/Node.js needed on production server!
```

**Setup:** See [GitHub Actions CI/CD Guide](docs/deployment/github-actions.md)

### Option 2: Manual Build & Upload

Build locally and upload to server:

```bash
# Build for Linux (from macOS)
GOOS=linux GOARCH=amd64 go build -o laju-go .
npm run build

# Upload to VPS
scp laju-go dist/ templates/ migrations/ public/ user@vps:/opt/laju-go/

# Restart service
ssh user@vps "sudo systemctl restart laju-go"
```

### Option 3: Build on Server

Pull source code and build on VPS:

```bash
# SSH to server
ssh user@vps

# Pull code
cd /opt/laju-go && git pull

# Build (requires Go + Node.js on server)
npm install && npm run build:all

# Restart service
sudo systemctl restart laju-go
```

---

## Deployment Notes

### What to Include When Shipping

| Include | Purpose |
|---------|---------|
| `laju-go` binary | Compiled Go executable |
| `templates/` | HTML templates (inertia.html, index.html) |
| `dist/` | Compiled frontend assets (JS, CSS from Vite) |
| `migrations/` | Database migrations |
| `public/` | Static assets |
| `.env` | Production configuration |

| Exclude | Reason |
|---------|--------|
| `data/` | Database created at runtime |
| `storage/` | User uploads created at runtime |
| Source code | Already compiled into binary |
| `node_modules/` | Not needed for runtime |
| `.git/` | Not needed for runtime |
| `.github/` | CI/CD workflows (deployment only) |

### Build Commands

| Command | Purpose |
|---------|---------|
| `npm run build` | Frontend only (Vite → dist/) |
| `npm run build:go` | Go binary only |
| `npm run build:all` | Frontend + Go binary (for production) |
| `GOOS=linux GOARCH=amd64 go build` | Cross-compile for Linux VPS |
