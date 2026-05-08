---
name: gortex-middlewares
description: "Work in the middlewares area — 7 symbols across 1 files (94% cohesion)"
---

# middlewares

7 symbols | 1 files | 94% cohesion

## When to Use

Use this skill when working on files in:
- `app/middlewares/rate-limit.go`

## Key Files

| File | Symbols |
|------|---------|
| `app/middlewares/rate-limit.go` | RateLimiterConfig, Message, StatusCode, Window, SkipFailed, ... |

## How to Explore

```
get_communities with id: "community-26"
smart_context with task: "understand middlewares", format: "gcx"
```

_`format: "gcx"` returns the [GCX1 compact wire format](../../docs/wire-format.md) — round-trippable, ~27% fewer tokens than JSON. Drop it for JSON output; agents using `@gortex/wire` or the Go `github.com/gortexhq/gcx-go` package decode either._
