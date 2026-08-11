import { createInertiaApp } from "@inertiajs/react";
import { createRoot } from "react-dom/client";

// React Fast Refresh — must load before any component imports in dev.
// @vitejs/plugin-react normally injects this preamble via transformIndexHtml,
// but with Inertia the Go server serves HTML, so Vite never sees it.
// We load it here manually and use lazy glob so components load after.
//
// Dynamic import is required: /@react-refresh is a Vite virtual module
// that only exists in dev — a static import would break production builds.
if (import.meta.env.DEV) {
	const { injectIntoGlobalHook } = await import("/@react-refresh");
	injectIntoGlobalHook(window);
	window.$RefreshReg$ = () => {};
	window.$RefreshSig$ = () => {};
}

const pages = import.meta.glob("./pages/**/*.tsx");

createInertiaApp({
	resolve: async (name) => (await pages[`./pages/${name}.tsx`]()).default,
	setup({ el, App, props }) {
		createRoot(el).render(<App {...props} />);
	},
});
