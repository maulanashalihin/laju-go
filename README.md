# Laju Go

High-performance SaaS boilerplate built with **Go Fiber** + **Svelte 5** + **Inertia.js 3** + **SQLite** (CGO).

Build production-ready web applications faster with clean layered architecture that combines the speed of Go with the DX of modern frontend frameworks. Ships with **Svelte 5** by default, but Inertia.js makes it trivial to swap to **React** or **Vue** without changing any Go code.

## 🚀 Quick Start

```bash
git clone https://github.com/maulanashalihin/laju-go.git
cd laju-go
cp .env.example .env
go mod download && npm install
npm run dev:all
```

Visit `http://localhost:8080` to see your application running.

## ✨ Features

### Authentication & Security

- **Email/Password** — Argon2id hashing, session management
- **Google OAuth 2.0** — One-click social login
- **Password Reset** — Email-based recovery with secure tokens
- **CSRF Protection** — Built-in per-session token validation
- **Rate Limiting** — Configurable throttling for auth, API, upload endpoints
- **Session Fixation Protection** — Session ID regenerated on privilege escalation

### User Management

- **Role-Based Access Control** — Admin/User roles with middleware guards
- **Profile Management** — Update name, email, avatar
- **Avatar Upload** — With file type/size validation

### Development Experience

- **Hot Module Replacement** — Vite HMR for instant frontend updates
- **Go Hot Reload** — Air rebuilds on Go file changes (~1-2s)
- **Clean Architecture** — Handler → Service → Query (sqlc-generated)
- **Full TypeScript** — Every `.svelte` file uses `<script lang="ts">`
- **Type-Safe Templates** — Go HTML via [templ](https://templ.guide/)

### Performance & Database

- **SQLite (mattn/go-sqlite3)** — CGO-based, 2x throughput vs pure-Go drivers
- **WAL Mode + mmap** — Optimized for production workloads
- **In-Memory Caching** — Session & user profile TTL caches avoid DB lookups
- **Background Cleanup** — Expired sessions & tokens auto-purged every hour

## 📁 Project Structure

```
laju-go/
├── cmd/laju-go/main.go        # Application entry point
├── app/                       # Backend Go code
│   ├── handlers/              # HTTP request handlers
│   ├── services/              # Business logic layer
│   ├── queries/               # sqlc-generated query code
│   ├── middlewares/           # Auth, CSRF, rate limiting
│   ├── cache/                 # In-memory TTL caches
│   ├── models/                # Data structures + DTOs
│   ├── session/               # Session store (SQLite + cache)
│   └── config/                # Env-based configuration
├── frontend/                  # Svelte 5 frontend
│   └── src/
│       ├── components/        # Header, DarkModeToggle
│       ├── pages/auth/        # Login, Register, ForgotPassword, ResetPassword
│       ├── pages/app/         # Dashboard, Profile
│       └── lib/i18n/          # Internationalization (en/id)
├── queries/                   # SQL source files (write here → sqlc)
├── routes/                    # Route definitions
├── migrations/                # Goose SQL migrations (1 table per file)
├── templates/                 # templ HTML components
├── docs/                      # Documentation
└── systemd/                   # Production service file
```

## 🛠️ Tech Stack

| Layer | Technology |
|-------|------------|
| **Backend** | Go 1.26+, Fiber v2 |
| **Database** | SQLite via `mattn/go-sqlite3` (CGO) |
| **Query Builder** | sqlc — compile-time type-safe SQL |
| **Migrations** | Goose — embedded in binary via `go run` |
| **Frontend** | Svelte 5 (rune-based) |
| **Build Tool** | Vite 8 |
| **Styling** | Tailwind CSS 4 |
| **Templating** | templ — type-safe Go HTML |
| **SPA Bridge** | Inertia.js 3 |
| **Icons** | Lucide Svelte |

### Why `mattn/go-sqlite3` (CGO)?

The `mattn/go-sqlite3` driver delivers **~1.3–1.9x higher throughput** than pure-Go alternatives in real-world benchmarks (100K+ RPS on a $24/mo Vultr instance). The trade-off:

- ✅ **2x faster** on production workloads
- 🛠️ **Cross-compile** via `make build-linux` (requires `brew install zig` for `zig cc`)
- ➡️ **Static binary** still possible — `zig cc` links `libsqlite3` statically

For development (macOS), CGO works out of the box — no extra setup needed.

## ⚡ Quick Reference

```bash
# Development
npm run dev:all                # Vite + Air (hot reload both)

# Build (production)
npm run build:all              # vite build → go build

# Build for Linux (from macOS)
make build-linux               # requires zig cc

# Database
npm run db:migrate             # run pending migrations
npm run db:generate            # sqlc — generate Go from SQL

# Templates
templ generate                 # regenerate templ Go files

# Test
go test ./...                  # backend tests (services, queries, handlers)
# E2E: use pi agent_browser for manual flow testing
```

### Testing Strategy

| Approach | For | Command |
|----------|-----|---------|
| Go unit/integration | Services, queries, handlers | `go test ./...` |
| E2E / user flow | Visual regression, auth flows, form submission | `agent_browser` via pi |

> **E2E testing** dilakukan manual dengan `agent_browser` (buka browser, klik, isi form, verify redirect).
> Tidak perlu Cypress/Playwright — browser asli lebih realistik untuk project skala ini.

## 🚀 Deployment (Your Workflow)

```bash
# 1. Pull latest
git pull

# 2. Build
npm run build:all

# 3. Restart service
sudo systemctl restart laju-go
```

Only runtime artifacts needed on server:

- `laju-go` binary
- `dist/` — frontend assets
- `.env` — configuration
- `migrations/` — auto-run on startup

> **Note**: No Go, Node, or npm needed on the server — just the binary + assets.

## 🗄️ Database

Migrations run **automatically** on startup. Manual:

```bash
# Run pending migrations
go run github.com/pressly/goose/v3/cmd/goose@latest \
  -dir migrations sqlite3 ./data/app.db up

# Status
go run github.com/pressly/goose/v3/cmd/goose@latest \
  -dir migrations sqlite3 ./data/app.db status
```

Write SQL in `queries/*.sql`, then:

```bash
npm run db:generate    # sqlc generates Go code into app/queries/
```

## 📖 Documentation

| Section | Description |
|---------|-------------|
| [Architecture](docs/guide/architecture.md) | Layered design, patterns, conventions |
| [Database](docs/guide/database.md) | SQLite setup, migrations, sqlc |
| [Frontend](docs/guide/frontend.md) | Svelte 5 + Inertia.js patterns |
| [Deployment](docs/deployment/production.md) | Systemd, Nginx, production setup |
| [Benchmarks](docs/benchmark/) | SQLite driver performance data |

## 📄 License

MIT — see [LICENSE](LICENSE).
