---
name: gortex-new
description: "Work in the new area — 24 symbols across 7 files (69% cohesion)"
---

# new

24 symbols | 7 files | 69% cohesion

## When to Use

Use this skill when working on files in:
- `app/cache/user_cache.go`
- `app/handlers/app.go`
- `app/handlers/upload.go`
- `app/services/asset.go`
- `app/services/auth.go`
- `app/session/session.go`
- `main.go`

## Key Files

| File | Symbols |
|------|---------|
| `app/cache/user_cache.go` | NewUserCache, ttl |
| `app/handlers/app.go` | NewAppHandler, store, inertiaService, userService |
| `app/handlers/upload.go` | NewUploadHandler, userService, store |
| `app/services/asset.go` | IsDevelopment, GetViteServerURL |
| `app/services/auth.go` | NewAuthService, querier, cfg |
| `app/session/session.go` | New, querier |
| `main.go` | main, setupViteProxy, viteURL, app, runMigrations, ... |

## Entry Points

- `main.go::main`

## Connected Communities

- **queries** (2 cross-edges)
- **services** (1 cross-edges)
- **config** (1 cross-edges)
- **.** (1 cross-edges)
- **services** (1 cross-edges)
- **services** (1 cross-edges)
- **handlers** (1 cross-edges)
- **handlers** (1 cross-edges)

## How to Explore

```
get_communities with id: "community-104"
smart_context with task: "understand new", format: "gcx"
find_usages with id: "main.go::main", format: "gcx"
```

_`format: "gcx"` returns the [GCX1 compact wire format](../../docs/wire-format.md) — round-trippable, ~27% fewer tokens than JSON. Drop it for JSON output; agents using `@gortex/wire` or the Go `github.com/gortexhq/gcx-go` package decode either._
