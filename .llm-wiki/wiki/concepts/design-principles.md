# Design Principles

> Panduan design untuk agent saat generate halaman baru di project Go Fiber + Svelte 5 + Inertia.js + Tailwind CSS v4.

---

## ⚠️ Anti AI-Slop (BACA INI DULU)

AI punya pola default yang bikin semua output keliatan sama. **Hindari ini sebelum nulis kode.**

### 🎨 Warna — AI Defaults

| ❌ AI Slop | Kenapa | ✅ Ganti dengan |
|-----------|--------|----------------|
| **Purple/blue glow + dark mesh** hero | Tanda paling jelas "AI generated" | Gradient brand (cyan→violet) atau radial gradient lembut |
| **Warm beige/cream bg** (`#f5f1ea`, `#f7f5f1`, `#efeae0`) | LLM default "premium consumer" — semua proyek sama | Pakai token `neutral-*` dari `@theme` |
| **Brass/ochre/oxblood** (`#b08947`, `#9a2436`, `#9c6e2a`) | Sepaket sama beige bg | `brand-400` (cyan) atau `secondary-500` (violet) |
| **Random gradient di setiap section** | Tanpa intent, murahan | Satu gradient di hero aja, sisanya solid |
| **Pure black `#000` untuk dark mode** | Crushing detail, kasar | `neutral-950` (`#070b16`) |
| **Opacity/transparency stacking** 3+ layer | Visual noise, aksesibilitas jelek | Max 2 layer overlay, sisanya solid |

### 🔤 Tipografi — AI Tells

| ❌ AI Slop | Kenapa | ✅ Ganti dengan |
|-----------|--------|----------------|
| **Inter default** | AI selalu fallback ke Inter | OK pakai Inter (sudah di `@theme`), **jangan import ulang** via Google Fonts |
| **Fraunces / Instrument Serif** | Dua font favorit LLM — ketahuan | Jangan dipakai. Serif jarang dibutuhkan |
| **Serif buat "creative" / "premium"** | AI always reach for serif when brief says "creative" | **Default sans-serif.** Serif cuma kalo brand explicit nyebut nama font |
| **Mixed-family emphasis** (sans headline satu kata serif) | Amatir, AI tell klasik | Italic/bold dari font yang SAMA |
| **Headline > 8 words** | Gak bisa dibaca cepat | Max 8 words untuk display headline |

### 📐 Layout — AI Patterns

| ❌ AI Slop | Kenapa | ✅ Ganti dengan |
|-----------|--------|----------------|
| **Centered hero** (text + CTA di tengah) | Default AI tiap kali bikin landing | Split screen, left-aligned + asset, atau asymmetric |
| **3 equal feature cards** | Paling obvious AI layout ever | Variasi ukuran (1 card besar + 2 kecil), grid asimetris |
| **Cards-inside-cards-inside-cards** | AI suka nesting cards tak terbatas | Flat hierarchy, border/divider cukup |
| **"Innovate" / "Empower" / "Revolutionize"** di headline | Kata buzzword AI | Bicara konkret: "Track your team's weekly velocity" |
| **Left text + right image** tiap section | Bolak-balik pattern membosankan | Variasi: full-width, grid, overlap, background-image |
| **Glassmorphism di semua card** | 2023 trend, udah mati | Border tipis + bg solid, glass cuma untuk overlay/navbar |
| **Infinite scroll + micro-animations everywhere** | Gak semua perlu gerak | Animasi cuma untuk hirarki & entry point |

### 🧠 Mental Model

Sebelum nulis kode, tanya: *"Apakah ini default yang bakal AI lain hasilkan juga?"*

Kalau jawabannya "iya" — **ubah pendekatannya.** Jangan puas sama layout pertama yang keluar dari prompt.

---

## Stack & Conventions

| Layer | Teknologi | Aturan |
|-------|-----------|--------|
| **Frontend** | Svelte 5 (runes) | Wajib `<script lang="ts">`, rune `$state`/`$derived`/`$props`, **jangan** `$effect` untuk derived state |
| **Routing SPA** | Inertia.js | Internal link wajib `use:inertia` dari `@inertiajs/svelte` |
| **HTTP Forms** | Inertia `router.post()`/`router.put()` | Jangan `<form>` biasa |
| **fetch()** | Manual fetch wajib `X-XSRF-TOKEN` header via `getCSRFToken()` | `lib/utils/csrf.ts` |
| **Styling** | Tailwind CSS v4 | File di `frontend/src/app.css` — semua token warna sudah di `@theme` |
| **Icons** | Lucide Svelte (`lucide-svelte`) | Satu keluarga icon, `stroke-width="2"` global |
| **Animasi** | Svelte `transition:` + CSS | GSAP untuk landing page canggih, Svelte transition cukup untuk dashboard |
| **Font** | Inter (sans) + JetBrains Mono (mono) — `@theme` di `app.css` |

## Color System

Semua token sudah terdefinisi di `frontend/src/app.css`. **Jangan define ulang.**

### Brand

| Token | Value | Usage |
|-------|-------|-------|
| `brand-400` | `#22d3ee` (cyan-teal) | Buttons, links, focus rings, primary accent |
| `brand-600` / `brand-700` | Darker cyan | Hover/active states, dark mode button bg |
| `secondary-500` | `#a855f7` (violet) | Secondary accent, premium highlights, gradients |

### Neutrals (cool navy-tinted)

| Token | Value | Usage |
|-------|-------|-------|
| `neutral-50` | `#f8fafc` | Light mode page bg |
| `neutral-950` | `#070b16` | Dark mode page bg — **jangan pure black** |
| `neutral-925` | `#0b111f` | Dark mode raised surface (card, sidebar) |
| `neutral-850` | `#172033` | Extra surface step |
| `neutral-800`..`600` | Slate scale | Text — `neutral-900` untuk headline, `neutral-600` body |

### Semantic

| Token | Value |
|-------|-------|
| `success` | `#10b981` (green) |
| `warning` | `#f59e0b` (amber) |
| `error` | `#ef4444` (red) |
| `info` | `#3b82f6` (blue) |

### Pola Penggunaan

- **Satu accent per viewport** — treat accent seperti highlighter yang cuma bisa dipakai sekali
- **CTA buttons**: `bg-brand-600 hover:bg-brand-700 text-white` + `shadow-lg shadow-brand-600/25`
- **Ghost buttons**: `bg-neutral-200/80 dark:bg-neutral-800 hover:bg-neutral-300/80`
- **Borders**: `border-neutral-200/80 dark:border-white/[0.04]` — invisible di dark mode, subtle di light
- **Cards**: bg-white dark:bg-neutral-925/50 + border + rounded-2xl
- **Gradient hero**: pakai utility `bg-gradient-hero` atau inline `linear-gradient(135deg, ...)`

## Shadow System

| Token | Value | Usage |
|-------|-------|-------|
| `shadow-soft` | Custom soft shadow | Default card shadow |
| `shadow-glow-brand` | Cyan glow | Premium highlight, hero section |
| `shadow-glow-brand-lg` | Larger cyan glow | Landing page hero CTA |

Dark mode tidak pakai shadow — pakai `border-white/[0.04-0.06]` subtle.

## Typography

### Scale

| Element | Class | Keterangan |
|---------|-------|-----------|
| Page H1 | `text-3xl font-bold tracking-tight` | Dashboard page title |
| Section title | `text-base font-semibold` | Card headers |
| Card title | `text-sm font-medium` | Item titles |
| Body | `text-sm text-neutral-600 dark:text-neutral-400` | Body text |
| Small | `text-xs text-neutral-500` | Meta, timestamps |
| Mono | `font-mono` | Data numbers |

### Hero (landing page)

- Max **2 lines** untuk headline di desktop
- Subtext max **20 words**, max 3-4 lines
- Default range: `text-4xl md:text-5xl lg:text-6xl`
- CTA harus visible tanpa scroll
- Top padding max `pt-24` — lebih dari itu keliatan floating

## Layout Patterns

### Bento Grid

```
grid md:grid-cols-2 lg:grid-cols-3 gap-5
```

- Mix cell sizes (col-span-2 untuk primary card)
- Rounded-2xl border cards
- Konsisten spacing `gap-5`
- **Anti-pattern:** jangan 3 card sama persis — variasi ukuran

### App Shell

```
div min-h-screen bg-neutral-50 dark:bg-neutral-950
  + fixed sidebar (w-72, hidden lg:flex)
  + fixed header (h-16, left-72)
  + content area (max-w-6xl mx-auto px-6 py-8)
```

### Auth Page

```
min-h-screen bg-neutral-50 dark:bg-neutral-950
  + centered card (max-w-md mx-auto)
```

### Page Structure (Inertia)

```svelte
<script lang="ts">
  import AppLayout from "@layouts/AppLayout.svelte";
  import type { User } from "@lib/types";

  interface Props {
    user?: User;
    success?: string;
    error?: string;
  }
  let { user }: Props = $props();
</script>

<AppLayout {user} group="dashboard">
  <div class="max-w-6xl mx-auto px-6 py-8 space-y-6">
    <!-- page content -->
  </div>
</AppLayout>
```

### Form Patterns

- Label **above** input, helper text optional, error text **below**
- `gap-2` untuk input blocks
- **Never** placeholder-as-label
- Submit via `router.post()` / `router.put()` / `router.delete()`
- Flash messages dari `$page.props.flash`

### Animasi & Transisi

- Entry animasi: `in:fly={{ y: 20, duration: 600 }}`
- Staggered list: `delay: 100 + i * 50`
- Dark mode toggle ada di `@components/DarkModeToggle.svelte`
- Responsive: mobile menu drawer pake `transition:fly={{ x: 300 }}`
- Hover card: `hover:border-brand-400/30 hover:shadow-xl hover:shadow-brand-400/5`
- Button press: `active:scale-[0.98]`

## Anti-Patterns (Teknis)

| ❌ Jangan | ✅ Ganti dengan |
|-----------|----------------|
| `$effect` buat derived state | `$derived()` |
| `h-screen` buat hero | `min-h-[100dvh]` |
| `<form>` biasa untuk Inertia | `router.post()` / `router.put()` |
| `<a href>` internal link | `<a use:inertia href="...">` |
| Inter font via Google link | `@theme` di `app.css` (Inter sudah di set) |
| Pure black `#000` / pure white `#fff` | `neutral-950` / `neutral-50` |
| `lucide-react` | `lucide-svelte` |
| Campur >1 keluarga icon | Satu: **Lucide Svelte** |
| fetch() ke /app/ tanpa CSRF | `X-XSRF-TOKEN` header via `getCSRFToken()` |
| Custom SVG icon manual | Pakai Lucide, cari glyph yang sesuai |
| Button text wrap di desktop | Perpendek label (max 3 words untuk CTA) |
| Duplicate CTA intent | Satu label per intent |

## Component File Convention

```
frontend/src/
├── components/        # Shared UI (DarkModeToggle, Logo)
├── layouts/           # AppLayout, AuthLayout — pakai Snippet children
├── pages/app/         # Authenticated pages (Dashboard, Profile)
├── pages/auth/        # Auth pages (Login, Register, ForgotPassword, ResetPassword)
└── lib/
    ├── types.ts       # User, Flash interfaces
    ├── utils/csrf.ts  # getCSRFToken()
    └── i18n/          # Terjemahan (en/id)
```

### Inertia Props Flow

```
Handler (Go) → inertiaService.Render(c, "app/Dashboard", fiber.Map{
    "user": user,         // dari session cache
    "flash": {...},       // auto-merged dari session.flash
    "errors": {...},      // validation errors
  }) → Inertia response → Svelte page ($page.props.user)
```

### Halaman Baru Checklist

- [ ] Handler file terpisah (`app/handlers/`) — ikut pattern module
- [ ] Svelte page di `frontend/src/pages/` — Svelte 5 runes
- [ ] Layout sesuai role (`AppLayout` untuk auth, `AuthLayout` untuk guest)
- [ ] Inertia internal links pakai `use:inertia`
- [ ] Form pakai `router.post()` / `router.put()` / `router.delete()`
- [ ] CSRF header untuk manual `fetch()`
- [ ] `@theme` tokens untuk warna — jangan hardcode hex
- [ ] Satu accent color per viewport
- [ ] max-w-[1400px] mx-auto untuk page container
- [ ] **Cek anti-slop section** — apakah output keliatan "AI generated"?
