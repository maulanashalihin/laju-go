import { defineConfig } from "vitest/config";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { resolve } from "path";

export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      "@": resolve(__dirname, "frontend/src"),
      "@components": resolve(__dirname, "frontend/src/components"),
      "@layouts": resolve(__dirname, "frontend/src/layouts"),
      "@pages": resolve(__dirname, "frontend/src/pages"),
    },
    conditions: ["browser"],
  },
  test: {
    environment: "happy-dom",
    include: ["frontend/src/**/*.{test,spec}.{js,ts}"],
    clearMocks: true,
  },
});
