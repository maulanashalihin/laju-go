---
name: gortex-routes
description: "Work in the routes area — 8 symbols across 1 files (71% cohesion)"
---

# routes

8 symbols | 1 files | 71% cohesion

## When to Use

Use this skill when working on files in:
- `routes/web.go`

## Key Files

| File | Symbols |
|------|---------|
| `routes/web.go` | setupStaticRoutes, app, SetupRoutes, csrfMiddleware, mailerService, ... |

## Connected Communities

- **setup** (3 cross-edges)

## How to Explore

```
get_communities with id: "community-18"
smart_context with task: "understand routes", format: "gcx"
```

_`format: "gcx"` returns the [GCX1 compact wire format](../../docs/wire-format.md) — round-trippable, ~27% fewer tokens than JSON. Drop it for JSON output; agents using `@gortex/wire` or the Go `github.com/gortexhq/gcx-go` package decode either._
