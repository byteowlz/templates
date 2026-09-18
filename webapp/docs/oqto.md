# docs/oqto.md — optional oqto-app packaging (ADR-0038)

This repo ships **structured for oqto discovery**: an `oqto-app.toml` manifest
at the root, and a build recipe that produces the self-contained
`sandboxed-web` bundle the manifest points at.

Honest status: the oqto runtime side of ADR-0038 (filesystem discovery,
Installation/binding records, capability grants, the Gate) is **designed, not
implemented** (`oqto-171p`). Today the app runs standalone; the manifest is a
pin for the future resolver, and `oqto-check` keeps the repo conformant so the
future migration is a move, not a rewrite.

## What's here

- `oqto-app.toml` — manifest (TOML, never YAML), schema pin `oqto-app/v0`,
  `presentations = ["sandboxed-web"]`, `entry = "bundle/index.html"`,
  `requested_capabilities = []`.
- `just oqto-bundle` — `bun run build` → copies `dist/` to `bundle/`
  (relative base, no external origins).
- `just oqto-check` — manifest sanity + bundle self-containment gate
  (fails on any external URL in the bundle).

## Hard rules (oqto-apps skill)

1. Surface to the host only via host capabilities (`files`, `kv`,
   `notifications`, `theme`, `user`) — no `fetch` to oqto, no DOM/window
   escapes, no imports from the oqto shell.
2. Theme via the `theme` capability (`Base24Scheme`/`ThemeMode`) — never
   hardcoded colors. This repo's token layer already speaks slots/roles; that
   is exactly what the capability expects.
3. Capabilities used = capabilities declared. Touching `host.kv` means adding
   `"kv"` to `requested_capabilities` in `oqto-app.toml`. Nothing in the seed
   app uses capabilities, so the manifest declares none.
4. Bundles are self-contained: no CDN, no network — `just oqto-check` enforces.
5. Manifests are TOML — never YAML.

## Promoting into an oqto workdir (when the resolver lands / for workbench dev)

The oqto frontend (mock host + real UI, no oqto login) is the verification
loop: `http://localhost:3000/workbench.html`. To develop against it, move this
app's source into `<oqto-frontend>/mini-apps/<app-id>/`, define the app with
the mini-apps SDK, and register it in the standalone workbench registry:

```ts
import { defineOqtoApp } from "@/mini-apps/sdk";

export const app = defineOqtoApp({
	id: "{{project_name}}",
	title: "{{project_name}}",
	requestedCapabilities: [], // keep in sync with oqto-app.toml
	component: App,            // the same root component; host access via useOqtoHost()
});
```

```ts
// mini-apps/workbench/registry.ts
export const standaloneApps: ReadonlyArray<OqtoApp> = [
	/* ...existing */ app,
];
```

Verify with the oqto-apps skill quick-start gates (typecheck, biome, vitest,
build, screenshot loop at 1440×900 and 390×844). A declarative presentation
(`ui/entry.json`, profile `oqto/tables-v1`) is opt-in later — add the
`[presentation.declarative]` block to the manifest only when real payloads
exist.
