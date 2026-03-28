# Project Structure

Complete reference for the Laju Go directory structure and file organization.

## Root Directory

```
laju-go/
├── main.go                    # Application entry point
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
├── repositories/
│   ├── session.repository.go  # Session database operations
│   └── user.repository.go     # User database operations
├── services/
│   ├── asset.go               # Vite asset management
│   ├── auth.go                # Authentication logic
│   ├── inertia.go             # Inertia.js rendering
│   ├── mailer.go              # Email service
│   └── user.go                # User business logic
└── session/
    └── session.go             # Session infrastructure
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
    authService *services.AuthService
    store *session.Store
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    // Handle login
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

#### `/app/repositories/`

| File | Purpose |
|------|---------|
| `session.repository.go` | Session CRUD operations |
| `user.repository.go` | User CRUD operations |

**Pattern**: One repository per entity using Squirrel query builder

```go
type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
    // Query database
}
```

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
    userRepo *repositories.UserRepository
}

func (s *AuthService) Login(email, password string) (*User, error) {
    // Business logic
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
│   │   └── (empty - for future use)
│   ├── lib/
│   │   ├── i18n/
│   │   │   ├── en.json
│   │   │   ├── id.json
│   │   │   └── translation.js
│   │   └── utils/
│   │       └── helpers.js
│   ├── pages/
│   │   ├── admin/
│   │   ├── app/
│   │   │   ├── Dashboard.svelte
│   │   │   └── Profile.svelte
│   │   └── auth/
│   │       ├── Login.svelte
│   │       ├── Register.svelte
│   │       ├── ForgotPassword.svelte
│   │       └── ResetPassword.svelte
│   ├── main.ts
│   └── app.css
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
func SetupRoutes(app *fiber.App) {
    // Public routes
    app.Get("/", PublicHandler.Index)
    
    // Protected routes
    app.Get("/app", AuthRequired, AppHandler.Dashboard)
}
```

---

### `/migrations` - Database Migrations

```
migrations/
├── 0001_create_users_table.sql
└── 0002_create_sessions_table.sql
```

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

### `/templates` - HTML Templates

```
templates/
├── index.html                 # Landing page template
└── inertia.html               # Inertia.js base template
```

**Example**:
```html
<!-- inertia.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{ .title }}</title>
    {{ .inertiaHead }}
</head>
<body>
    {{ .inertia }}
</body>
</html>
```

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

### `main.go`

Application entry point:

```go
func main() {
    // Load configuration
    config := config.Load()
    
    // Initialize database
    db := initDatabase(config.DBPath)
    
    // Initialize session store
    store := initSession(config.SessionSecret)
    
    // Initialize repositories
    userRepo := repositories.NewUserRepository(db)
    
    // Initialize services
    authService := services.NewAuthService(userRepo)
    
    // Initialize handlers
    authHandler := handlers.NewAuthHandler(authService, store)
    
    // Setup routes
    app := fiber.New()
    routes.SetupRoutes(app, authHandler)
    
    // Start server
    app.Listen(":" + config.AppPort)
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
  cmd = "go build -o ./tmp/main ."
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
module laju-go

go 1.26

require (
    github.com/gofiber/fiber/v2 v2.52.0
    github.com/mattn/go-sqlite3 v1.14.22
    github.com/Masterminds/squirrel v1.5.4
)
```

---

## File Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Go handlers | `{feature}.go` | `auth.go`, `app.go` |
| Go services | `{feature}.go` | `auth.go`, `user.go` |
| Go repositories | `{entity}.repository.go` | `user.repository.go` |
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
│  Repositories (app/repositories/)   │
│  - Database queries                 │
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
  ├── config/
  ├── session/
  ├── repositories/
  │     └── (depends on: database/sql)
  ├── services/
  │     └── (depends on: repositories/)
  ├── handlers/
  │     └── (depends on: services/, session/)
  └── routes/
        └── (depends on: handlers/, middlewares/)
```

---

## Next Steps

- [Architecture Guide](../guide/architecture.md) - Understanding the layers
- [Routing Guide](../guide/routing.md) - Route definitions
- [Development Workflow](../deployment/development.md) - Working with the codebase
