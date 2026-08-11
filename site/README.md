# Laju Go site

The Laju Go landing page + documentation + blog, built with **Astro + Starlight**
(fully static — no server adapter).

## Develop

```bash
cd site
npm install
npm run dev        # http://localhost:4321
```

## Build

```bash
npm run build      # static output → site/dist
```

## Deploy (Cloudflare Pages)

- **Dashboard**: root directory `site`, build `npm run build`, output `dist`.
- **CLI**: `npx wrangler pages deploy dist --project-name laju-go`
  (`site/wrangler.toml` pins the config).

Set the production domain in `astro.config.mjs` (`site`) before going live.

## Content

- Landing page: `src/pages/index.astro` (custom dark-tech design, outside Starlight).
- Docs: `src/content/docs/**/*.mdx` — sidebar configured in `astro.config.mjs`.
- Blog: `src/content/docs/blog/*.mdx` — powered by starlight-blog plugin.
- Theme overrides: `src/styles/custom.css` (cyan-teal accent, navy-tinted dark mode).

Keep docs in sync with code changes in the same PR.
