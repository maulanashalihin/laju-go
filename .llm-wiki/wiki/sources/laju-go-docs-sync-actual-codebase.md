---
type: source
title: "Laju Go docs synced with actual codebase"
slug: laju-go-docs-sync-actual-codebase
status: insight
created: 2026-07-06
updated: 2026-07-06
category: maintenance
---
# Laju Go docs synced with actual codebase
Audit of Laju Go docs/ found 18 doc files, 10 of which needed fixes because they didn't reflect the actual codebase. Fixes include: [[concept-cgo-cross-compilation]] (fix "zero CGO" claim to mattn/go-sqlite3 CGO), [[entity-go-fiber]] route signatures, [[entity-sqlite]] MailerService constructor, handler constructor signatures, and references to files that no longer exist in docs/README.md. Landing page template also fixed from modernc claim to mattn which is actually used in production. Git-based deployment (clone → build → systemd restart) added as the primary deployment strategy.
*Category: maintenance*
---
*Captured: 2026-07-06*
## Related
_Add links to related pages._