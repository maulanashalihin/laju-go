# Project Structure

```
laju-go/
├── main.go                          # Application entry point (Go server bootstrap)
├── go.mod                           # Go module dependencies
├── go.sum                           # Go dependency checksums
├── package.json                     # Node.js dependencies & scripts
├── package-lock.json                # Node.js dependency lock file
├── vite.config.js                   # Vite build configuration
├── tsconfig.json                    # TypeScript configuration (frontend)
├── tsconfig.node.json               # TypeScript Node configuration
├── .env                             # Environment variables (gitignored)
├── .env.example                     # Environment template
├── .gitignore                       # Git ignore rules
├── .vite-port                       # Vite dev server port (auto-generated)
├── .air.toml                        # Air hot reload configuration
├── README.md                        # Quick start documentation
├── docs/                            # Documentation folder
├── QWEN.md                          # AI assistant context file
├── laju-go                          # Compiled binary (gitignored)
│
├── app/                             # Go backend application code
│   ├── config/
│   │   └── config.go                # Environment configuration loader
│   ├── handlers/
│   │   ├── app.go                   # Dashboard & profile page handlers
│   │   ├── auth.go                  # Authentication handlers (login, register, OAuth)
│   │   ├── public.go                # Public page handlers (home, about)
│   │   ├── upload.go                # File upload handler
│   │   └── password-reset.go        # Password reset request & completion handlers
│   ├── middlewares/
│   │   ├── auth.go                  # AuthRequired, AdminRequired, Guest middleware
│   │   ├── csrf.go                  # CSRF protection middleware
│   │   └── rate-limit.go            # Rate limiting middleware
│   ├── models/
│   │   ├── dto.go                   # Request/Response DTOs
│   │   ├── session.go               # Session domain model
│   │   └── user.go                  # User domain model
│   ├── repositories/
│   │   ├── session.repository.go    # Session database operations
│   │   └── user.repository.go       # User database operations (Squirrel SQL builder)
│   ├── services/
│   │   ├── asset.go                 # Vite manifest/asset management
│   │   ├── auth.go                  # Authentication business logic
│   │   ├── inertia.go               # Inertia.js rendering service
│   │   ├── mailer.go                # Email sending service (SMTP)
│   │   └── user.go                  # User business logic
│   └── session/
│       └── session.go               # Session infrastructure (cookie encoding/decoding)
│
├── routes/
│   └── web.go                       # Route definitions & setup
│
├── frontend/                        # Svelte 5 frontend application
│   ├── src/
│   │   ├── components/
│   │   │   ├── Button.svelte        # Reusable button component
│   │   │   ├── DarkModeToggle.svelte # Dark/light theme toggle
│   │   │   ├── Header.svelte        # Application header/navigation
│   │   │   └── Input.svelte         # Reusable input component
│   │   ├── layouts/                 # Page layout components (empty - for future use)
│   │   ├── lib/
│   │   │   ├── i18n/
│   │   │   │   └── translation.js   # Internationalization translations
│   │   │   └── utils/
│   │   │       └── helpers.js       # Utility helper functions
│   │   ├── pages/
│   │   │   ├── admin/               # Admin pages (empty - for future use)
│   │   │   ├── app/
│   │   │   │   ├── Dashboard.svelte # Main dashboard page
│   │   │   │   └── Profile.svelte   # User profile page
│   │   │   └── auth/
│   │   │       ├── Login.svelte     # Login page
│   │   │       ├── Register.svelte  # Registration page
│   │   │       ├── ForgotPassword.svelte  # Password reset request page
│   │   │       └── ResetPassword.svelte   # Password reset completion page
│   │   ├── main.ts                  # Inertia.js entry point
│   │   └── app.css                  # Global styles (Tailwind CSS)
│   └── (build artifacts in dist/)
│
├── migrations/
│   ├── 0001_create_users_table.sql  # Users table migration
│   └── 0002_create_sessions_table.sql # Sessions table migration
│
├── data/
│   ├── app.db                       # SQLite database (gitignored)
│   ├── app.db-shm                   # SQLite shared memory file
│   └── app.db-wal                   # SQLite write-ahead log file
│
├── dist/                            # Production build output
│   ├── .vite/
│   │   └── manifest.json            # Vite asset manifest
│   └── assets/
│       ├── app-*.css                # Compiled CSS
│       ├── main-*.js                # Main JS bundle
│       ├── Dashboard-*.js           # Dashboard page chunk
│       ├── Login-*.js               # Login page chunk
│       ├── Profile-*.js             # Profile page chunk
│       └── [page]-*.js              # Other page chunks
│
├── templates/
│   ├── index.html                   # Landing page template
│   └── inertia.html                 # Inertia.js base template
│
├── public/
│   └── .gitkeep                     # Placeholder for static assets
│
├── storage/
│   └── .gitkeep                     # Placeholder for user uploads (avatars)
│
└── tmp/
    └── main                         # Temporary build artifact (Air)
```

---

## Directory Descriptions

| Directory | Purpose |
|-----------|---------|
| `app/` | Core Go backend code following layered architecture |
| `app/config/` | Environment configuration loading |
| `app/handlers/` | HTTP request/response handlers (app, auth, public, upload, password-reset) |
| `app/middlewares/` | Request middleware (auth, csrf, rate-limit) |
| `app/models/` | Domain models (user, session) and DTOs |
| `app/repositories/` | Database access layer (user, session) |
| `app/services/` | Business logic layer (auth, user, mailer, inertia, asset) |
| `app/session/` | Session infrastructure (cookie-based storage) |
| `routes/` | Route definitions and middleware setup |
| `frontend/` | Svelte 5 frontend source code |
| `frontend/src/components/` | Reusable UI components (Button, Input, Header, DarkModeToggle) |
| `frontend/src/lib/` | Utility modules (i18n translations, helpers) |
| `frontend/src/pages/` | Page components organized by feature (auth, app, admin) |
| `migrations/` | Database migration scripts (Goose) |
| `data/` | SQLite database files |
| `dist/` | Production build artifacts (Vite output) |
| `templates/` | HTML templates for server-side rendering |
| `public/` | Static assets served directly |
| `storage/` | User uploaded files (avatars) |
| `tmp/` | Temporary build files (Air hot reload) |

---

## Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Request                          │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│  Routes (routes/web.go)                                 │
│  - Maps URLs to handlers                                │
│  - Applies middleware                                   │
│  - Sets up CSRF protection                              │
│  - Configures mailer service                            │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│  Middleware (app/middlewares/)                          │
│  - auth.go: AuthRequired, AdminRequired, Guest          │
│  - csrf.go: CSRF token validation                       │
│  - rate-limit.go: Request rate limiting                 │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│  Handlers (app/handlers/)                               │
│  - app.go: Dashboard, Profile                           │
│  - auth.go: Login, Register, OAuth, Logout              │
│  - public.go: Index, About                              │
│  - upload.go: File uploads                              │
│  - password-reset.go: Reset request & completion        │
│  - Parse requests, validate, call services, respond     │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│  Services (app/services/)                               │
│  - AuthService: Authentication flows                    │
│  - UserService: User management                         │
│  - InertiaService: Inertia rendering                    │
│  - AssetService: Vite asset resolution                  │
│  - MailerService: Email sending (SMTP)                  │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│  Repositories (app/repositories/)                       │
│  - UserRepository: User CRUD (Squirrel SQL builder)     │
│  - SessionRepository: Session persistence               │
│  - Parameterized queries for SQL injection prevention   │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│  Database (data/app.db - SQLite)                        │
│  - users: User accounts                                 │
│  - sessions: User sessions                              │
│  - Optimized: WAL mode, connection pooling              │
└─────────────────────────────────────────────────────────┘
```

---

## Layer Separation: Why `session/` is Separate from `services/`

The `app/session/` folder is kept separate from `app/services/` because they represent
different architectural layers:

| Layer | Folder | Purpose | Example |
|-------|--------|---------|---------|
| **Infrastructure** | `session/` | Technical implementation (cookie encoding, storage) | `Session.Save()`, `Session.Destroy()` |
| **Business Logic** | `services/` | Domain-specific rules and flows | `AuthService.Login()`, `UserService.Update()` |

### Key Differences

**`app/session/` - Infrastructure Layer**
- Generic session management (can be reused across projects)
- Handles cookie encoding/decoding with `securecookie`
- No knowledge of business domain (users, auth, etc.)
- Used BY services, not part of business logic

**`app/services/` - Business Logic Layer**
- Project-specific business rules
- Authentication flows, user management, email sending
- Uses session infrastructure to persist data
- Contains domain knowledge

### Dependency Relationship

```
services/  →  session/
   │            │
   │            └─→ Infrastructure (cookie storage)
   │
   └─→ Business logic uses session to store user data
```

```go
// Example: AuthService (business logic) uses Session (infrastructure)
func (s *AuthService) Login(email, password string) (*models.User, error) {
    // Business logic: validate credentials
    user, err := s.userRepo.GetByEmail(email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    // Business logic: check password hash
    if !checkPassword(user.Password, password) {
        return nil, ErrInvalidCredentials
    }

    // Infrastructure: store session (generic, no business logic)
    session.Set("user_id", user.ID)
    session.Set("email", user.Email)
    session.Save()

    return user, nil
}
```

### Benefits of Separation

1. **Clear Responsibilities**: Session doesn't know about users/auth, just stores data
2. **Reusability**: `session/` can be used in any Fiber project
3. **Testability**: Can mock session infrastructure when testing services
4. **Flexibility**: Easy to swap session implementation (e.g., cookie → Redis) without changing services

---

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    password TEXT,
    avatar TEXT DEFAULT '',
    role TEXT NOT NULL DEFAULT 'user',
    google_id TEXT UNIQUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_google_id ON users(google_id);
```

### Sessions Table
```sql
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    data TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
```

---

## File Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Go handlers | `{feature}.go` | `auth.go`, `app.go`, `upload.go` |
| Go services | `{feature}.go` | `auth.go`, `mailer.go`, `user.go` |
| Go repositories | `{entity}.repository.go` | `user.repository.go`, `session.repository.go` |
| Go models | `{entity}.go` | `user.go`, `session.go`, `dto.go` |
| Go middlewares | `{feature}.go` | `auth.go`, `csrf.go`, `rate-limit.go` |
| Svelte pages | `{Page}.svelte` | `Login.svelte`, `Dashboard.svelte` |
| Svelte components | `{Component}.svelte` | `Button.svelte`, `Input.svelte` |
| Frontend lib | `{feature}.js` | `translation.js`, `helpers.js` |
| Migrations | `{sequence}_{description}.sql` | `0001_create_users_table.sql` |

---

## Frontend Structure Details

### Components (`frontend/src/components/`)
Reusable UI components used across multiple pages:
- `Button.svelte` - Styled button with variants
- `Input.svelte` - Form input with label and error display
- `Header.svelte` - Application header with navigation
- `DarkModeToggle.svelte` - Theme switcher (light/dark mode)

### Library (`frontend/src/lib/`)
Utility modules and helper functions:
- `i18n/translation.js` - Internationalization translations
- `utils/helpers.js` - Common utility functions

### Pages (`frontend/src/pages/`)
Page components organized by feature:
- `auth/` - Authentication pages (Login, Register, Password Reset)
- `app/` - Authenticated app pages (Dashboard, Profile)
- `admin/` - Admin-only pages (empty, for future development)
