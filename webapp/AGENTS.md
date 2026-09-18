# AGENTS.md

Guidance for coding agents working on this focused webapp (built from the
byteowlz `webapp` template).

## Domain and Architecture

- Read `CONTEXT.md` before domain work and use its canonical terms consistently.
- Update `CONTEXT.md` when domain language is resolved. Keep it a glossary only.
- Read `docs/adr/` before architectural changes. Add an ADR only for decisions
  that are hard to reverse, surprising without context, and based on a real
  trade-off.

## Design system discipline (non-negotiable)

- `src/index.css` is the token layer: base24 slots → roles → framework
  vocabulary. See `docs/design-system.md`. The house scheme `oqto-dark` is the
  committed default; `.light` carries `oqto-light`.
- Components use role tokens only (`bg-background`, `text-muted-foreground`,
  `border-border`, `text-success`, ...). Never hardcode hex/oklch literals in
  components. Slots (`--base0B`) are for the token layer and generated CSS only.
- Radius via `--radius-sm/md/lg/xl` (proportional to the `--radius` dial, R1).
  Default dial is `0` (sharp) — the house identity. Changing it is a per-tool
  identity decision recorded in `DESIGN.md`.
- Status colors: `success`/`warning`/`danger`/`info` — no ad-hoc reds/greens.
- Dark-first. If a light mode is needed, drive it via the `.light` class or the
  omarchy adapter — never by scattering literals.

## UI work: impeccable skill

- Use the `impeccable` skill for any UI design/build/polish work. `PRODUCT.md`
  and `DESIGN.md` are its durable context — keep them current.
- These workbenches are **Operate**-mode surfaces: scanability, consistency and
  the real usage scene outrank expression. Brand lives in precise details.
- Verify in bounded passes (screenshots at 1440×900 and 390×844, one batched
  defect pass, one confirm round — then stop).
- Responsive by default; 44px touch targets on interactive elements; visible
  `:focus-visible`; honor `prefers-reduced-motion` (the base layer already does).

## Workflow

- Bun for everything: `just install`, `just dev`, `just check`.
- Run `just check` (typecheck + oxlint + build) after significant edits.
- Keep dependencies minimal; this is a focused tool, not a platform. New deps
  need a reason in the PR/commit message.
- Never publish artifacts to public registries without explicit user approval.
- Do exactly what the user asks — no unsolicited files or docs.

## Service registry protocol

- `just dev` and `just serve` register this UI at the byteowlz service registry
  (`http://100.64.0.12:8774`, default group `misc`, default port `8790`) and
  unregister on exit. If you start a dev server by hand, run `just register`
  and `just unregister` around it.
- Pick a free `PORT` if 8790 is taken; set `GROUP=` to the tool's domain
  (pixelgym | voice | ops | infra | gpu | misc).

## Omarchy adapter (optional)

- `src/theme/omarchy.generated.css` is **machine-local, generated** — never hand
  edit, never commit personal theme values. The committed state is the no-op
  passthrough. Regenerate with `just omarchy` (see `docs/omarchy.md`).

## Oqto app (optional)

- `oqto-app.toml` is TOML, never YAML. Capabilities used = capabilities
  declared. Bundles must be self-contained (`just oqto-check` gates this).
- The oqto runtime (discovery, grants) is designed, not implemented
  (`oqto-171p`): today the app runs standalone. See `docs/oqto.md` before
  touching capabilities or the manifest.

## Justfile Commands

Run `just` to see all recipes. Core: `dev`, `build`, `check`, `serve`,
`omarchy`, `oqto-bundle`, `oqto-check`, `register`/`unregister`, `release-tag`.
