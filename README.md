# Laju Go

High-performance SaaS boilerplate built with **Go Fiber** + **Inertia.js 3** + **SQLite**.

Build production-ready web applications faster with a clean, layered architecture that combines the speed of Go with the developer experience of modern frontend frameworks. Ships with **Svelte 5** by default, but Inertia.js makes it trivial to swap to **React** or **Vue 3** without changing any Go code.

## 🚀 Quick Start

```bash
git clone https://github.com/maulanashalihin/laju-go.git
cd laju-go
cp .env.example .env
go mod download && npm install
npm run dev:all
```

Visit `http://localhost:8080` to see your application running.

> 👶 **Never set up Go before?** See [Option 1: AI Setup](#-option-1-ai-setup-recommended-for-first-timers) below — copy a prompt, paste to your AI assistant, done.

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
- **Clean Architecture** - Separated layers (handlers, services, queries)
- **TypeScript Ready** - Full type safety in frontend code
- **Framework Agnostic** - Swap Svelte, React, or Vue without touching Go

### Production Ready
- **SQLite Optimized** - WAL mode, connection pooling, production-tuned
- **Database Migrations** - Goose-based schema version control
- **Docker Support** - Multi-stage builds for efficient containerization
- **Systemd Ready** - Production deployment with process management
- **Litestream DR** - Continuous SQLite replication to S3 for disaster recovery

## 📚 Documentation

| Section | Description |
|---------|-------------|
| [Getting Started](docs/getting-started/introduction.md) | Introduction, installation, and configuration |
| [Architecture Guide](docs/guide/architecture.md) | Layered architecture, design patterns, and best practices |
| [Routing & Handlers](docs/guide/routing.md) | Route definitions, middleware, and request handling |
| [Database](docs/guide/database.md) | SQLite setup, migrations, and query building |
| [Authentication](docs/guide/authentication.md) | Auth flows, OAuth, sessions, and password reset |
| [Frontend](docs/guide/frontend.md) | Svelte 5 (default), React & Vue support via Inertia.js |
| [Deployment](docs/deployment/development.md) | Development workflow, production deployment, Docker, Litestream DR |
| **[Benchmark](docs/benchmark/)** | **SQLite driver performance across Vultr servers** |
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
├── frontend/                  # Svelte 5 frontend (swappable to React/Vue)
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
| **Frontend** | Svelte 5 (default) | Reactive UI framework — swap to React or Vue via Inertia.js |
| **Build Tool** | Vite 5 | Fast build tooling and dev server |
| **Styling** | Tailwind CSS 4 | Utility-first CSS framework |
| **Templating** | templ | Type-safe HTML components for Go |
| **SPA Bridge** | Inertia.js 3 | Server-driven single-page apps |
| **Icons** | Lucide Svelte | Beautiful, consistent icons |

### Why SQLite (`modernc.org/sqlite`)?

We intentionally chose `modernc.org/sqlite` (pure Go) over `mattn/go-sqlite3` (CGO). Here's why — backed by [real benchmarks](https://github.com/maulanashalihin/go-sqlite-benchmark-mattn-vs-modernc).

| Factor | `modernc.org/sqlite` ✅ | `mattn/go-sqlite3` ❌ |
|--------|------------------------|----------------------|
| **Cross-compile** | `GOOS=linux GOARCH=amd64 go build` — just works | Needs Docker, musl-cross, or server-side GCC |
| **Static binary** | Single self-contained binary | Links to `libsqlite3`, dynamic dependency hell |
| **Docker/CI** | `FROM golang:alpine` works | Must install `gcc`, `libsqlite3-dev`, image bloat |
| **Debug production** | Full Go stack traces | CGO stack traces are opaque and painful |
| **Restart safety** | No stale CGO state — clean restarts every time | Stale C threads can cause crashes after restart |

### Benchmark Reality (Vultr Servers + Go Fiber + wrk)

Comprehensive benchmark across 3 Vultr server types (shared & dedicated CPU):

| Server | Type | mattn RPS | modernc RPS | Gap |
|--------|------|-----------|-------------|:---:|
| 1 vCPU Shared | Budget | 16,414 | 12,175 | 1.35x |
| 4 vCPU Dedicated | Mid-tier | 84,946 | 63,991 | 1.33x |
| 6 vCPU Shared | High-end | 101,555 | 53,009 | **1.92x** |

**Key findings**:
- **Dedicated CPU** gives 25-81% better per-vCPU performance than shared CPU
- **mattn gap grows on shared CPU** (1.92x) vs dedicated (1.33x) due to Go scheduler contention
- **Server A (6v Shared $96/mo)** offers best RPS/$ for production (100K RPS)
- **modernc scales linearly** up to 4 cores, then drops on shared 6+ cores

### The Real Trade-off

For most SaaS apps, both drivers handle **100K+ RPS** — far beyond what a typical app needs. The practical difference is deployment:

- **modernc**: Single static binary, CI just works, `FROM alpine` or even `scratch`
- **mattn**: Must install `gcc`, configure `CC` env, bloated Docker images, fragile CGO stack traces

> **Bottom line**: Use modernc for development (simpler), switch to mattn for production if you need >50K RPS (2x throughput).

**📋 [Full Benchmark Report →](docs/benchmark/)**

| Document | Description |
|----------|-------------|
| [Strategic Insights](docs/benchmark/sqlite-driver-benchmark-insights-2026-05-08.md) | 10 key findings, decision matrix, cost analysis |
| [All Servers Comparison](docs/benchmark/sqlite-driver-benchmark-all-servers-2026-05-08.md) | Complete data across 3 Vultr servers |
| [Server A Results](docs/benchmark/sqlite-driver-benchmark-2026-05-08.md) | 6v Shared detailed benchmark |
| [Server B vs C](docs/benchmark/sqlite-driver-benchmark-comparison-2026-05-08.md) | 1v Shared vs 4v Dedicated comparison |

## 📦 Installation

### Prerequisites

- **Go** 1.26 or higher
- **Node.js** 18 or higher
- **SQLite3** (usually pre-installed on macOS/Linux)
- **Git** for version control

### ⭐ Option 1: AI Setup (Recommended for first-timers)

No Go installed? No problem. Copy this prompt and paste it to your AI coding assistant (Claude, ChatGPT, Gemini). It will install everything and get the project running.

```text
Set up and run the Laju Go project (https://github.com/maulanashalihin/laju-go) on this machine.

1. Check what OS I'm on (macOS/Linux/Windows) and install prerequisites if missing:
   - Go 1.26+ — install from https://go.dev/dl/ if not found
   - Node.js 18+ — install from https://nodejs.org/ if not found
   - Git — install if not found
   - SQLite3 — macOS/Linux usually have it pre-installed

2. Clone the repo and install dependencies:
   git clone https://github.com/maulanashalihin/laju-go.git
   cd laju-go
   go mod download
   npm install

3. Set up environment:
   cp .env.example .env
   Generate a random 32-character string for SESSION_SECRET in .env

4. Install Go dev tools (Air for hot reload, templ for templates):
   go install github.com/air-verse/air@latest
   go install github.com/a-h/templ/cmd/templ@latest
   Make sure ~/go/bin is in PATH

5. Start the dev server:
   npm run dev:all

6. Confirm it's running by visiting http://localhost:8080

Pause after each step if there are errors. Don't skip steps.
```

### Option 2: Using create-laju-go CLI

```bash
npx create-laju-go my-app
cd my-app
npm run dev:all
```

The CLI will check for Go/Git, choose package manager, clone the template, install dependencies, and set up environment.

### Option 3: Manual Installation

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
npm run build        # Build frontend only (vite build)
npm run build:all    # Full production build: vite + go build
npm run build:linux  # Cross-compile binary for Linux (pure Go, no CGO)
npm run serve        # Run production binary (./laju-go)

# Testing
npm run test:run     # Run frontend tests
```

### Development Workflow

| You Edit | What Happens |
|----------|--------------|
| `.svelte` / `.tsx` / `.vue` files | Vite HMR updates instantly |
| `.go` files | Air rebuilds and restarts (~1-2 sec) |
| `.css` files | Hot reload (instant) |
| `migrations/` | Auto-run on server start |

## 🚀 Production Deployment

> **Important**: Build everything locally, then upload only the runtime artifacts to your server.
> No build tools (Go, Node, npm) are needed on the server — just the binary and assets.

### 1. Build Locally

```bash
# Full production build (frontend + Go binary)
npm run build:all

# Or for Linux deployment from macOS:
npm run build:linux
```

This produces two things:
- **`laju-go`** — Static Go binary (pure Go SQLite, no CGO)
- **`dist/`** — Frontend assets (CSS/JS built by Vite)

### 2. Deploy Artifacts to Server

Only these files are needed at runtime:

| Artifact | Purpose |
|----------|---------|
| `laju-go` | Go binary (the application) |
| `dist/` | Frontend assets |
| `migrations/` | SQL migrations (auto-run on startup) |
| `.env` | Environment configuration |

```bash
# Example: upload via scp
scp laju-go user@server:/opt/laju-go/
scp -r dist user@server:/opt/laju-go/dist
scp -r migrations user@server:/opt/laju-go/migrations
scp .env user@server:/opt/laju-go/.env
```

### 3. Run with systemd (Production)

Set up a systemd service for auto-start and process management:

```bash
# On the server, create service file
sudo nano /etc/systemd/system/laju-go.service
```

```ini
[Unit]
Description=Laju Go Application
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/laju-go
ExecStart=/opt/laju-go/laju-go
Restart=always
RestartSec=5
EnvironmentFile=/opt/laju-go/.env

[Install]
WantedBy=multi-target.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable laju-go
sudo systemctl start laju-go
```

### One-Click Deploy Script

For automated deployment, configure and use the deploy script:

```bash
cp .deploy.example .deploy
# Edit .deploy with your server details and APP_NAME

npm run deploy
```

This script:
1. ✅ Builds frontend + Go binary **locally**
2. ✅ Uploads only runtime artifacts (binary, `dist/`, `migrations/`) to server
3. ✅ Detects first deploy or update
4. ✅ Sets up `.env` and systemd service (first deploy)
5. ✅ Restarts the service (update)

> See [One-Click Deployment Guide](docs/deployment/one-click-deployment.md) for full instructions.

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

### Ubuntu/Debian Server Setup

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
