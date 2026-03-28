# Laju Go - Complete Documentation

## Table of Contents

1. [Overview](#overview)
2. [Technology Stack](#technology-stack)
3. [Getting Started](#getting-started)
4. [Project Structure](#project-structure)
5. [Backend Architecture](#backend-architecture)
6. [Frontend Architecture](#frontend-architecture)
7. [Database](#database)
8. [Authentication](#authentication)
9. [Routes](#routes)
10. [Development Workflow](#development-workflow)
11. [Production Deployment](#production-deployment)
12. [Security](#security)
13. [Environment Variables](#environment-variables)
14. [Features](#features)

---

## Overview

**Laju Go** is a high-performance SaaS boilerplate built with modern technologies:

- **Backend**: Go Fiber (fasthttp-based web framework)
- **Frontend**: Svelte 5 with Vite
- **Database**: SQLite with optimizations for production
- **SPA Bridge**: Inertia.js for server-driven single-page applications

This stack provides the performance of Go with the developer experience of modern frontend frameworks, without the complexity of building a separate API.

### Key Features

- ✅ User authentication (email/password + Google OAuth)
- ✅ Password reset via email (SMTP)
- ✅ Role-based access control (Admin/User)
- ✅ Database-backed session management (persistent)
- ✅ CSRF protection middleware
- ✅ Rate limiting middleware (auth, password reset)
- ✅ File upload support (avatars)
- ✅ Hot module replacement (HMR) in development
- ✅ Production-ready build pipeline
- ✅ SQLite optimized for production use
- ✅ Clean layered architecture
- ✅ Dark mode support

---

## Technology Stack

### Backend

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.26+ | Programming language |
| Fiber | v2.52.0 | Web framework (fasthttp) |
| SQLite3 | - | Database |
| Squirrel | v1.5.4 | SQL query builder |
| Goose | v3.20.0 | Database migrations |
| OAuth2 | v0.18.0 | Google OAuth |
| Crypto | v0.21.0 | Password hashing (bcrypt) |
| Godotenv | v1.5.1 | Environment variables |

### Frontend

| Technology | Version | Purpose |
|------------|---------|---------|
| Svelte | 5.0.0 | UI framework |
| Vite | 5.0.0 | Build tool |
| Inertia.js | 3.0.0 | SPA bridge |
| Tailwind CSS | 4.0.0 | Styling |
| Lucide Svelte | 1.0.0 | Icons |
| Axios | 1.13.6 | HTTP client |
| Day.js | 1.11.20 | Date handling |
| Dotenv | 17.3.1 | Environment variables (frontend) |
| Vitest | 4.1.2 | Testing framework |

---

## Getting Started

### Prerequisites

- Go 1.26 or higher
- Node.js 18 or higher
- SQLite3

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/maulanashalihin/laju-go.git
   cd laju-go
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Install Node.js dependencies**
   ```bash
   npm install
   ```

4. **Copy environment file**
   ```bash
   cp .env.example .env
   ```

5. **Update environment variables**
   Edit `.env` with your configuration (see [Environment Variables](#environment-variables))

6. **Database migrations**
   Migrations run automatically on startup. Manual control:
   ```bash
   goose -dir migrations sqlite3:data/app.db up
   ```

7. **Start development servers**

   Terminal 1 - Vite dev server:
   ```bash
   npm run dev
   ```

   Terminal 2 - Go server:
   ```bash
   go run main.go
   ```
   
   Or use Air for hot reload:
   ```bash
   air
   # or
   npm run dev:go
   ```

8. **Open browser**
   Navigate to `http://localhost:8080`

---

## Project Structure

See [FOLDER.md](FOLDER.md) for detailed directory structure.

### Quick Reference

```
laju-go/
├── main.go              # Entry point
├── app/                 # Backend code
│   ├── handlers/        # HTTP handlers (auth, app, public, upload, password-reset)
│   ├── services/        # Business logic (auth, user, mailer, inertia, asset)
│   ├── repositories/    # Database layer (user, session)
│   ├── models/          # Data models
│   ├── middlewares/     # Request middleware
│   └── session/         # Session infrastructure
├── routes/              # Route definitions
├── frontend/            # Svelte frontend
│   └── src/
│       ├── components/  # Reusable components (Button, Input, Header, DarkModeToggle)
│       └── pages/       # Page components (auth/, app/, admin/)
├── migrations/          # DB migrations (users, sessions)
├── templates/           # HTML templates
└── storage/             # Uploaded files (avatars)
```

---

## Backend Architecture

### Layered Architecture

```
Routes → Middleware → Handler → Service → Repository → Database
```

### Layers Explained

1. **Routes** (`routes/web.go`)
   - Defines URL endpoints
   - Maps routes to handlers
   - Applies middleware chains
   - Sets up CSRF protection
   - Configures mailer service

2. **Middleware** (`app/middlewares/`)
   - `AuthRequired`: Protects routes, requires authenticated user
   - `AdminRequired`: Requires admin role
   - `Guest`: Redirects authenticated users away from auth pages

3. **Handlers** (`app/handlers/`)
   - `auth.go`: Login, register, OAuth callbacks, logout
   - `app.go`: Dashboard, profile pages
   - `public.go`: Home, about pages
   - `upload.go`: File upload handling
   - `password-reset.go`: Password reset requests and completion
   - Parse HTTP requests, validate input, call services, return responses

4. **Services** (`app/services/`)
   - `AuthService`: Authentication flows (email/password, OAuth)
   - `UserService`: User management (CRUD, profile updates)
   - `MailerService`: Email sending via SMTP (password resets)
   - `InertiaService`: Inertia.js rendering and responses
   - `AssetService`: Vite asset manifest resolution

5. **Repositories** (`app/repositories/`)
   - `UserRepository`: User database operations using Squirrel
   - `SessionRepository`: Session persistence
   - Parameterized queries for SQL injection prevention

6. **Models** (`app/models/`)
   - `User`: Domain model
   - `DTOs`: Request/Response data transfer objects

7. **Session** (`app/session/`)
   - Infrastructure layer for session storage
   - Cookie encoding/decoding with securecookie
   - Generic implementation (reusable across projects)
   - Database-backed session persistence

### Example Flow: User Login

```
POST /login/login
    ↓
routes/web.go → AuthHandler.Login()
    ↓
AuthService.LoginByEmail()
    ↓
UserRepository.GetByEmail()
    ↓
bcrypt.CompareHashAndPassword()
    ↓
Session created in database → Inertia redirect to /app
```

### Example Flow: Password Reset

```
POST /password-reset/request
    ↓
PasswordResetHandler.RequestReset()
    ↓
UserService.GenerateResetToken()
    ↓
MailerService.SendResetEmail()
    ↓
Email sent with reset link → Success response
```

---

## Frontend Architecture

### Inertia.js Pattern

Laju Go uses Inertia.js to bridge backend and frontend without building a separate API:

1. **Initial Load**: Server renders HTML via `inertia.html` template
2. **Subsequent Navigation**: Inertia makes XHR requests with `X-Inertia: true` header
3. **Server Response**: Returns JSON with component name and props
4. **Client Rendering**: Svelte dynamically loads the component

### Component Structure

```
frontend/src/
├── main.ts                  # Inertia initialization
├── app.css                  # Global styles (Tailwind)
├── components/              # Reusable UI components
│   ├── Button.svelte        # Button component
│   ├── Input.svelte         # Form input component
│   ├── Header.svelte        # Application header
│   └── DarkModeToggle.svelte # Theme toggle
└── pages/                   # Page components
    ├── auth/
    │   ├── Login.svelte     # Login page
    │   ├── Register.svelte  # Registration page
    │   ├── ForgotPassword.svelte  # Password reset request
    │   └── ResetPassword.svelte   # Password reset completion
    ├── app/
    │   ├── Dashboard.svelte # Main dashboard
    │   └── Profile.svelte   # User profile
    └── admin/               # Admin pages (future)
```

### Making Inertia Requests

```svelte
<script>
  import { router } from '@inertiajs/svelte'

  function handleSubmit() {
    router.post('/login/login', {
      email: formData.email,
      password: formData.password,
    }, {
      onError: (errors) => {
        // Handle validation errors
      }
    })
  }
</script>
```

### Form Validation

Inertia automatically handles validation errors returned from the server:

```svelte
<script>
  import { page } from '@inertiajs/svelte'
  
  $: errors = page.props.errors || {}
</script>

{#if errors.email}
  <span class="error">{errors.email}</span>
{/if}
```

---

## Database

### Schema

**Users Table**

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

**Sessions Table**

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

### SQLite Production Optimizations

Applied in `main.go`:

```go
db.Exec(`PRAGMA journal_mode=WAL`)           // Write-Ahead Logging
db.Exec(`PRAGMA synchronous=NORMAL`)          // Balance speed/durability
db.Exec(`PRAGMA cache_size=-64000`)           // 64MB cache
db.Exec(`PRAGMA temp_store=MEMORY`)           // Memory temp tables
db.Exec(`PRAGMA busy_timeout=5000`)           // 5s lock wait timeout
```

### Connection Pool

```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

### Migrations

Using Goose for database migrations:

```bash
# Create new migration
goose -dir migrations create add_users_table

# Run all migrations
goose -dir migrations sqlite3:data/app.db up

# Check status
goose -dir migrations sqlite3:data/app.db status

# Rollback last migration
goose -dir migrations sqlite3:data/app.db down
```

---

## Authentication

### Authentication Methods

1. **Email/Password**
   - Registration with email, name, password
   - Password hashed with bcrypt
   - Login validates credentials

2. **Google OAuth**
   - OAuth 2.0 flow
   - Automatic user creation if not exists
   - Links Google ID to user account

### Password Reset

1. User requests reset via `/password-reset/request`
2. System generates token and stores in session
3. Email sent with reset link containing token
4. User clicks link → `/password-reset/reset?token=xxx`
5. User submits new password
6. Password updated, sessions invalidated

### Session Management

**Architecture**: `app/session/` is a separate infrastructure layer (not part of `app/services/`)

| Layer | Folder | Purpose |
|-------|--------|---------|
| Infrastructure | `session/` | Cookie encoding/decoding, storage mechanics |
| Business Logic | `services/` | Authentication rules, user management |

This separation allows:
- Reusability: `session/` works in any Fiber project
- Flexibility: Easy to swap implementation (cookie → Redis)
- Clear responsibilities: Session doesn't know about business domain

**Session Configuration**:
- Database-backed sessions (persistent across restarts)
- HTTPOnly cookies for security
- SameSite=Lax for CSRF protection
- SecureCookie encoding with secret key
- Configurable expiration

### Session Data

```go
// Session storage (infrastructure layer)
session.Set("user_id", user.ID)
session.Set("email", user.Email)
session.Set("role", user.Role)
session.Save()
```

### Middleware Protection

```go
// Protected route - requires authentication
app.Get("/app", middlewares.AuthRequired(store), AppHandler.Dashboard)

// Admin only
app.Get("/admin", middlewares.AdminRequired(store), AdminHandler.Dashboard)

// Guest only (redirect if authenticated)
app.Get("/login", middlewares.Guest(store), AuthHandler.ShowLoginForm)
```

---

## Routes

### Public Routes

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/` | `PublicHandler.Index` | Landing page |
| GET | `/about` | `PublicHandler.About` | About page |

### Authentication Routes

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/login` | `AuthHandler.ShowLoginForm` | Login form (Guest only) |
| POST | `/login` | `AuthHandler.Login` | Login submission (rate-limited) |
| GET | `/register` | `AuthHandler.ShowRegisterForm` | Registration form (Guest only) |
| POST | `/register` | `AuthHandler.Register` | Registration (rate-limited) |
| GET | `/auth/google` | `AuthHandler.GoogleLogin` | OAuth start |
| GET | `/auth/google/callback` | `AuthHandler.GoogleCallback` | OAuth callback |
| POST | `/logout` | `AuthHandler.Logout` | Logout (requires auth) |
| GET | `/api/me` | `AuthHandler.Me` | Current user API (requires auth) |

### Password Reset Routes

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/forgot-password` | `PasswordResetHandler.ShowForgotPasswordForm` | Request reset form |
| POST | `/forgot-password` | `PasswordResetHandler.SendResetLink` | Send reset email (rate-limited) |
| GET | `/reset-password/:token` | `PasswordResetHandler.ShowResetPasswordForm` | Reset password form |
| POST | `/reset-password/:token` | `PasswordResetHandler.ResetPassword` | Complete reset |

### Protected App Routes

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/app` | `AppHandler.Dashboard` | Dashboard (CSRF protected) |
| GET | `/app/profile` | `AppHandler.Profile` | Profile page (CSRF protected) |
| PUT | `/app/profile` | `AppHandler.UpdateProfile` | Update profile (CSRF protected) |
| PUT | `/app/profile/password` | `AppHandler.UpdatePassword` | Update password (CSRF protected) |
| POST | `/upload` | `UploadHandler.Upload` | File upload (CSRF protected) |

### Admin Routes

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/admin` | (inline) | Admin dashboard |

---

## Development Workflow

### Development Mode

1. **Start Vite dev server** (auto-detects port):
   ```bash
   npm run dev
   ```
   - Writes port to `.vite-port`
   - Enables HMR (Hot Module Replacement)

2. **Start Go server**:
   
   Option A - Direct run:
   ```bash
   go run main.go
   ```
   
   Option B - With Air (hot reload):
   ```bash
   air
   # or
   npm run dev:go
   ```
   
   Option C - Both concurrently:
   ```bash
   npm run dev:all
   ```

3. **Auto-reload**:
   - Air automatically rebuilds on `.go` file changes (~1-2 sec)
   - Vite HMR updates frontend instantly

### Available Scripts

```bash
# Development
npm run dev          # Start Vite dev server
npm run dev:go       # Start Go server with Air
npm run dev:all      # Run both Vite and Air concurrently
go run main.go       # Start Go server directly

# Production build
npm run build        # Build frontend + Go binary
npm run serve        # Run production binary

# Testing
npm run test:run     # Run frontend tests
npm run test:ui      # Run tests with UI
```

### Vite Port Detection

The application uses a custom Vite plugin to automatically detect the dev server port:

1. Vite writes port to `.vite-port` file on startup
2. Go server reads `.vite-port` to proxy requests to Vite
3. Cleanup on Vite exit

This allows multiple instances to run without port conflicts.

---

## Production Deployment

### Build Steps

1. **Build frontend**:
   ```bash
   npm run build
   ```
   - Compiles Svelte components
   - Generates hashed assets in `dist/`
   - Creates manifest in `dist/.vite/manifest.json`

2. **Build Go binary**:
   ```bash
   go build -o laju-go .
   ```

3. **Run migrations**:
   ```bash
   goose -dir migrations sqlite3:data/app.db up
   ```

4. **Start server**:
   ```bash
   ./laju-go
   ```

### Environment Configuration

Set production environment variables:

```bash
APP_ENV=production
APP_PORT=8080
DB_PATH=data/app.db
SESSION_SECRET=<strong-random-secret>
# ... other variables
```

### Process Management (systemd)

```ini
# /etc/systemd/system/laju-go.service
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
StandardOutput=journal
StandardError=journal
NoNewPrivileges=true
PrivateTmp=true
Environment="PATH=/usr/local/go/bin:/usr/bin:/bin"

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable laju-go
sudo systemctl start laju-go
sudo systemctl status laju-go
```

### Docker Deployment

```dockerfile
# Multi-stage build
FROM node:20-alpine AS frontend
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM golang:1.22-alpine AS backend
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
COPY --from=frontend /app/dist ./dist
RUN go build -o laju-go .

FROM alpine:latest
WORKDIR /app
COPY --from=backend /app/laju-go .
EXPOSE 8080
CMD ["./laju-go"]
```

```bash
# Build and run
docker build -t laju-go .
docker run -p 8080:8080 -v $(pwd)/data:/app/data laju-go
```

---

## Security

### Implemented Security Measures

| Feature | Implementation |
|---------|----------------|
| Password Hashing | bcrypt with cost factor |
| SQL Injection | Squirrel parameterized queries |
| XSS Protection | HTML escaping in templates |
| CSRF Protection | SameSite cookies, OAuth state, CSRF tokens |
| Session Security | HTTPOnly cookies, SecureCookie encoding |
| File Upload | Type validation, size limits (5MB) |
| Role-Based Access | Admin/User middleware guards |
| Foreign Keys | SQLite foreign key constraints |

### Best Practices

1. **Never commit `.env`** - Contains secrets
2. **Use HTTPS in production** - Set `APP_ENV=production`
3. **Rotate session secrets** - Change `SESSION_SECRET` regularly
4. **Validate all input** - File uploads, form data
5. **Rate limiting** - Consider adding for auth endpoints
6. **Secure cookies** - Set `Secure` flag in production

---

## Environment Variables

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `APP_ENV` | Environment | `development`, `production` |
| `APP_PORT` | Server port | `8080` |
| `DB_PATH` | Database path | `data/app.db` |
| `SESSION_SECRET` | Session encryption key | `your-32-char-secret` |

### OAuth Variables (Optional)

| Variable | Description | Example |
|----------|-------------|---------|
| `GOOGLE_CLIENT_ID` | Google OAuth client ID | `xxx.apps.googleusercontent.com` |
| `GOOGLE_CLIENT_SECRET` | Google OAuth secret | `GOCSPX-xxx` |
| `GOOGLE_REDIRECT_URL` | OAuth callback URL | `http://localhost:8080/auth/google/callback` |

### Email Variables (Password Reset)

| Variable | Description | Example |
|----------|-------------|---------|
| `SMTP_HOST` | SMTP server | `smtp.gmail.com` |
| `SMTP_PORT` | SMTP port | `587` |
| `SMTP_USER` | SMTP username | `your-email@gmail.com` |
| `SMTP_PASS` | SMTP password | `your-app-password` |
| `FROM_EMAIL` | Sender email | `noreply@example.com` |
| `FROM_NAME` | Sender name | `Laju Go` |

### Example `.env`

```bash
# Application
APP_ENV=development
APP_PORT=8080

# Database
DB_PATH=data/app.db

# Session
SESSION_SECRET=your-32-character-secret-key-change-in-production

# Google OAuth (optional)
GOOGLE_CLIENT_ID=xxx.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-xxx
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# Email (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
FROM_EMAIL=noreply@example.com
FROM_NAME=Laju Go
```

---

## Features

### User Management
- Registration with email/password
- Profile editing (name, email)
- Avatar upload
- Role-based access (user/admin)

### Authentication
- Email/password login
- Google OAuth integration
- Session persistence (database-backed)
- Remember me functionality
- Password reset via email

### Frontend
- Svelte 5 with reactive components
- Inertia.js for SPA experience
- Tailwind CSS for styling
- Dark mode toggle
- Responsive design
- Form validation with error display

### Development Experience
- Hot module replacement (Vite)
- Go hot reload (Air)
- TypeScript support
- Component-based architecture
- Automatic asset versioning

### Production Ready
- SQLite optimizations (WAL, connection pooling)
- Database migrations (Goose)
- Process management (systemd)
- Docker support
- Environment-based configuration

---

## Troubleshooting

### Common Issues

**1. Port already in use**
```bash
# Kill process on port 8080
lsof -ti:8080 | xargs kill -9
```

**2. Database locked**
```bash
# Remove WAL files
rm data/app.db-shm data/app.db-wal
```

**3. Vite port detection fails**
```bash
# Remove .vite-port and restart
rm .vite-port
npm run dev
```

**4. Migration errors**
```bash
# Check migration status
goose -dir migrations sqlite3:data/app.db status

# Reset if needed
goose -dir migrations sqlite3:data/app.db down
```

**5. Google OAuth not working**
- Verify redirect URI matches exactly
- Check Google Cloud Console credentials
- Ensure OAuth consent screen is configured

**6. Email not sending**
- Use app-specific password for Gmail
- Verify SMTP credentials
- Check firewall/port blocking

---

## Contributing

1. Create feature branch
2. Make changes
3. Run tests
4. Submit pull request

---

## License

MIT License - see [LICENSE](LICENSE) for details.

---

## Acknowledgments

- [Go Fiber](https://gofiber.io/) - Fast web framework
- [Svelte](https://svelte.dev/) - Cybernetically enhanced web apps
- [Inertia.js](https://inertiajs.com/) - Server-driven SPA
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS
- [Lucide Icons](https://lucide.dev/) - Beautiful icons
