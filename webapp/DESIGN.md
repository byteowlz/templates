# DESIGN.md — {{project_name}}

(Impeccable durable context. `impeccable document`/`new-work` owns this file.
Canonical token contract lives in `docs/design-system.md`; this file records
THIS tool's look — its identity.)

## Look

- **Scheme**: house default `oqto-dark` (committed in `src/index.css`);
  `.light` carries `oqto-light`. Optional omarchy follow per `docs/omarchy.md`.
- **Radius dial**: `0px` (sharp) — house identity. Only change with intent,
  and record the new value here.
- **Type**: monospace-first (`ui-monospace` stack), 14px/1.5 body. Titles
  lowercase, bold, wide tracking — matching byteowlz services.
- **Density**: compact; 44px minimum interactive targets for touch safety.
- **Motion**: quiet. 150ms color transitions only; honor reduced-motion.

## Layout grammar

Single-column workbench shell: header (title · one-line desc · status dot) →
main (cards/panels on `surface`, app on `background`) → no footer chrome.
Elevated surfaces use `card`; recessed rails/terminals use `surface-sunken`.

## Bans

- No hardcoded colors in components (role tokens only).
- No radius other than the dial tiers.
- No decorative gradients/shadows; elevation is a border + surface role change.
