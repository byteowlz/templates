import tailwindcss from "@tailwindcss/vite";
import path from "node:path";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

const port = Number(process.env.PORT) || 8790;

// host: true makes the dev/preview server reachable from the tailnet, which is
// what the byteowlz service registry (http://100.64.0.12:8774) links to.
export default defineConfig({
	plugins: [react(), tailwindcss()],
	resolve: {
		alias: { "@": path.resolve(import.meta.dirname, "./src") },
	},
	server: { port, host: true, strictPort: true },
	preview: { port, host: true, strictPort: true },
	// Relative base so a production build also works when opened from a
	// subpath or packaged as a self-contained oqto sandboxed-web bundle.
	base: "./",
});
