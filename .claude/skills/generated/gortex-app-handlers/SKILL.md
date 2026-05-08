---
name: gortex-app-handlers
description: "Work in the app/handlers area — 28 symbols across 6 files (71% cohesion)"
---

# app/handlers

28 symbols | 6 files | 71% cohesion

## When to Use

Use this skill when working on files in:
- `app/handlers/auth.go`
- `app/middlewares/auth.go`
- `app/middlewares/csrf.go`
- `app/services/inertia.go`
- `app/session/session.go`
- `main.go`

## Key Files

| File | Symbols |
|------|---------|
| `app/handlers/auth.go` | Login, c, GoogleCallback, c |
| `app/middlewares/auth.go` | CORS, closure@104 |
| `app/middlewares/csrf.go` | setToken, generateToken, c |
| `app/services/inertia.go` | RenderWithMeta, renderJSON, component, props, props, ... |
| `app/session/session.go` | Set, value, key, Save |
| `main.go` | Render, component, c |

## Entry Points

- `app/handlers/auth.go::AuthHandler.GoogleCallback`
- `app/services/inertia.go::InertiaService.RenderWithMeta`
- `app/handlers/auth.go::AuthHandler.Login`
- `app/middlewares/auth.go::CORS`
- `main.go::Render`

## Connected Communities

- **setup** (5 cross-edges)
- **handlers** (2 cross-edges)
- **handlers** (2 cross-edges)
- **session** (1 cross-edges)
- **queries** (1 cross-edges)
- **queries** (1 cross-edges)
- **middlewares** (1 cross-edges)
- **queries** (1 cross-edges)

## How to Explore

```
get_communities with id: "community-29"
smart_context with task: "understand app/handlers", format: "gcx"
find_usages with id: "app/handlers/auth.go::AuthHandler.GoogleCallback", format: "gcx"
```

_`format: "gcx"` returns the [GCX1 compact wire format](../../docs/wire-format.md) — round-trippable, ~27% fewer tokens than JSON. Drop it for JSON output; agents using `@gortex/wire` or the Go `github.com/gortexhq/gcx-go` package decode either._
