# {{project_name}}

A focused byteowlz webapp (workbench/tool) built from the `webapp` template.

## Stack

- **Bun + Vite + TypeScript (strict)** — one command to dev, one to build.
- **React 19** — component model agents already know from oqto/foxline apps.
- **Tailwind CSS v4** via `@tailwindcss/vite` — token layer in `src/index.css`.
- **byteowlz design system** — house scheme `oqto-dark` committed as slot → role
  tokens (see `docs/design-system.md`). No component-library lock-in; `Button`
  and `Card` ship as token-pure exemplars.

## Quickstart

```bash
bun install
just dev          # vite dev server + register in the service registry
```

`just` shows every recipe: `check` (typecheck+lint+build), `serve` (build+preview,
registered), `omarchy` (apply the machine's omarchy theme), `oqto-bundle` /
`oqto-check` (oqto opt-in artifacts), `release-tag`.

## Optional integrations

| Opt-in | Entry point | Doc |
|---|---|---|
| Omarchy theme follow | `just omarchy` → `src/theme/omarchy.generated.css` | `docs/omarchy.md` |
| Oqto app (ADR-0038) | `oqto-app.toml` + `just oqto-bundle` | `docs/oqto.md` |
| Design-system engine (JS) | `byt design-system sync` + `applyScheme()` | `docs/design-system.md` |

## Service registry

Dev and serve register the UI at `http://100.64.0.12:8774` (group `misc`,
port `8790` by default — override with `PORT=`, `GROUP=`, `NAME=`) and
unregister on exit. Keep that contract: every byteowlz web UI should be
discoverable fleet-wide.

## Conventions

Standard byteowlz template furniture: `AGENTS.md` (agent guidance), `CONTEXT.md`
(domain glossary), `docs/adr/` (durable decisions), `PRODUCT.md`/`DESIGN.md`
(impeccable design context), `CHANGELOG.md`.
