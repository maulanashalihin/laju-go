---
name: gortex-handlers
description: "Work in the handlers area — 29 symbols across 7 files (71% cohesion)"
---

# handlers

29 symbols | 7 files | 71% cohesion

## When to Use

Use this skill when working on files in:
- `app/handlers/app.go`
- `app/handlers/auth.go`
- `app/handlers/password-reset.go`
- `app/handlers/public.go`
- `app/services/inertia.go`
- `app/services/mailer.go`
- `app/session/session.go`

## Key Files

| File | Symbols |
|------|---------|
| `app/handlers/app.go` | UpdatePassword, c |
| `app/handlers/auth.go` | ShowLoginForm, ShowRegisterForm, c, c |
| `app/handlers/password-reset.go` | ShowResetPasswordForm, ResetPassword, c, SendResetLink, c, ... |
| `app/handlers/public.go` | About, c |
| `app/services/inertia.go` | Render, c, component, props |
| `app/services/mailer.go` | ValidateResetToken, token, InvalidateResetToken, token |
| `app/session/session.go` | GetFlash, c, key |

## Entry Points

- `app/handlers/password-reset.go::PasswordResetHandler.ResetPassword`
- `app/handlers/password-reset.go::PasswordResetHandler.SendResetLink`
- `app/handlers/password-reset.go::PasswordResetHandler.ShowResetPasswordForm`
- `app/handlers/app.go::AppHandler.UpdatePassword`

## Connected Communities

- **handlers** (2 cross-edges)
- **services** (1 cross-edges)
- **setup** (1 cross-edges)
- **app/handlers** (1 cross-edges)
- **services** (1 cross-edges)
- **services** (1 cross-edges)
- **queries** (1 cross-edges)

## How to Explore

```
get_communities with id: "community-8"
smart_context with task: "understand handlers", format: "gcx"
find_usages with id: "app/handlers/password-reset.go::PasswordResetHandler.ResetPassword", format: "gcx"
```

_`format: "gcx"` returns the [GCX1 compact wire format](../../docs/wire-format.md) — round-trippable, ~27% fewer tokens than JSON. Drop it for JSON output; agents using `@gortex/wire` or the Go `github.com/gortexhq/gcx-go` package decode either._
