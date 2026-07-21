---
type: source
title: "Observation: Docs audit complete — 18 files reviewed, 10 fixed"
slug: obs-2026-07-06-docs-audit-complete-18-files-reviewed-10-fixed
status: observation
created: 2026-07-06
updated: 2026-07-06
relevance: high
observed_at: 2026-07-06T01:45:02.433Z
tags: ["docs", "audit", "laju-go", "cleanup"]
source_context: "Review docs/ at user's request after seeing capture results to llm-wiki"
---
# ⭐ Observation: Docs audit complete — 18 files reviewed, 10 fixed
Audit of docs/ against the actual Laju Go codebase. Findings: docs/README.md references 20+ files that don't exist; docs/guide/email.md MailerService signature differs (no SendTemplate/SendWelcomeEmail, all inline HTML); docs/guide/architecture.md & handlers.md have outdated constructor signatures; docs/guide/templ.md LandingPage & InertiaPage signatures outdated; docs/guide/routing.md route examples don't match; docs/guide/storage.md upload handler code outdated; docs/guide/data-protection.md mentions BackupService that doesn't exist; templates/index.templ says "zero CGO via modernc" but actually uses mattn/go-sqlite3 (CGO). All have been fixed.
*Relevance: high*

*Context: Review docs/ at user's request after seeing capture results to llm-wiki*

*Tags: docs audit laju-go cleanup*
---
*Observed: 2026-07-06T01:45:02.433Z*