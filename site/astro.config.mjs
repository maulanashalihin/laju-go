// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import starlightBlog from "starlight-blog";
import mdx from "@astrojs/mdx";

export default defineConfig({
	site: "https://laju.dev",
	integrations: [
		starlight({
			plugins: [
				starlightBlog({
					title: "Blog",
					postCount: 20,
					recentPostCount: 20,
					authors: {
						maulana: {
							name: "Maulana Shalihin",
							title: "Laju Go author",
							url: "https://github.com/maulanashalihin",
						},
					},
				}),
			],
			title: "Laju Go",
			description:
				"High-performance SaaS boilerplate: Go Fiber + Svelte 5 / React 19 / Vue 3 + Inertia.js + SQLite + templ. Auth, uploads, migrations, tests, Docker — wired end to end.",
			favicon: "/favicon.svg",
			logo: {
				src: "./src/assets/logo.svg",
				alt: "Laju Go logo",
				replacesTitle: false,
			},
			social: [
				{
					icon: "github",
					label: "GitHub",
					href: "https://github.com/maulanashalihin/laju-go",
				},
			],
			customCss: ["./src/styles/custom.css"],
			components: {
				Head: "./src/components/Head.astro",
				ThemeSelect: "./src/components/ThemeSelect.astro",
			},
			sidebar: [
				{
					label: "Getting Started",
					items: [
						{ label: "Introduction", slug: "getting-started/introduction" },
						{ label: "Installation", slug: "getting-started/installation" },
						{
							label: "Building with AI agents",
							slug: "getting-started/ai-agents",
						},
					],
				},
				{ label: "Philosophy", slug: "philosophy" },
				{
					label: "Architecture",
					items: [
						{ label: "Overview", slug: "architecture/overview" },
						{ label: "Three-tier rule", slug: "architecture/three-tier" },
						{ label: "Conventions", slug: "architecture/conventions" },
					],
				},
				{
					label: "Frontend",
					items: [
						{ label: "Svelte 5 + Inertia", slug: "frontend/svelte-inertia" },
						{ label: "React 19 + Inertia", slug: "frontend/react-inertia" },
						{ label: "Vue 3 + Inertia", slug: "frontend/vue-inertia" },
						{ label: "templ components", slug: "frontend/templ" },
					],
				},
				{
					label: "Auth",
					items: [
						{ label: "Sessions & middleware", slug: "auth/sessions" },
						{ label: "Google OAuth", slug: "auth/google-oauth" },
						{ label: "Password reset", slug: "auth/password-reset" },
					],
				},
				{
					label: "Database",
					items: [
						{
							label: "Schema & migrations",
							slug: "database/schema-migrations",
						},
						{ label: "sqlc code generation", slug: "database/sqlc" },
						{ label: "SQLite configuration", slug: "database/sqlite-configuration" },
						{ label: "Data protection & recovery", slug: "database/data-protection" },
					],
				},
				{
					label: "Uploads",
					items: [
					{ label: "Overview", slug: "uploads/overview" },
					{ label: "Avatar upload", slug: "uploads/avatar" },
					{ label: "TUS resumable upload", slug: "uploads/tus" },
					],
				},
				{
					label: "Deployment",
					items: [
						{ label: "Overview", slug: "deployment/overview" },
						{ label: "AI agent", slug: "deployment/ai-agent" },
						{ label: "Docker", slug: "deployment/docker" },
						{ label: "Linux VPS", slug: "deployment/vps" },
						{ label: "Reverse proxy", slug: "deployment/reverse-proxy" },
						{ label: "Configuration", slug: "deployment/configuration" },
					],
				},
				{ label: "Testing", slug: "testing" },
				{ label: "Troubleshooting", slug: "troubleshooting" },
				{ label: "Contributing", slug: "contributing" },
			],
		}),
		mdx(),
	],
});
