#!/usr/bin/env node
// Read-only exhibit. Exact names, not filesystem paths supplied by the client.
import { createServer } from "node:http";
import { readFile, realpath, stat } from "node:fs/promises";
import { dirname, resolve, sep } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(fileURLToPath(import.meta.url));
// Curate here AND in index.html. No directory listings, globbing or user input.
export const media = Object.freeze({
  "/media/method-a.gif": "demo/method-a.gif",
  "/media/method-b.gif": "demo/method-b.gif",
  "/media/source.svg": "demo/source.svg",
  "/media/diff.svg": "demo/diff.svg",
  "/media/repair.svg": "demo/repair.svg",
});
const contentTypes = { gif: "image/gif", svg: "image/svg+xml" };

export function createReviewServer() {
  return createServer(async (req, res) => {
    const send = (status, content, type = "text/plain; charset=utf-8") => {
      res.writeHead(status, {
        "Content-Type": type,
        "Cache-Control": "private, max-age=120",
        "X-Content-Type-Options": "nosniff",
        "Referrer-Policy": "no-referrer",
        "Content-Security-Policy":
          "default-src 'none'; img-src 'self'; style-src 'unsafe-inline'; script-src 'unsafe-inline'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'",
      });
      res.end(req.method === "HEAD" ? undefined : content);
    };
    if (req.method !== "GET" && req.method !== "HEAD")
      return send(405, "Read-only");
    // Match the unnormalized request target: decoded/normalized dot segments
    // must never be interpreted as file paths.
    const target = req.url?.split("?")[0];
    if (target === "/health") return send(200, "ok");
    const relative =
      target === "/"
        ? "index.html"
        : Object.hasOwn(media, target)
          ? media[target]
          : undefined;
    if (!relative) return send(404, "Not available");
    const path = resolve(root, relative);
    const approvedRoot = resolve(
      root,
      relative === "index.html" ? "." : "demo",
    );
    try {
      const [actual, safe] = await Promise.all([
        realpath(path),
        realpath(approvedRoot),
      ]);
      if (
        relative === "index.html"
          ? actual !== path
          : !actual.startsWith(safe + sep)
      )
        return send(404, "Not available");
      if (!(await stat(actual)).isFile()) return send(404, "Not available");
      const type =
        relative === "index.html"
          ? "text/html; charset=utf-8"
          : contentTypes[relative.split(".").at(-1)];
      return send(200, await readFile(actual), type);
    } catch {
      return send(404, "Not available");
    }
  });
}

if (
  process.argv[1] &&
  resolve(process.argv[1]) === fileURLToPath(import.meta.url)
) {
  const port = Number(process.env.PORT || 8792);
  const hosts = [
    ...new Set([
      "127.0.0.1",
      ...(process.env.REVIEW_TAILNET_IP ? [process.env.REVIEW_TAILNET_IP] : []),
    ]),
  ];
  const servers = [];
  for (const host of hosts) {
    const server = createReviewServer();
    server.on("error", (error) => {
      console.error(`review ${host}:${port}: ${error.message}`);
      process.exitCode = 1;
      servers.forEach((active) => active.close());
    });
    server.listen(port, host, () =>
      console.log(`review http://${host}:${port}/`),
    );
    servers.push(server);
  }
  for (const signal of ["SIGINT", "SIGTERM"])
    process.on(signal, () => servers.forEach((server) => server.close()));
}
