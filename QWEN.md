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

| Layer | Technology |
|-------|------------|
| Language | Go 1.26+ |
| Web Framework | Go Fiber v2 |
| Database | SQLite3 |
| Query Builder | Squirrel |
| Migrations | Goose |
| Frontend | Svelte 5 |
| Build Tool | Vite |
| Styling | Tailwind CSS |
| SPA Bridge | Inertia.js |
| OAuth | golang.org/x/oauth2 |
| Session | Database-backed (SQLite) |
| Email | SMTP (MailerService) |

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
│   ├── config/              # Environment configuration
│   ├── handlers/            # HTTP controllers
│   ├── middlewares/         # Auth guards (AuthRequired, AdminRequired, Guest)
│   ├── models/              # Data structures (User, DTOs)
│   ├── repositories/        # Database queries (Squirrel SQL builder)
│   ├── services/            # Business logic (Auth, User, Inertia, Asset)
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
│   │   ├── pages/           # Page components (Auth/, App/, Admin/)
│   │   ├── main.ts          # Inertia.js entry point
│   │   └── app.css          # Global styles (Tailwind)
│
├── migrations/              # Database migrations (Goose)
├── templates/               # HTML templates (inertia.html base)
├── data/                    # SQLite database (gitignored)
├── storage/                 # User uploads (gitignored)
├── dist/                    # Production build assets
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

### Development Setup

```bash
# 1. Install Go dependencies
go mod download

# 2. Install Node.js dependencies
npm install

# 3. Copy environment file
cp .env.example .env

# 4. Start Vite dev server (Terminal 1)
npm run dev

# 5. Start Go server (Terminal 2)
go run main.go
# Or with hot reload: air
```

### Available Commands

| Command | Description |
|---------|-------------|
| `npm run dev` | Start Vite dev server (frontend HMR) |
| `npm run dev:go` | Start Go server with Air (hot reload) |
| `npm run dev:all` | Run both Vite and Air concurrently |
| `go run main.go` | Run Go server directly |
| `npm run build` | Build frontend + Go binary for production |
| `go test ./...` | Run Go tests |
| `npm run test:run` | Run frontend tests |

### Environment Configuration

Required `.env` variables:
```bash
APP_PORT=8080
APP_ENV=development
DB_PATH=./data/app.db
SESSION_SECRET=your-secret-key-change-this
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
```

---

## Development Conventions

### Code Organization
- **Handlers**: Parse requests, call services, return responses
- **Services**: Business logic, authentication flows, user management
- **Repositories**: Database operations using Squirrel query builder
- **Models**: Domain models and DTOs

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
    router.post('/login/login', {
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
# Run all migrations
goose -dir migrations sqlite3 ./data/app.db up

# Check status
goose -dir migrations sqlite3 ./data/app.db status

# Create new migration
goose -dir migrations create migration_name
```

### SQLite Production Optimizations
Applied automatically in `main.go`:
- `journal_mode=WAL` - Write-Ahead Logging
- `synchronous=NORMAL` - Balance speed/durability
- `cache_size=-64000` - 64MB cache
- `temp_store=MEMORY` - Memory temp tables
- `busy_timeout=5000` - 5s lock wait timeout

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
- **Frontend**: Vitest with Happy-DOM, run with `npm run test:run`
- **Test isolation**: Use separate test database if needed

---

## Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| Port 8080 in use | `lsof -ti:8080 | xargs kill -9` |
| Database locked | Remove `data/app.db-shm` and `data/app.db-wal` |
| Vite port detection fails | Delete `.vite-port` and restart Vite |
| Migration errors | Check status with `goose status`, reset if needed |

---

## Documentation References

- `README.md` - Quick start guide
- `docs/DOKUMEN.md` - Complete project documentation
- `docs/FOLDER.md` - Detailed directory structure and architecture explanation
