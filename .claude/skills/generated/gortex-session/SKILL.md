---
name: gortex-session
description: "Work in the session area — 8 symbols across 1 files (64% cohesion)"
---

# session

8 symbols | 1 files | 64% cohesion

## When to Use

Use this skill when working on files in:
- `app/session/session.go`

## Key Files

| File | Symbols |
|------|---------|
| `app/session/session.go` | Session, userID, dirty, id, expiresAt, ... |

## How to Explore

```
get_communities with id: "community-53"
smart_context with task: "understand session", format: "gcx"
```

_`format: "gcx"` returns the [GCX1 compact wire format](../../docs/wire-format.md) — round-trippable, ~27% fewer tokens than JSON. Drop it for JSON output; agents using `@gortex/wire` or the Go `github.com/gortexhq/gcx-go` package decode either._
