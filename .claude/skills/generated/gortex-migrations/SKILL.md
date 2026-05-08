---
name: gortex-migrations
description: "Work in the migrations area — 11 symbols across 1 files (100% cohesion)"
---

# migrations

11 symbols | 1 files | 100% cohesion

## When to Use

Use this skill when working on files in:
- `migrations/0001_create_users_table.sql`

## Key Files

| File | Symbols |
|------|---------|
| `migrations/0001_create_users_table.sql` | users, password, updated_at, name, email_verified, ... |

## How to Explore

```
get_communities with id: "community-103"
smart_context with task: "understand migrations", format: "gcx"
```

_`format: "gcx"` returns the [GCX1 compact wire format](../../docs/wire-format.md) — round-trippable, ~27% fewer tokens than JSON. Drop it for JSON output; agents using `@gortex/wire` or the Go `github.com/gortexhq/gcx-go` package decode either._
