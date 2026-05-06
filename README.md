# Laju Go

High-performance SaaS boilerplate built with **Go Fiber** + **Svelte 5** + **Inertia.js 3** + **SQLite**.

Build production-ready web applications faster with a clean, layered architecture that combines the speed of Go with the developer experience of modern frontend frameworks.

## 🚀 Quick Start

### Option 1: Using create-laju-go CLI (Recommended)

```bash
# Create new project with CLI
npx create-laju-go my-app

# Navigate to project
cd my-app

# Start development
npm run dev:all
```

### Option 2: Clone Repository

```bash
# Clone the repository
git clone https://github.com/maulanashalihin/laju-go.git
cd laju-go

# Install dependencies
go mod download && npm install

# Set up environment
cp .env.example .env

# Start development servers
npm run dev:all
```

Visit `http://localhost:8080` to see your application running.

## ✨ Features

### Authentication & Security
- **Email/Password Authentication** - Secure login with bcrypt password hashing
- **Google OAuth 2.0** - One-click social login integration
- **Password Reset** - Email-based password recovery with secure tokens
- **Session Management** - Database-backed persistent sessions
- **CSRF Protection** - Built-in cross-site request forgery prevention
- **Rate Limiting** - Configurable request throttling for sensitive endpoints

### User Management
- **Role-Based Access Control** - Admin/User roles with middleware guards
- **Profile Management** - Update profile, change password, avatar upload
- **File Upload** - Avatar upload with validation and secure storage

### Development Experience
- **Hot Module Replacement** - Vite HMR for instant frontend updates
- **Go Hot Reload** - Air automatically rebuilds on Go file changes
- **Clean Architecture** - Separated layers (handlers, services, repositories)
- **TypeScript Ready** - Full type safety in frontend code

### Production Ready
- **SQLite Optimized** - WAL mode, connection pooling, production-tuned
- **Database Migrations** - Goose-based schema version control
- **Docker Support** - Multi-stage builds for efficient containerization
- **Systemd Ready** - Production deployment with process management

## 📚 Documentation

| Section | Description |
|---------|-------------|
| [Getting Started](docs/getting-started/introduction.md) | Introduction, installation, and configuration |
| [Architecture Guide](docs/guide/architecture.md) | Layered architecture, design patterns, and best practices |
| [Routing & Handlers](docs/guide/routing.md) | Route definitions, middleware, and request handling |
| [Database](docs/guide/database.md) | SQLite setup, migrations, and query building |
| [Authentication](docs/guide/authentication.md) | Auth flows, OAuth, sessions, and password reset |
| [Frontend](docs/guide/frontend.md) | Svelte 5 components and Inertia.js integration |
| [Deployment](docs/deployment/development.md) | Development workflow, production deployment, Docker |
| [API Reference](docs/reference/api-reference.md) | Complete endpoint documentation |
| [Troubleshooting](docs/reference/troubleshooting.md) | Common issues and solutions |

## 📁 Project Structure

```
laju-go/
├── main.go                    # Application entry point
├── app/                       # Backend Go code
│   ├── handlers/              # HTTP request handlers
│   ├── services/              # Business logic layer
│   ├── queries/               # Generated SQL query code (sqlc)
│   ├── middlewares/           # Request middleware
│   └── models/                # Data structures
├── frontend/                  # Svelte 5 frontend
│   └── src/
│       ├── components/        # Reusable UI components
│       ├── pages/             # Page components
│       └── lib/               # Utilities and helpers
├── queries/                   # SQL query source files (write queries here)
├── routes/                    # Route definitions
├── migrations/                # Database migrations
├── templates/                 # Templ templates (HTML + Go typed components)
└── docs/                      # Documentation
```

> 📖 See [Project Structure](docs/reference/project-structure.md) for a complete directory reference.

## 🛠️ Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| **Backend** | Go 1.26+ | Programming language |
| **Web Framework** | Fiber v2 | High-performance HTTP framework (fasthttp) |
| **Database** | SQLite3 | Embedded SQL database |
| **Query Builder** | sqlc | Compile-time type-safe SQL code generation |
| **Migrations** | Goose | Database schema management |
| **Frontend** | Svelte 5 | Reactive UI framework |
| **Build Tool** | Vite 5 | Fast build tooling and dev server |
| **Styling** | Tailwind CSS 4 | Utility-first CSS framework |
| **Templating** | templ | Type-safe HTML components for Go |
| **SPA Bridge** | Inertia.js 3 | Server-driven single-page apps |
| **Icons** | Lucide Svelte | Beautiful, consistent icons |

### Why SQLite (`modernc.org/sqlite`)?

We intentionally chose `modernc.org/sqlite` (pure Go) over `mattn/go-sqlite3` (CGO). Here's why this decision is locked in:

| Factor | `modernc.org/sqlite` ✅ | `mattn/go-sqlite3` ❌ |
|--------|------------------------|----------------------|
| **Cross-compile** | `GOOS=linux GOARCH=amd64 go build` — just works | Needs Docker, musl-cross, or server-side GCC |
| **Static binary** | Single self-contained binary | Links to `libsqlite3`, dynamic dependency hell |
| **Docker/CI** | `FROM golang:alpine` works | Must install `gcc`, `libsqlite3-dev`, image bloat |
| **Debug production** | Full Go stack traces | CGO stack traces are opaque and painful |
| **Raw DB speed** | ~20-50% slower in benchmarks | Faster at micro-benchmark level |

**The catch:** At the full HTTP stack level (Fiber routing + JSON marshal + auth + network), the SQLite driver difference is **less than 1.5%** of total request latency. You're bottlenecked by JSON/auth/network long before the driver. The deployment simplicity of pure Go wins every time.

**Decision is final.** Don't migrate to `mattn/go-sqlite3` unless you have a very specific reason (e.g. need SQLite extensions, or doing batch ETL where DB is 90% of CPU).

## 📦 Installation

### Prerequisites

- **Go** 1.26 or higher
- **Node.js** 18 or higher
- **SQLite3** (usually pre-installed on macOS/Linux)
- **Git** for version control

### Method 1: Using create-laju-go CLI (Recommended)

The easiest way to create a new Laju Go project:

```bash
# Create new project
npx create-laju-go my-app

# Navigate to project
cd my-app

# Install Air for hot reload (recommended)
go install github.com/air-verse/air@latest

# Start development
npm run dev:all
```

The CLI will:
- Check for Go and Git installation
- Let you choose package manager (npm, yarn, bun)
- Clone the template from GitHub
- Install all dependencies
- Set up environment configuration

### Method 2: Manual Installation

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

4. **Configure environment**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` with your settings. At minimum, set:
   ```bash
   APP_ENV=development
   SESSION_SECRET=your-32-character-secret-key
   ```

5. **Set up Google OAuth (Optional)**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project and enable Google+ API
   - Create OAuth 2.0 credentials
   - Add `http://localhost:8080/auth/google/callback` to authorized redirect URIs
   - Copy Client ID and Secret to `.env`

6. **Set up Email/SMTP (Optional - for password reset)**
   - Configure SMTP settings in `.env`
   - For Gmail, use an [App Password](https://support.google.com/accounts/answer/185833)

## 🏃 Development

### Option 1: Run Everything Together (Recommended)

Start both Vite and Go servers with hot reload:

```bash
npm run dev:all
```

### Option 2: Run Servers Separately

**Terminal 1** - Vite dev server (frontend HMR):
```bash
npm run dev
```

**Terminal 2** - Go server with hot reload:
```bash
air
# Or via npm
npm run dev:go
```

### Option 3: Manual Run

```bash
# Go server (manual restart after changes)
go run main.go

# Vite dev server
npm run dev
```

### Available Scripts

```bash
# Development
npm run dev          # Start Vite dev server
npm run dev:go       # Start Go server with Air hot reload
npm run dev:all      # Run both Vite and Air concurrently

# Production
npm run build        # Build frontend and Go binary
npm run serve        # Run production binary

# Testing
npm run test:run     # Run frontend tests
```

### Development Workflow

| You Edit | What Happens |
|----------|--------------|
| `.svelte` files | Vite HMR updates instantly |
| `.go` files | Air rebuilds and restarts (~1-2 sec) |
| `.css` files | Hot reload (instant) |
| `migrations/` | Auto-run on server start |

## 🚀 Production Deployment

### Quick Deploy

```bash
# Build frontend assets
npm run build

# Build Go binary
go build -o laju-go .

# Run the server
./laju-go
```

### Docker Deployment

```bash
# Build the image
docker build -t laju-go .

# Run the container
docker run -p 8080:8080 \
  -v $(pwd)/data:/root/data \
  -v $(pwd)/storage:/root/storage \
  laju-go
```

### Ubuntu/Debian Server

For complete production deployment instructions including systemd service setup, Nginx reverse proxy, and SSL configuration, see [Production Deployment Guide](docs/deployment/production.md).

## 🔐 Default Admin Setup

After your first registration, promote your user to admin via SQLite:

```bash
sqlite3 data/app.db "UPDATE users SET role = 'admin' WHERE email = 'your@email.com';"
```

## 🗄️ Database Migrations

Migrations run automatically on startup. Manual commands:

```bash
# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run all migrations
goose -dir migrations sqlite3 data/app.db up

# Check migration status
goose -dir migrations sqlite3 data/app.db status

# Rollback last migration
goose -dir migrations sqlite3 data/app.db down
```

## 📝 SQL Queries (sqlc)

This project uses [sqlc](https://sqlc.dev/) for compile-time type-safe SQL queries. Write your SQL in `queries/*.sql`, then generate Go code:

```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate Go code from SQL files
npm run db:generate
```

### Directory Structure

| Directory | Purpose |
|-----------|---------|
| `queries/` | SQL source files — **write your queries here** |
| `app/queries/` | Generated Go code + wrapper — **do not edit manually** |

### Adding a New Query

1. Add the query to `queries/user.sql` or create a new `.sql` file:
```sql
-- name: GetUserCount :one
SELECT COUNT(*) FROM users;
```

2. Regenerate:
```bash
npm run db:generate
```

3. Use in your service:
```go
count, err := s.querier.GetUserCount(ctx)
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage report
go test -cover ./...
```

## 📊 Performance Optimizations

### SQLite Production Settings

The application includes these optimizations by default, tuned for **Vultr High Frequency 1-2GB RAM**:

| Setting | Value | Benefit |
|---------|-------|---------|
| `journal_mode` | WAL | Better write concurrency |
| `synchronous` | NORMAL | Faster writes with safety |
| `cache_size` | 16MB | Reduced disk I/O (optimized for 1-2GB RAM) |
| `mmap_size` | 256MB | NVMe memory-mapped I/O |
| `temp_store` | MEMORY | Faster temp table operations |
| `busy_timeout` | 5000ms | Automatic retry on locks |
| Connection Pool | 15 max | Efficient connection reuse |

### Tune for Your Server

Different RAM size? See the complete **[SQLite Configuration Guide](docs/deployment/sqlite-configuration.md)** for optimal settings:

| Server RAM | MaxOpenConns | cache_size | mmap_size |
|------------|--------------|------------|-----------|
| 512MB ⚠️ | 10 | 8MB | 128MB |
| **1-2GB ✅** | **15** | **16MB** | **256MB** |
| 4GB | 25 | 32MB | 512MB |
| 8GB | 50 | 256MB | 1GB |
| 16GB+ | 100 | 500MB+ | 2GB |

> 📖 **Full guide**: [SQLite Configuration Guide](docs/deployment/sqlite-configuration.md) - Complete reference for tuning SQLite based on RAM, CPU, storage type, and workload patterns.

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- [Go Fiber](https://gofiber.io/) - Fast web framework
- [Svelte](https://svelte.dev/) - Cybernetically enhanced web apps
- [Inertia.js](https://inertiajs.com/) - Server-driven SPA
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS
- [Lucide Icons](https://lucide.dev/) - Beautiful, consistent icons

## 📞 Support

- **Documentation**: [docs/](docs/) folder
- **Issues**: [GitHub Issues](https://github.com/maulanashalihin/laju-go/issues)
- **Discussions**: [GitHub Discussions](https://github.com/maulanashalihin/laju-go/discussions)
