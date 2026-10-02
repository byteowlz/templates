# AGENTS.md

Guidance for coding agents working on this GPUI desktop template. Deliberately
short: it keeps the enforceable boundaries and workflows, and points at
machine-checked configuration instead of duplicating any inventory.

## Source of truth

Static facts are machine-checked, not copied here:

- Crates, dependencies, lint settings -> `Cargo.toml` (the authoritative
  enumeration is `cargo metadata --no-deps --format-version 1`).
- Task commands -> `just` (run `just` to list). `just check-all` is what a PR
  must pass (fmt + clippy with warnings denied + drift-check + tests).
- Issues -> `trx` (see below), never markdown TODOs.
- Exact gpui/gpui-component APIs -> the pinned crates in the local Cargo
  registry checkout, not the website. This is critical: gpui APIs change
  between releases, so read the source of the *pinned* version you depend on.

If `cargo metadata`, `just`, or `scripts/drift-check.sh` disagree with anything
here, the command is right and this file is wrong.

## GPUI: one dependency, all of it under `gpui_kit`

`gpui-kit` is the **single** dependency for a GPUI desktop app. An application
never lists the `gpui-pre-*` crates directly; gpui-kit pulls the matching set
and re-exports it as `gpui_kit::*`:

| Path                  | Crate            | Feature     |
| --------------------- | ---------------- | ----------- |
| `gpui_kit::*`         | `gpui`           | always      |
| `gpui_kit::platform`  | `gpui_platform`  | always      |
| `gpui_kit::base`      | `gpui-base`      | always      |
| `gpui_kit::component` | `gpui-component` | `component` |
| `gpui_kit::assets`    | `gpui-kit-assets`| `assets`    |

Key facts:

- `gpui_kit::application().run(|cx| ...)` opens the platform; `gpui_kit::init(cx)`
  initializes the enabled layers; `component::Root` is the **first** view in a
  window (it owns Sheet/Dialog/Notification/implicit native menu).
- `cx.theme()` returns the active `component::Theme`; it **derefs to
  `ThemeColor`**, so `cx.theme().foreground`, `.primary`, `.border`,
  `.sidebar`, `.muted_foreground`, `.danger`, … resolve directly.
- Views implement `Render`; layout uses the `div()` element and the `Styled`
  / `ParentElement` builders (`flex`, `gap_3`, `px`, `w`, `bg`, `rounded`…).
- A `cx.listener(|this, event, window, cx| …)` callback takes **four**
  arguments. Use it with `Button::new(id).on_click(…)`.
- Read the pinned source for the exact method set. `FontWeight` is
  `FontWeight::BOLD` (consts are uppercase); `Pixels` exposes `.as_f32()`.

## Domain and architecture (the two-crate seam)

- Read `CONTEXT.md` before domain work; keep it a glossary only.
- `{{project_name}}-design` is the **design-system mechanism**: base24 scheme
  (data) -> closed abstract `Role` -> a gpui-component `ThemeSet` JSON
  document. It holds **no gpui types** so it is target-agnostic (the candidate
  for `design-system/impls/gpui-rs`) and can move there unchanged.
- `{{project_name}}` is the only crate that imports gpui-kit. It hosts the app
  (`main.rs`), the root view (`app.rs`), configuration (`config.rs`) and the
  theme wiring (`theme.rs`).
- Keep the seam: never put gpui types into `{{project_name}}-design`; never
  duplicate the design-system mechanism into the app crate.

## Design-system rules

- A **scheme** is data (`crates/{{project_name}}-design/schemes/*.json`). It is
  the only layer a look fills.
- A **role** is an intent bound to a slot (`spec/roles.md`, 16 closed roles).
  Consumers read the role layer, never a scheme slot or a host hex.
- **Never hardcode a color that exists as a role.** Read it from `cx.theme()`
  (e.g. `theme.primary`, `theme.border`, `theme.muted_foreground`). The
  `example-*.json` files are starter palette, not a house look: replace them
  with your tool's schemes and set the `Identity` (fonts, radius dial, shadow).
- Alpha is **derived** at the implementation layer (`color.alpha(0.2)`), never
  an extra slot.
- Radius is a single proportional dial (`RadiusScale`), sharp-at-0. Pick a dial
  for your identity; `0` is a legitimate choice.
- The theme set must parse with gpui-component's schema. `theme.rs::init`
  loads it via `ThemeRegistry::load_themes_from_str`, and `Theme::change`
  refreshes both the component `Theme` and the gpui-base copy (markdown,
  scrollbars). Only updating the former leaves markdown on toolkit defaults.

## Workflow

- Before anything significant: `just check-all`.
- Config: `crates/{{project_name}}/src/config.rs` loads
  `$XDG_CONFIG_HOME/{{project_name}}/config.toml` (a default is written on first
  run). **Configuration and machine output are JSON/TOML only** — never add a
  YAML mode or a YAML dependency.
- `unsafe_code` is forbidden; the scoped allow is only the AppKit/objc2 bridge
  if you add one (see how `oqto-desktop` scopes its `mac_chrome` and `live`
  probes if you need it).
- Run the app with `just run`. A GPUI window takes focus, so on a shared
  machine prefer building/checking and let the human launch it.

## Issue tracking (trx)

Use `trx` for all issue tracking — never markdown TODOs.

```bash
trx ready --json
trx create "Title" -t task -p 2 --json
trx update <id> --status in_progress --json
trx close <id> -r "reason" --json
```

Priorities: 0=critical, 1=high, 2=medium (default), 3=low, 4=backlog. Issue
state lives in `.trx/` (JSONL) — commit it with code changes.

## House rules

- Do exactly what the user asks — no unsolicited files.
- Keep README updates concise and emoji-free.
- Never commit secrets or sensitive paths; scrub logs.
- `Cargo.lock` is committed; bump manifests + lock together.
- Screenshots of real session/history content are user data: inspect locally,
  never commit or copy text into shared docs.