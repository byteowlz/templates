import tailwindcss from "@tailwindcss/vite";
import { execSync } from "node:child_process";
import path from "node:path";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

const port = Number(process.env.PORT) || 8790;

// Bind to the tailnet interface so the app is reachable fleet-wide (and
// registered at http://100.64.0.12:8774) without also exposing it on public
// / LAN interfaces. Resolve the machine's tailnet IP at startup; fall back to
// localhost off-tailnet. Override explicitly with HOST.
function resolveHost(): string {
	if (process.env.HOST) return process.env.HOST;
	try {
		const out = execSync("tailscale ip -4", {
			encoding: "utf8",
			stdio: ["ignore", "pipe", "ignore"],
		});
		const ip = out.trim().split(/\s+/)[0];
		if (ip) return ip;
	} catch {
		/* no tailscale — fall through */
	}
	return "localhost";
}

const host = resolveHost();

export default defineConfig({
	plugins: [react(), tailwindcss()],
	resolve: {
		alias: { "@": path.resolve(import.meta.dirname, "./src") },
	},
	// host is the tailnet interface IP (see resolveHost above), so the app is
	// reachable from the service registry and other machines — not 0.0.0.0.
	server: { port, host, strictPort: true },
	preview: { port, host, strictPort: true },
	// Relative base so a production build also works when opened from a
	// subpath or packaged as a self-contained oqto sandboxed-web bundle.
	base: "./",
});
