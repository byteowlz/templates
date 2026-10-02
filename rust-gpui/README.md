# {{project_name}}

A native GPUI desktop application built on the byteowlz design system.

- **One dependency** for GPUI: `gpui-kit` (it pulls and re-exports the matching
  `gpui-pre-*` set as `use gpui_kit::*;`).
- **Design system built in**: base24 scheme -> closed role layer -> gpui theme
  set. A scheme is data; reskin the scheme and the whole app follows.
- **Two-crate seam**: `{{project_name}}-design` (the target-agnostic mechanism,
  no gpui) and `{{project_name}}` (the only crate that imports gpui-kit).

> The bundled `example-*.json` schemes are **starter palette**, not a house
> look. Replace them with your tool's schemes and set the `Identity` (fonts,
> radius dial, shadow).

## Quick start

Install the latest stable Rust toolchain (`rustup default stable`). GPUI needs
a native platform backend — on macOS this is a normal `cargo run`; on Linux you
need the platform system deps, on Windows the WebView2/Vulkan prerequisites.

```bash
just check-all   # fmt + clippy(with warnings denied) + drift-check + tests
just run         # launch the native window
```

The window shows the design-system role vocabulary live: a sunken sidebar, an
elevated card, primary/outline buttons, status badges and a dark/light toggle.

## Structure

```
crates/
  {{project_name}}-design/    # design-system mechanism (scheme -> roles -> theme set), no gpui
    src/{color,scheme,roles,radius,theme_set}.rs
    schemes/{example-dark,example-light}.json
  {{project_name}}/           # the app; the only crate that imports gpui-kit
    src/{main,app,config,theme}.rs
examples/
  config.toml                 # documented default
  config.schema.json          # editor/LSP schema for the config file
docs/adr/                     # architecture decision records (template + README)
scripts/drift-check.sh        # verifies documented claims against real manifests
justfile                      # command runner (just check-all is the PR gate)
```

## The design-system flow

1. A **scheme** JSON fills only the chromatic slots (`base00`-`base17`).
2. [`{{project_name}}-design`](crates/{{project_name}}-design/src/lib.rs) resolves
   the 16 closed **roles** against the scheme.
3. [`theme_set.rs`](crates/{{project_name}}-design/src/theme_set.rs) emits a
   gpui-component `ThemeSet` (role-mapped colors + syntax) for a scheme +
   identity.
4. [`theme.rs`](crates/{{project_name}}/src/theme.rs) loads it into the
   `ThemeRegistry` and applies the configured scheme.
5. Views read colors from `cx.theme()` (the role layer) — never a host hex.

## Configuration

Settings live in `$XDG_CONFIG_HOME/{{project_name}}/config.toml` (a default is
written on first run). See `examples/config.toml` and `examples/config.schema.json`.

## GPUI resources

`gpui-kit` is the binding and `gpui-kit.com` is the guide. The source of the
pinned crates in your Cargo registry checkout is the precise API reference.

## License

MIT — see the workspace `Cargo.toml`.