import type { APIRoute, GetStaticPaths } from "astro";
import { getCollection } from "astro:content";

export const prerender = true;

export const getStaticPaths: GetStaticPaths = async () => {
	const docs = await getCollection("docs");
	return docs.map((entry) => ({
		params: { slug: entry.id },
	}));
};

export const GET: APIRoute = async ({ params }) => {
	const slug = params.slug;
	if (!slug) return new Response("Not found", { status: 404 });

	const svg = `<svg width="1200" height="630" xmlns="http://www.w3.org/2000/svg">
		<rect width="1200" height="630" fill="#070b16"/>
		<rect x="0" y="0" width="1200" height="630" fill="url(#grad)"/>
		<defs>
			<radialGradient id="grad" cx="50%" cy="0%" r="60%">
				<stop offset="0%" stop-color="rgba(34,211,238,0.12)"/>
				<stop offset="100%" stop-color="transparent"/>
			</radialGradient>
		</defs>
		<g transform="translate(80, 200)">
			<rect x="0" y="0" width="48" height="48" rx="10" fill="#0b111f" stroke="#1e293b"/>
			<path d="M12 34L24 12L36 34" stroke="#22d3ee" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
			<path d="M18 34L24 23L30 34" stroke="#a855f7" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" fill="none" opacity="0.7"/>
		</g>
		<text x="80" y="310" font-family="Inter, system-ui, sans-serif" font-size="52" font-weight="800" fill="#e2e8f0" letter-spacing="-2">
			Laju Go
		</text>
		<text x="80" y="370" font-family="Inter, system-ui, sans-serif" font-size="28" font-weight="400" fill="#64748b">
			High-performance Go SaaS boilerplate
		</text>
		<text x="80" y="420" font-family="JetBrains Mono, monospace" font-size="18" fill="#22d3ee">
			Go Fiber · Svelte 5 · Inertia.js · SQLite · templ
		</text>
		<text x="80" y="540" font-family="JetBrains Mono, monospace" font-size="16" fill="#475569">
			${slug}
		</text>
	</svg>`;

	return new Response(svg, {
		headers: { "Content-Type": "image/svg+xml" },
	});
};
