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
- One-click deployment via SSH

---

## Tech Stack

| Layer | Technology | Version |
|-------|------------|---------|
| Language | Go | 1.26+ |
| Web Framework | Go Fiber v2 | 2.52.5 |
| Database | SQLite (modernc.org) | 1.39.1 |
| Query Builder | Squirrel | 1.5.4 |
| Migrations | Goose | 3.20.0 |
| Frontend | Svelte | 5.55.0 |
| Build Tool | Vite | 8.0.3 |
| Styling | Tailwind CSS | 4.2.2 |
| SPA Bridge | Inertia.js | 3.0.0 |

---

## Project Structure

```
laju-go/
├── main.go                  # Entry point
├── go.mod                   # Go dependencies
├── package.json             # Node.js dependencies
├── .deploy.example          # Deployment config template
│
├── app/                     # Backend code
│   ├── config/              # Environment configuration
│   ├── handlers/            # HTTP controllers
│   ├── middlewares/         # Auth guards, rate limiting
│   ├── models/              # Data structures
│   ├── repositories/        # Database queries (Squirrel)
│   ├── services/            # Business logic
│   └── session/             # Session infrastructure
│
├── routes/
│   └── web.go               # Route definitions & static files
│
├── frontend/                # Svelte 5 + Inertia.js
│   └── src/
│       ├── components/      # Reusable UI components
│       ├── pages/           # Page components
│       └── main.ts          # Entry point
│
├── migrations/              # Database migrations (Goose)
├── templates/               # HTML templates
├── scripts/                 # Deployment scripts
│   ├── deploy.sh            # Main deployment script
│   ├── first-deploy.sh      # First deploy setup
│   └── update-deploy.sh     # Update deployment
├── data/                    # SQLite database (gitignored)
├── storage/                 # User uploads (gitignored)
├── dist/                    # Production build assets
└── public/                  # Static assets (served at /public)
```

---

## Building and Running

### Prerequisites
- Go 1.26+
- Node.js 18+

### Development Setup

```bash
go mod download && npm install
cp .env.example .env
npm run dev:all
```

### Available Commands

| Command | Description |
|---------|-------------|
| `npm run dev` | Vite dev server (HMR) |
| `npm run dev:go` | Go server with Air (hot reload) |
| `npm run dev:all` | Both Vite + Air |
| `npm run build` | Frontend only |
| `npm run build:go` | Go binary only |
| `npm run build:all` | Frontend + Go binary |
| `npm run build:linux` | Cross-compile for Linux |
| `npm run serve` | Run production binary |
| `npm run deploy` | One-click deployment |
| `npm run db:migrate` | Run migrations |
| `npm run db:migrate:status` | Check migration status |
| `npm run db:migrate:down` | Rollback last migration |
| `npm run db:migrate:create` | Create new migration |
| `npm run db:refresh` | Reset database |

---

## Development Conventions

### Code Organization
- **Handlers**: Parse requests, call services, return responses
- **Services**: Business logic, auth flows, email sending
- **Repositories**: Database operations (Squirrel SQL builder)
- **Middlewares**: Auth checks, rate limiting, CSRF

### Inertia.js Pattern
- Initial load: Server renders HTML via `inertia.html`
- Subsequent: XHR with `X-Inertia: true` header → JSON response
- Frontend dynamically loads components

### Protected Routes

```go
// Requires authentication
app.Get("/app", middlewares.AuthRequired(store), AppHandler.Dashboard)

// Admin only
app.Get("/admin", middlewares.AdminRequired(store), AdminHandler.Dashboard)

// Guest only
app.Get("/login", middlewares.Guest(store), AuthHandler.ShowLoginForm)
```

### Static File Serving

Configured in `routes/web.go`:
- `/dist` → Built frontend assets
- `/public` → Public assets (images, etc.)
- `/storage` → User uploads (avatars)

---

## Deployment

### One-Click Deployment (Recommended)

```bash
# Configure deployment
cp .deploy.example .deploy
nano .deploy  # Set SERVER_USER, SERVER_HOST, SERVER_PATH, REPO_URL

# Deploy
npm run deploy
```

**What happens:**
1. Builds frontend + Go binary (linux/amd64)
2. Commits and pushes to GitHub
3. Auto-detects first deploy vs update
4. **First deploy**: Creates directory, clones repo, sets up systemd, configures .env
5. **Update deploy**: Pulls changes, rebuilds, restarts service

### Manual Deployment

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o laju-go .
npm run build

# Upload to server
scp laju-go dist/ templates/ migrations/ public/ user@vps:/opt/laju-go/

# Restart service
ssh user@vps "systemctl restart laju-go"
```

---

## Environment Configuration

### Required `.env` variables

```bash
# Server
APP_PORT=8080
APP_ENV=development

# Database
DB_PATH=./data/app.db

# Session
SESSION_SECRET=<generate-random-secret>

# Google OAuth (optional)
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# SMTP (for password reset)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
FROM_EMAIL=noreply@example.com
```

---

## Database

### Migrations

```bash
# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir migrations sqlite ./data/app.db up

# Check status
goose -dir migrations sqlite ./data/app.db status

# Rollback
goose -dir migrations sqlite ./data/app.db down
```

### SQLite Optimizations (Auto-applied)

- `journal_mode=WAL` - Write-Ahead Logging
- `synchronous=NORMAL` - Balance speed/durability
- `cache_size=-16000` - 16MB cache
- `mmap_size=268435456` - 256MB memory-mapped I/O
- `busy_timeout=5000` - 5s lock wait timeout
- `max_open_conns=15` - Connection pool

---

## Common Issues

| Issue | Solution |
|-------|----------|
| Port 8080 in use | `lsof -ti:8080 | xargs kill -9` |
| Database locked | Remove `data/app.db-shm` and `data/app.db-wal` |
| Migration errors | `goose -dir migrations sqlite ./data/app.db status` |
| Google OAuth fails | Check redirect URL in Google Cloud Console |
| Email not sending | Use Gmail App Password, not regular password |

---

## Documentation

See `docs/` directory for detailed guides:

- [One-Click Deployment](docs/deployment/one-click-deployment.md)
- [Production Deployment](docs/deployment/production.md)
- [Architecture Guide](docs/guide/architecture.md)
- [API Reference](docs/reference/api-reference.md)
- [Troubleshooting](docs/reference/troubleshooting.md)
