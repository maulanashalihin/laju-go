---
name: gortex-setup
description: "Work in the setup area — 50 symbols across 8 files (80% cohesion)"
---

# setup

50 symbols | 8 files | 80% cohesion

## When to Use

Use this skill when working on files in:
- `app/cache/user_cache.go`
- `app/handlers/auth.go`
- `app/handlers/upload.go`
- `app/middlewares/auth.go`
- `app/middlewares/csrf.go`
- `app/middlewares/rate-limit.go`
- `app/session/session.go`
- `routes/web.go`

## Key Files

| File | Symbols |
|------|---------|
| `app/cache/user_cache.go` | Size |
| `app/handlers/auth.go` | GetAvatar, c, Logout, c, Me, ... |
| `app/handlers/upload.go` | Upload, c |
| `app/middlewares/auth.go` | AdminRequired, Guest, AuthRequired, store, closure@52, ... |
| `app/middlewares/csrf.go` | validateToken, constantTimeCompare, Protect, CSRFMiddleware, a, ... |
| `app/middlewares/rate-limit.go` | Limit, closure@55, getClientKey, c |
| `app/session/session.go` | Get, c |
| `routes/web.go` | setupAppRoutes, uploadHandler, setupPublicRoutes, setupAuthRoutes, appHandler, ... |

## Entry Points

- `app/handlers/upload.go::UploadHandler.Upload`
- `app/handlers/auth.go::AuthHandler.GetAvatar`
- `app/handlers/auth.go::AuthHandler.Me`
- `app/middlewares/auth.go::AdminRequired`
- `app/cache/user_cache.go::UserCache.Size`

## Connected Communities

- **app/handlers** (4 cross-edges)
- **services** (2 cross-edges)
- **handlers** (1 cross-edges)
- **queries** (1 cross-edges)
- **queries** (1 cross-edges)
- **middlewares** (1 cross-edges)
- **services** (1 cross-edges)

## How to Explore

```
get_communities with id: "community-47"
smart_context with task: "understand setup", format: "gcx"
find_usages with id: "app/handlers/upload.go::UploadHandler.Upload", format: "gcx"
```

_`format: "gcx"` returns the [GCX1 compact wire format](../../docs/wire-format.md) — round-trippable, ~27% fewer tokens than JSON. Drop it for JSON output; agents using `@gortex/wire` or the Go `github.com/gortexhq/gcx-go` package decode either._
