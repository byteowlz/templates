# Working vocabulary

A glossary only — cross-repo vocabulary lives in the owning project's over
CONTEXT.md. Add terms here when a word gains a precise meaning in this codebase.

## Application

- **scheme**: the chromatic base (`base00`-`base17`), the only layer a look
  fills. A base24 scheme (or a base16 scheme widened by the derivation policy).
- **role**: a closed semantic intent bound to a slot (e.g. `primary` ->
  `base0B`, `danger` -> `base08`). The portable surface consumers read.
- **identity**: per-tool, non-color presentation facts — font families, radius
  dial, shadow. Distinct from the scheme.
- **theme set**: a gpui-component `ThemeSet` document (one entry per scheme)
  that an app loads to activate a look. Produced from scheme(s) + identity.
- **look**: a concrete scheme + identity pair worn by one tool.

## GPUI

- **Root**: the first view in a window in gpui-component; owns Sheet/Dialog/
  Notification/native-menu management. Wrap your root view in it.
- **view**: a Rust struct implementing `Render` (and usually `Into`-able to an
  `Entity`). Layout is the `div()` element + `Styled`/`ParentElement` builders.
- **theme** (`cx.theme()`): the active `Theme`; derefs to `ThemeColor` so
  `foreground`, `primary`, `border`, `sidebar`, `muted_foreground`, `danger`
  etc. are the role-mapped colors to read.
- **window**: the native window a view renders into; `WindowOptions` sets
  bounds, min size, titlebar and `app_id`.

## Config

- **config.toml**: the app's user setting file at
  `$XDG_CONFIG_HOME/{{project_name}}/config.toml`. Resolution order: CLI args >
  explicit file > env vars > this file.