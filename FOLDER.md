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
├── README.md                        # Quick start documentation
├── DOKUMEN.md                       # Full project documentation
├── laju-go                          # Compiled binary (gitignored)
│
├── app/                             # Go backend application code
│   ├── config/
│   │   └── config.go                # Environment configuration loader
│   ├── handlers/
│   │   ├── app.handler.go           # Dashboard & profile page handlers
│   │   ├── auth.handler.go          # Authentication handlers (login, register, OAuth)
│   │   ├── public.handler.go        # Public page handlers (home, about)
│   │   └── upload.handler.go        # File upload handler
│   ├── middlewares/
│   │   └── auth.middleware.go       # Auth, Admin, Guest middleware
│   ├── models/
│   │   ├── dto.model.go             # Request/Response DTOs
│   │   └── user.model.go            # User domain model
│   ├── repositories/
│   │   └── user.repository.go       # User database operations (Squirrel SQL builder)
│   ├── services/
│   │   ├── asset.service.go         # Vite manifest/asset management
│   │   ├── auth.service.go          # Authentication business logic
│   │   ├── inertia.service.go       # Inertia.js rendering service
│   │   └── user.service.go          # User business logic
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
│   │   │   └── Input.svelte         # Reusable input component
│   │   ├── layouts/                 # Page layout components (empty - for future use)
│   │   ├── pages/
│   │   │   ├── Admin/               # Admin pages (empty - for future use)
│   │   │   ├── App/
│   │   │   │   ├── Dashboard.svelte # Main dashboard page
│   │   │   │   └── Profile.svelte   # User profile page
│   │   │   └── Auth/
│   │   │       ├── Login.svelte     # Login page
│   │   │       └── Register.svelte  # Registration page
│   │   ├── main.js                  # Inertia.js entry point
│   │   └── app.css                  # Global styles (Tailwind CSS)
│   └── (build artifacts in dist/)
│
├── migrations/
│   └── 0001_create_users_table.sql  # Database migration (Goose)
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
│       └── Profile-*.js             # Profile page chunk
│
├── templates/
│   ├── index.html                   # Landing page template
│   └── inertia.html                 # Inertia.js base template
│
├── public/
│   └── images/                      # Static images (empty - for future use)
│
├── storage/
│   └── avatars/                     # User uploaded avatars (gitignored)
│
└── tmp/
    └── main                         # Temporary build artifact
```

---

## Directory Descriptions

| Directory | Purpose |
|-----------|---------|
| `app/` | Core Go backend code following layered architecture |
| `app/config/` | Environment configuration loading |
| `app/handlers/` | HTTP request/response handlers |
| `app/middlewares/` | Request middleware (auth, admin, guest) |
| `app/models/` | Domain models and DTOs |
| `app/repositories/` | Database access layer |
| `app/services/` | Business logic layer |
| `app/session/` | Session infrastructure (cookie-based storage, cross-cutting concern) |
| `routes/` | Route definitions |
| `frontend/` | Svelte 5 frontend source code |
| `frontend/src/components/` | Reusable UI components |
| `frontend/src/pages/` | Page components organized by feature |
| `migrations/` | Database migration scripts |
| `data/` | SQLite database files |
| `dist/` | Production build artifacts |
| `templates/` | HTML templates for server-side rendering |
| `public/` | Static assets |
| `storage/` | User uploaded files |
| `tmp/` | Temporary files |

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
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│  Middleware (app/middlewares/)                          │
│  - AuthRequired                                         │
│  - AdminRequired                                        │
│  - Guest                                                │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│  Handlers (app/handlers/)                               │
│  - Parse requests                                       │
│  - Call services                                        │
│  - Return responses                                     │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│  Services (app/services/)                               │
│  - Business logic                                       │
│  - Authentication                                       │
│  - User management                                      │
│  - Inertia rendering                                    │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│  Repositories (app/repositories/)                       │
│  - SQL queries (Squirrel)                               │
│  - CRUD operations                                      │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│  Database (data/app.db - SQLite)                        │
└─────────────────────────────────────────────────────────┘
```

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
- Authentication flows, user management
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
