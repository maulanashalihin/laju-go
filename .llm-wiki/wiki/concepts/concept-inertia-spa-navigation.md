---
type: concept
status: stub
---

# Inertia SPA Navigation

Inertia.js enables single-page application (SPA) navigation in Laju Go without building a separate API.

## Navigation Types

| Action | Method | Notes |
|--------|--------|-------|
| Internal links | `<a href="/path" use:inertia>` | Uses `use:inertia` action from `@inertiajs/svelte` |
| OAuth links | `<a href="/auth/google">` | Plain `<a>` — must redirect to external provider |
| Form submission | `router.post()` / `router.put()` | From `@inertiajs/svelte` |
| File upload | `fetch()` + FormData | Must include CSRF header, then `router.put()` to save URL |

## Redirect Rules

- POST/PUT handlers must use `c.Redirect(path, fiber.StatusSeeOther)` (303)
- Inertia does not follow 302 correctly for form submissions — 303 changes POST/PUT to GET on redirect

## Source

Captured from [[sources/SRC-2026-07-06-001]].
