# docs/design-system.md — the byteowlz design system in this repo

Source of truth: `~/byteowlz/design-system` (spec + `@byteowlz/design-system`
engine). This file is the consumer contract for THIS repo.

## How the token layer works

`src/index.css` commits a three-tier token architecture:

```
base24 slots (--base00..--base17)   ← the only place color literals live
        ↓
abstract roles (16, closed set)     ← background, surface, surface-sunken,
        ↓                              foreground, muted-foreground, primary,
                                       secondary, muted, accent, success,
                                       warning, danger, info, border, ring, input
        ↓
framework vocabulary (shadcn-ish)   ← card/popover pairs, sidebar, charts,
                                       destructive (non-portable convenience)
```

Tailwind v4's `@theme inline` block exposes all of the above as utilities
(`bg-surface`, `text-muted-foreground`, `border-border`, `bg-primary/20`, ...).
Alpha (`/20`, scrims, hovers) is derived from the resolved role — never a new
slot.

Committed defaults: slots of the house scheme **oqto-dark** on `:root`,
**oqto-light** under `.light`. Radius dial `--radius: 0px` (sharp, house
identity); the tier scale is proportional (R1: sm = 0.5×, md = 0.75×, lg = 1×,
xl = 1.25×), so dial 0 → everything sharp.

Slot var names use the engine's exact casing (`--base0B`, not `--base0b`) —
CSS custom properties are case-sensitive, and matching the engine lets
`applyScheme()` override every slot inline later.

## Recoloring this tool (its *look*)

Three sanctioned mechanisms, in order of preference:

1. **Pick a different house scheme**: copy slots from
   `~/byteowlz/design-system/impls/shadcn-ts/src/schemes/*.json` into
   `src/index.css` (`:root` / `.light`).
2. **Omarchy follow**: `just omarchy` maps the machine's active omarchy theme
   onto the slots at runtime — see `docs/omarchy.md`.
3. **Engine (dynamic switching)**: vendor the engine and call `applyScheme()` —
   see below.

Identity (fonts, radius dial, spacing) is per-tool by design: it lives in
`src/index.css` `:root` and `DESIGN.md`, never in the mechanism.

## Vendoring the JS engine (optional)

Only needed if the app must switch schemes at runtime (e.g. a scheme picker).
For static looks, the committed token layer is enough.

```bash
# from the repo root (byt discovers design-system.toml)
cat > design-system.toml <<'EOF'
# Vendored design-system snapshot — managed by `byt design-system sync`.
ref  = "v0.1.0"
impl = "shadcn-ts"
path = "vendor/design-system"
EOF
byt design-system sync
bun add @byteowlz/design-system@file:./vendor/design-system
```

```ts
import { applyScheme } from "@byteowlz/design-system";
applyScheme(scheme, { radius: "0px" }); // writes slots+roles+vocabulary inline
```

Keep the `design-system.toml` pin updated (`byt design-system status`).
