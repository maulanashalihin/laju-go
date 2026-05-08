---
name: gortex-queries
description: "Work in the queries area — 9 symbols across 1 files (93% cohesion)"
---

# queries

9 symbols | 1 files | 93% cohesion

## When to Use

Use this skill when working on files in:
- `app/queries/user.sql.go`

## Key Files

| File | Symbols |
|------|---------|
| `app/queries/user.sql.go` | CreateUserWithGoogleIDParams, Email, UpdatedAt, GoogleID, Name, ... |

## How to Explore

```
get_communities with id: "community-112"
smart_context with task: "understand queries", format: "gcx"
```

_`format: "gcx"` returns the [GCX1 compact wire format](../../docs/wire-format.md) — round-trippable, ~27% fewer tokens than JSON. Drop it for JSON output; agents using `@gortex/wire` or the Go `github.com/gortexhq/gcx-go` package decode either._
