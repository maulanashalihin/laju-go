---
name: gortex-config
description: "Work in the config area — 19 symbols across 1 files (94% cohesion)"
---

# config

19 symbols | 1 files | 94% cohesion

## When to Use

Use this skill when working on files in:
- `app/config/config.go`

## Key Files

| File | Symbols |
|------|---------|
| `app/config/config.go` | Config, SMTPPass, DBPath, IsDevelopment, FrontendURL, ... |

## How to Explore

```
get_communities with id: "community-95"
smart_context with task: "understand config", format: "gcx"
```

_`format: "gcx"` returns the [GCX1 compact wire format](../../docs/wire-format.md) — round-trippable, ~27% fewer tokens than JSON. Drop it for JSON output; agents using `@gortex/wire` or the Go `github.com/gortexhq/gcx-go` package decode either._
