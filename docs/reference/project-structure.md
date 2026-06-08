# Project Structure

Complete reference for the Laju Go directory structure and file organization.

## Root Directory

```
laju-go/
├── cmd/laju-go/main.go        # Application entry point
├── go.mod                     # Go module dependencies
├── go.sum                     # Go dependency checksums
├── package.json               # Node.js dependencies & scripts
├── package-lock.json          # Node.js dependency lock file
├── vite.config.js             # Vite build configuration
├── tsconfig.json              # TypeScript configuration
├── .env                       # Environment variables (gitignored)
├── .env.example               # Environment template
├── .gitignore                 # Git ignore rules
├── .air.toml                  # Air hot reload configuration
├── README.md                  # Project documentation
└── docs/                      # Documentation folder
```

## Application Directories

### `/app` - Backend Go Code

Core application logic organized by architectural layer.

```
app/
├── cache/
│   └── user_cache.go          # In-memory TTL user profile cache
├── config/
│   └── config.go              # Environment configuration loader
├── handlers/
│   ├── app.go                 # Dashboard & profile handlers
│   ├── auth.go                # Authentication handlers
│   ├── public.go              # Public page handlers
│   ├── upload.go              # File upload handler
│   └── password-reset.go      # Password reset handlers
├── middlewares/
│   ├── auth.go                # Auth & role middleware
│   ├── csrf.go                # CSRF protection
│   └── rate-limit.go          # Rate limiting
├── models/
│   ├── dto.go                 # Request/Response DTOs
│   ├── session.go             # Session model
│   └── user.go                # User model
├── queries/                   # sqlc generated (DO NOT EDIT)
│   ├── db.go                  # Database connection wrapper
│   ├── models.go              # Generated models
│   ├── querier.go             # Querier interface + wrapper
│   ├── session.sql.go         # Session queries
│   ├── session_helpers.go     # Session helpers
│   └── user.sql.go            # User queries
├── services/
│   ├── asset.go               # Vite asset management
│   ├── auth.go                # Authentication logic
│   ├── inertia.go             # Inertia.js rendering
│   ├── mailer.go              # Email service
│   └── user.go                # User business logic
└── session/
    └── session.go             # Session infrastructure (separate from services)
```

#### `/app/config/`

| File | Purpose |
|------|---------|
| `config.go` | Loads and validates environment variables |

**Example**:
```go
// app/config/config.go
type Config struct {
    AppEnv      string
    AppPort     string
    DBPath      string
    SessionSecret string
}

func Load() *Config {
    // Load from .env or environment
}
```

#### `/app/handlers/`

| File | Purpose |
|------|---------|
| `app.go` | Dashboard, profile page handlers |
| `auth.go` | Login, register, OAuth, logout |
| `public.go` | Home, about pages |
| `upload.go` | File upload handling |
| `password-reset.go` | Password reset flow |

**Pattern**: Struct-based handlers with dependency injection

```go
type AuthHandler struct {
    authService    *services.AuthService
    userService    *services.UserService
    store          *session.Store
    inertiaService *services.InertiaService
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req models.LoginRequest
    c.BodyParser(&req)

    user, err := h.authService.Login(req.Email, req.Password)
    if err != nil {
        h.store.Flash(c, "error", "Invalid email or password")
        return c.Redirect("/login")
    }

    sess, _ := h.store.Get(c)
    sess.Set("user_id", user.ID)
    sess.Save()
    return c.Redirect("/app")
}
```

#### `/app/middlewares/`

| File | Purpose |
|------|---------|
| `auth.go` | `AuthRequired`, `AdminRequired`, `Guest` |
| `csrf.go` | CSRF token validation |
| `rate-limit.go` | Request rate limiting |

**Example**:
```go
func AuthRequired(store *session.Store) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Check session
        return c.Next()
    }
}
```

#### `/app/models/`

| File | Purpose |
|------|---------|
| `dto.go` | Data Transfer Objects for requests/responses |
| `session.go` | Session domain model |
| `user.go` | User domain model |

**Example**:
```go
type User struct {
    ID        int
    Email     string
    Name      string
    Password  string
    Role      string
    CreatedAt time.Time
}
```

#### `/app/queries/` (sqlc Generated)

| File | Purpose |
|------|---------|
| `db.go` | Database transaction helpers |
| `models.go` | Generated Go structs matching DB schema |
| `querier.go` | Querier wrapper struct |
| `user.sql.go` | User CRUD queries |
| `session.sql.go` | Session CRUD queries |
| `session_helpers.go` | Session helper functions |

**Important**: This directory is **auto-generated by sqlc** from `queries/*.sql`. Never edit manually.

**Pattern**: Write SQL → generate → use typed methods

```go
// Instead of hand-writing repositories, use sqlc-generated queries:
user, err := s.querier.GetUserByEmail(ctx, email)
```

#### `/app/cache/`

| File | Purpose |
|------|---------|
| `user_cache.go` | In-memory TTL cache for user profiles |

Used by `UserService` to reduce DB queries. Cache is automatically invalidated on profile updates.

#### `/app/services/`

| File | Purpose |
|------|---------|
| `asset.go` | Vite manifest parsing, asset URLs |
| `auth.go` | Authentication business logic |
| `inertia.go` | Inertia.js response rendering |
| `mailer.go` | SMTP email sending |
| `user.go` | User management logic |

**Example**:
```go
type AuthService struct {
    querier *queries.Querier
}

func (s *AuthService) Login(email, password string) (*User, error) {
    user, err := s.querier.GetUserByEmail(context.Background(), email)
    // ... business logic ...
}
```

#### `/app/session/`

| File | Purpose |
|------|---------|
| `session.go` | Session storage infrastructure |

**Note**: Separate from services for reusability

---

### `/frontend` - Svelte 5 Frontend

```
frontend/
├── src/
│   ├── components/
│   │   ├── Button.svelte
│   │   ├── Input.svelte
│   │   ├── Header.svelte
│   │   └── DarkModeToggle.svelte
│   ├── layouts/
│   │   └── (add layout components here)
│   ├── lib/
│   │   ├── api/                    # API client helpers
│   │   ├── i18n/                   # Internationalization (EN/ID)
│   │   │   ├── en.json
│   │   │   ├── id.json
│   │   │   └── translation.js
│   │   ├── types/                  # TypeScript type definitions
│   │   └── utils/                  # Helper functions
│   │       └── helpers.js
│   ├── pages/
│   │   ├── admin/                  # Admin-only pages (future)
│   │   ├── app/
│   │   │   ├── Dashboard.svelte
│   │   │   └── Profile.svelte
│   │   └── auth/
│   │       ├── Login.svelte
│   │       ├── Register.svelte
│   │       ├── ForgotPassword.svelte
│   │       └── ResetPassword.svelte
│   ├── main.ts                     # Inertia app initialization
│   └── app.css                     # Global styles (Tailwind)
├── package.json
└── vite.config.js
```

#### `/frontend/src/components/`

Reusable UI components:

| Component | Purpose |
|-----------|---------|
| `Button.svelte` | Styled button with variants |
| `Input.svelte` | Form input with label and error |
| `Header.svelte` | Application header/navigation |
| `DarkModeToggle.svelte` | Light/dark theme toggle |

#### `/frontend/src/pages/`

Page components organized by feature:

| Directory | Purpose |
|-----------|---------|
| `admin/` | Admin-only pages (future) |
| `app/` | Authenticated user pages |
| `auth/` | Authentication pages |

#### `/frontend/src/lib/`

Utility modules:

| Directory | Purpose |
|-----------|---------|
| `i18n/` | Internationalization (EN/ID) |
| `utils/` | Helper functions |

---

### `/routes` - Route Definitions

```
routes/
└── web.go                     # All route definitions
```

**Example**:
```go
func SetupRoutes(app *fiber.App, handlers Handlers, store *session.Store, mailerService *services.MailerService, csrfMiddleware *middlewares.CSRFMiddleware) {
    // Public routes
    app.Get("/", handlers.Public.Index)
    app.Get("/about", handlers.Public.About)

    // Auth routes (with Guest middleware)
    app.Get("/login", middlewares.Guest(store), handlers.Auth.ShowLoginForm)
    app.Post("/login", middlewares.Guest(store), handlers.Auth.Login, middlewares.AuthRateLimit.Limit())

    // Protected routes (AuthRequired + CSRF)
    protected := app.Group("/app", middlewares.AuthRequired(store))
    protected.Use(csrfMiddleware.Protect())
    protected.Get("/", handlers.App.Dashboard)
}
```

---

### `/migrations` - Database Migrations

```
migrations/
├── 0001_create_users_table.sql
└── 0002_create_sessions_table.sql
```

Migrations run **automatically on startup** via Goose (`goose.Up(db, "./migrations")`).

**Naming**: `NNNN_description.sql` (N = sequence number)

**Example**:
```sql
-- 0001_create_users_table.sql
-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    -- ...
);

-- +goose Down
DROP TABLE IF EXISTS users;
```

---

### `/data` - Database Files

```
data/
├── app.db                     # SQLite database (gitignored)
├── app.db-shm                 # Shared memory file
└── app.db-wal                 # Write-ahead log
```

**Note**: Gitignored - created at runtime

---

### `/dist` - Production Build

```
dist/
├── .vite/
│   └── manifest.json          # Asset manifest
└── assets/
    ├── app-*.css              # Compiled CSS
    ├── main-*.js              # Main bundle
    └── [page]-*.js            # Page chunks
```

**Note**: Generated by `npm run build`

---

### `/templates` - Templ Templates

HTML templates written as [templ](https://templ.guide/) components — type-safe, compiled Go templates with JSX-like syntax.

```
templates/
├── index.templ                # Landing page template
├── index_templ.go             # Generated Go code (do not edit)
├── inertia.templ              # Inertia.js base template
└── inertia_templ.go           # Generated Go code (do not edit)
```

**Example**:
```templ
// inertia.templ
package templates

templ InertiaPage(title string, pageJSON string, viteServerURL string, mainJS string, mainCSS string, styles []string) {
    <!doctype html>
    <html lang="en">
        <head>
            <meta charset="UTF-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
            <title>{ title } - Laju</title>
        </head>
        <body class="bg-gray-50 text-gray-900">
            <div id="app"></div>
            <script data-page="app" type="application/json">
                { pageJSON }
            </script>
        </body>
    </html>
}
```

**Workflow**:
1. Edit `.templ` files
2. Run `templ generate` (or `templ generate -watch` for auto-regeneration)
3. Commit both `.templ` source and `_templ.go` generated files

---

### `/public` - Static Assets

```
public/
└── .gitkeep                   # Placeholder
```

For static files served directly (images, fonts, etc.)

---

### `/storage` - User Uploads

```
storage/
└── avatars/                   # User avatar uploads
```

**Note**: Gitignored - created at runtime

---

### `/tmp` - Build Artifacts

```
tmp/
└── main                       # Air build output
```

**Note**: Gitignored - auto-generated by Air

---

## Configuration Files

### `cmd/laju-go/main.go`

Application entry point:

```go
func main() {
    cfg := config.Load()

    // Initialize database + run migrations
    db, _ := initDatabase(cfg.DBPath)
    runMigrations(db, "./migrations")

    // Initialize querier (sqlc-generated)
    querier := queries.NewQuerier(db)

    // Initialize session store + cache
    sessionStore := session.New(querier)
    userCache := cache.NewUserCache(cfg.UserCacheTTL)

    // Initialize services
    authService := services.NewAuthService(querier, services.AuthServiceConfig{
        SessionSecret: cfg.SessionSecret,
        // ... OAuth config ...
    })
    userService := services.NewUserService(querier, userCache)
    inertiaService := services.NewInertiaService(assetService, sessionStore)

    // Initialize handlers
    routeHandlers := routes.Handlers{
        Public: handlers.NewPublicHandler(authService, userService, inertiaService, assetService),
        Auth:   handlers.NewAuthHandler(authService, userService, sessionStore, inertiaService),
        App:    handlers.NewAppHandler(userService, sessionStore, inertiaService),
        Upload: handlers.NewUploadHandler(sessionStore, userService),
    }

    // Setup routes + start
    app := fiber.New()
    routes.SetupRoutes(app, routeHandlers, sessionStore, mailerService, csrfMiddleware)
    app.Listen(":" + cfg.AppPort)
}
```

### `vite.config.js`

Vite build configuration:

```javascript
import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  plugins: [svelte()],
  server: { port: 5173 },
  build: {
    outDir: 'dist',
    manifest: true,
  },
})
```

### `.air.toml`

Air hot reload configuration:

```toml
[build]
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/laju-go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "node_modules"]
```

### `package.json`

NPM scripts and dependencies:

```json
{
  "scripts": {
    "dev": "vite",
    "dev:go": "air",
    "dev:all": "concurrently \"npm run dev\" \"npm run dev:go\"",
    "build": "vite build"
  }
}
```

### `go.mod`

Go module dependencies:

```go
module github.com/maulanashalihin/laju-go

go 1.26

require (
    github.com/gofiber/fiber/v2 v2.52.13
    github.com/a-h/templ v0.3.1001
    github.com/pressly/goose/v3 v3.20.0
    modernc.org/sqlite v1.39.1
)
```

---

## File Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Go handlers | `{feature}.go` | `auth.go`, `app.go` |
| Go services | `{feature}.go` | `auth.go`, `user.go` |
| Go queries (sqlc) | `{entity}.sql.go` | `user.sql.go` |
| Go models | `{entity}.go` | `user.go`, `session.go` |
| Go middlewares | `{feature}.go` | `auth.go`, `csrf.go` |
| Svelte pages | `{Page}.svelte` | `Login.svelte` |
| Svelte components | `{Component}.svelte` | `Button.svelte` |
| Migrations | `{seq}_{desc}.sql` | `0001_create_users.sql` |

---

## Architecture Layers

```
┌─────────────────────────────────────┐
│         HTTP Request                │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Routes (routes/web.go)             │
│  - Map URLs to handlers             │
│  - Apply middleware                 │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Middleware (app/middlewares/)      │
│  - Auth, CSRF, Rate Limit           │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Handlers (app/handlers/)           │
│  - Parse request, call services     │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Services (app/services/)           │
│  - Business logic                   │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Queries — sqlc (app/queries/)      │
│  - Type-safe generated SQL methods  │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Database (data/app.db)             │
│  - SQLite storage                   │
└─────────────────────────────────────┘
```

---

## Dependency Graph

```
main.go
  ├── config/          → environment variables
  ├── database/sql     → modernc.org/sqlite (pure Go)
  ├── queries/         → sqlc generated (depends on: database/sql)
  │     └── (type-safe SQL methods)
  ├── session/         → (depends on: queries/)
  ├── cache/           → (standalone in-memory TTL cache)
  ├── services/        → (depends on: queries/, session/, cache/)
  ├── handlers/        → (depends on: services/, session/)
  └── routes/          → (depends on: handlers/, middlewares/)
```

---

## Next Steps

- [Architecture Guide](../guide/architecture.md) - Understanding the layers
- [Routing Guide](../guide/routing.md) - Route definitions
- [Development Workflow](../deployment/development.md) - Working with the codebase
