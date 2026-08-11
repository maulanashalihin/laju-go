# create-laju-go

CLI tool to quickly create a new Laju Go project.

## Usage

```bash
# Create new project
npx create-laju-go my-app

# Navigate to project
cd my-app

# Start development server
npm run dev:all
```

## Options

```bash
npx create-laju-go <project-name> [options]
```

| Option | Description | Default |
|--------|-------------|---------|
| `--package-manager` | Package manager to use (npm, yarn, bun) | Auto-detect |

## Prerequisites

Before using this installer, make sure you have:

- **Go 1.26+** installed and in PATH
- **Node.js 18+** installed
- **Git** installed and in PATH

## What's Included

The created project includes:

- **Go Fiber** - High-performance web framework (fasthttp)
- **Svelte 5** - Reactive UI framework
- **Inertia.js** - Server-driven SPA without client-side routing
- **TailwindCSS v4** - Utility-first CSS framework
- **SQLite** - Lightweight database with sqlc (compile-time type-safe SQL)
- **Authentication** - Email/Password + Google OAuth
- **Session Management** - Database-backed sessions
- **Role-Based Access** - Admin/User roles with middleware
- **File Upload** - Avatar upload with validation
- **Goose Migrations** - Database schema management

## Next Steps

After creating the project:

```bash
# Navigate to project
cd my-app

# Install Air for hot reload (recommended)
go install github.com/air-verse/air@latest

# Start development (runs both Vite and Go server)
npm run dev:all
```

Visit `http://localhost:8080` to view the application.

## Project Structure

```
my-app/
├── cmd/laju-go/main.go       # Go entry point
├── go.mod                    # Go dependencies
├── package.json              # Node.js dependencies
├── vite.config.js            # Vite configuration
├── .env.example              # Environment template
│
├── app/                      # Go backend code
│   ├── config/               # Environment configuration
│   ├── models/               # Data structures + DTOs
│   ├── queries/              # sqlc-generated query code
│   ├── services/             # Business logic layer
│   ├── handlers/             # HTTP request handlers
│   ├── middlewares/          # Auth, CSRF, rate limiting
│   ├── cache/                # In-memory session cache
│   └── session/              # Session store (SQLite + cache)
│
├── frontend/                 # Svelte 5 frontend
│   └── src/
│       ├── components/       # Reusable UI components
│       ├── pages/            # Page components
│       └── main.ts           # Entry point
│
├── routes/                   # Route definitions
├── templates/                # templ HTML components
├── migrations/               # Goose SQL migrations
├── data/                     # SQLite database (gitignored)
└── storage/                  # User uploads (gitignored)
```

## Development Commands

| Command | Description |
|---------|-------------|
| `npm run dev` | Start Vite dev server (frontend HMR) |
| `npm run dev:go` | Start Go server with Air (hot reload) |
| `npm run dev:all` | Run both Vite and Go server concurrently |
| `npm run build` | Build frontend assets and Go binary |
| `npm run serve` | Run the production binary |

## Documentation

For complete Laju Go documentation:

**[Laju Go Documentation →](https://github.com/maulanashalihin/laju-go)**

## Troubleshooting

### Go not found
Make sure Go is installed and added to your PATH:
```bash
go version
```

### Permission denied on Linux/Mac
```bash
chmod +x node_modules/.bin/*
```

### Database issues
Delete the database files and restart:
```bash
rm -rf data/*.db data/*.db-shm data/*.db-wal
npm run dev:all
```

## License

MIT
