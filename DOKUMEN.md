# VeloStack Go - Complete Documentation

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

---

## Overview

**VeloStack Go** is a high-performance SaaS boilerplate built with modern technologies:

- **Backend**: Go Fiber (fasthttp-based web framework)
- **Frontend**: Svelte 5 with Vite
- **Database**: SQLite with optimizations for production
- **SPA Bridge**: Inertia.js for server-driven single-page applications

This stack provides the performance of Go with the developer experience of modern frontend frameworks, without the complexity of building a separate API.

### Key Features

- ✅ User authentication (email/password + Google OAuth)
- ✅ Role-based access control (Admin/User)
- ✅ Session management with secure cookies
- ✅ File upload support
- ✅ Hot module replacement (HMR) in development
- ✅ Production-ready build pipeline
- ✅ SQLite optimized for production use
- ✅ Clean layered architecture

---

## Technology Stack

### Backend

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.22+ | Programming language |
| Fiber | v2.52.0 | Web framework (fasthttp) |
| SQLite3 | - | Database |
| Squirrel | v1.2.0 | SQL query builder |
| Gorilla Sessions | v1.2.1 | Session management |
| OAuth2 | v0.30.0 | Google OAuth |
| Bcrypt | v0.6.0 | Password hashing |
| Goose | v3.18.0 | Database migrations |
| Godotenv | v1.5.1 | Environment variables |

### Frontend

| Technology | Version | Purpose |
|------------|---------|---------|
| Svelte | 5.0.0 | UI framework |
| Vite | 5.0.0 | Build tool |
| Inertia.js | 1.0.0 | SPA bridge |
| Tailwind CSS | 4.0.0 | Styling |
| Lucide Svelte | 0.469.0 | Icons |
| Axios | 1.7.9 | HTTP client |
| Day.js | 1.11.13 | Date handling |
| Vitest | 2.1.8 | Testing |

---

## Getting Started

### Prerequisites

- Go 1.22 or higher
- Node.js 18 or higher
- SQLite3

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd velostack-go
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

6. **Run database migrations**
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

8. **Open browser**
   Navigate to `http://localhost:8080`

---

## Project Structure

See [FOLDER.md](FOLDER.md) for detailed directory structure.

### Quick Reference

```
velostack-go/
├── main.go              # Entry point
├── app/                 # Backend code
│   ├── handlers/        # HTTP handlers
│   ├── services/        # Business logic
│   ├── repositories/    # Database layer
│   ├── models/          # Data models
│   └── middleware/      # Request middleware
├── routes/              # Route definitions
├── frontend/            # Svelte frontend
│   └── src/
│       ├── pages/       # Page components
│       └── components/  # Reusable components
├── migrations/          # DB migrations
├── templates/           # HTML templates
└── storage/             # Uploaded files
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

2. **Middleware** (`app/middlewares/`)
   - `AuthRequired`: Protects routes, requires authenticated user
   - `AdminRequired`: Requires admin role
   - `Guest`: Redirects authenticated users away

3. **Handlers** (`app/handlers/`)
   - Parse HTTP requests
   - Validate input
   - Call appropriate services
   - Return HTTP responses (HTML or JSON for Inertia)

4. **Services** (`app/services/`)
   - Business logic
   - Authentication flows
   - User management
   - Inertia rendering

5. **Repositories** (`app/repositories/`)
   - Database operations
   - SQL queries using Squirrel builder
   - Parameterized queries for security

6. **Models** (`app/models/`)
   - Domain models (User)
   - DTOs for requests/responses

7. **Session** (`app/session/`)
   - Infrastructure layer for session storage
   - Cookie encoding/decoding with securecookie
   - Generic implementation (reusable across projects)
   - Used by services to persist user session data

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
Session created → Inertia redirect to /app
```

---

## Frontend Architecture

### Inertia.js Pattern

VeloStack uses Inertia.js to bridge backend and frontend without building a separate API:

1. **Initial Load**: Server renders HTML via `inertia.html` template
2. **Subsequent Navigation**: Inertia makes XHR requests with `X-Inertia: true` header
3. **Server Response**: Returns JSON with component name and props
4. **Client Rendering**: Svelte dynamically loads the component

### Component Structure

```
frontend/src/
├── main.js              # Inertia initialization
├── app.css              # Global styles (Tailwind)
├── components/          # Reusable UI components
│   ├── Button.svelte
│   └── Input.svelte
└── pages/               # Page components
    ├── Auth/
    │   ├── Login.svelte
    │   └── Register.svelte
    └── App/
        ├── Dashboard.svelte
        └── Profile.svelte
```

### Making Inertia Requests

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

See [FOLDER.md](FOLDER.md) for detailed layer separation explanation.

**Session Configuration**:
- Cookie-based sessions using Gorilla securecookie
- HTTPOnly cookies for security
- SameSite=Lax for CSRF protection
- SecureCookie encoding with secret key

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
// Protected route
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
| GET | `/login` | `AuthHandler.ShowLoginForm` | Login form |
| POST | `/login/login` | `AuthHandler.Login` | Login submission |
| POST | `/login/register` | `AuthHandler.Register` | Registration |
| GET | `/auth/google` | `AuthHandler.GoogleLogin` | OAuth start |
| GET | `/auth/google/callback` | `AuthHandler.GoogleCallback` | OAuth callback |
| POST | `/logout` | `AuthHandler.Logout` | Logout |
| GET | `/api/me` | `AuthHandler.Me` | Current user API |

### Protected App Routes

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/app` | `AppHandler.Dashboard` | Dashboard |
| GET | `/app/profile` | `AppHandler.Profile` | Profile page |
| PUT | `/app/profile` | `AppHandler.UpdateProfile` | Update profile |
| POST | `/upload` | `UploadHandler.Upload` | File upload |

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
   ```bash
   go run main.go
   ```
   - Reads `.vite-port` for dev server URL
   - Proxies requests to Vite

3. **Auto-reload**:
   - Go 1.22+ automatically reloads on code changes
   - Vite HMR updates frontend instantly

### Available Scripts

```bash
# Development
npm run dev          # Start Vite dev server
go run main.go       # Start Go server

# Production build
npm run build        # Build frontend
go build -o velostack-go .  # Build Go binary

# Database
goose -dir migrations sqlite3:data/app.db up  # Run migrations

# Testing
npm run test         # Run frontend tests
```

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
   go build -o velostack-go .
   ```

3. **Run migrations**:
   ```bash
   goose -dir migrations sqlite3:data/app.db up
   ```

4. **Start server**:
   ```bash
   ./velostack-go
   ```

### Environment Configuration

Set production environment variables:

```bash
APP_ENV=production
APP_PORT=8080
DB_PATH=data/app.db
# ... other variables
```

### Process Management

Use a process manager like PM2, systemd, or Supervisor:

```ini
# /etc/systemd/system/velostack-go.service
[Unit]
Description=VeloStack Go
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/velostack-go
ExecStart=/path/to/velostack-go/velostack-go
Restart=always

[Install]
WantedBy=multi-user.target
```

---

## Security

### Implemented Security Measures

| Feature | Implementation |
|---------|----------------|
| Password Hashing | bcrypt with cost factor |
| SQL Injection | Squirrel parameterized queries |
| XSS Protection | HTML escaping in templates |
| CSRF Protection | SameSite cookies, OAuth state |
| Session Security | HTTPOnly, SecureCookie encoding |
| File Upload | Type validation, size limits (5MB) |
| Role-Based Access | Admin/User middleware guards |

### Best Practices

1. **Never commit `.env`** - Contains secrets
2. **Use HTTPS in production** - Set `APP_ENV=production`
3. **Rotate session secrets** - Change `SESSION_SECRET` regularly
4. **Validate all input** - File uploads, form data
5. **Rate limiting** - Consider adding for auth endpoints

---

## Environment Variables

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `APP_ENV` | Environment | `development`, `production` |
| `APP_PORT` | Server port | `8080` |
| `APP_URL` | Application URL | `http://localhost:8080` |
| `DB_PATH` | Database path | `data/app.db` |
| `SESSION_SECRET` | Session encryption key | `your-secret-key` |
| `SESSION_NAME` | Session cookie name | `velostack_session` |

### OAuth Variables (Optional)

| Variable | Description | Example |
|----------|-------------|---------|
| `GOOGLE_CLIENT_ID` | Google OAuth client ID | `xxx.apps.googleusercontent.com` |
| `GOOGLE_CLIENT_SECRET` | Google OAuth secret | `GOCSPX-xxx` |
| `GOOGLE_REDIRECT_URL` | OAuth callback URL | `http://localhost:8080/auth/google/callback` |

### Example `.env`

```bash
# Application
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080

# Database
DB_PATH=data/app.db

# Session
SESSION_SECRET=your-32-character-secret-key
SESSION_NAME=velostack_session

# Google OAuth (optional)
GOOGLE_CLIENT_ID=xxx.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-xxx
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
```

---

## API Reference

### Authentication Endpoints

#### POST `/login/login`

Login with email/password.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (Success):**
```json
{
  "redirect": "/app"
}
```

#### POST `/login/register`

Register new user.

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "password_confirmation": "password123"
}
```

#### GET `/api/me`

Get current authenticated user.

**Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "role": "user",
  "avatar": ""
}
```

### Protected Endpoints

#### PUT `/app/profile`

Update user profile.

**Headers:**
- `X-Inertia: true`

**Request:**
```json
{
  "name": "John Doe Updated",
  "email": "john.updated@example.com"
}
```

#### POST `/upload`

Upload file (avatar).

**Request:** `multipart/form-data` with `file` field.

**Response:**
```json
{
  "file_url": "/storage/avatars/filename.jpg"
}
```

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
goose -dir migrations sqlite3:data/app.db reset
```

---

## Contributing

1. Create feature branch
2. Make changes
3. Run tests
4. Submit pull request

---

## License

[Your License Here]
