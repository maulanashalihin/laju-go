---
name: gortex-templates
description: "Work in the templates area — 8 symbols across 1 files (90% cohesion)"
---

# templates

8 symbols | 1 files | 90% cohesion

## When to Use

Use this skill when working on files in:
- `templates/inertia_templ.go`

## Key Files

| File | Symbols |
|------|---------|
| `templates/inertia_templ.go` | InertiaPage, isDev, title, styles, closure@15, ... |

## Entry Points

- `templates/inertia_templ.go::InertiaPage`

## Connected Communities

- **handlers** (1 cross-edges)

## How to Explore

```
get_communities with id: "community-74"
smart_context with task: "understand templates", format: "gcx"
find_usages with id: "templates/inertia_templ.go::InertiaPage", format: "gcx"
```

_`format: "gcx"` returns the [GCX1 compact wire format](../../docs/wire-format.md) — round-trippable, ~27% fewer tokens than JSON. Drop it for JSON output; agents using `@gortex/wire` or the Go `github.com/gortexhq/gcx-go` package decode either._
