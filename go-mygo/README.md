# {{project_name}}

MyGo **v0.3.2 system-webview** desktop starter: native Go backend, HTML/CSS/
TypeScript presentation. This is **not** MyGo's fully native GPU-rendered UI.
Vanilla TypeScript + Vite keeps the small demo dependency-light; no React,
Three.js, HTTP service or embedded model is required.

## Start

Go 1.27.1 (automatic toolchain download permitted), Bun and just are required.
macOS uses system WKWebView; Linux needs GTK3/WebKitGTK 4.1; Windows needs WebView2.

```sh
just install
just check
just agent-check
just dev
just build       # CLI + native package with embedded dist; skips macOS DMG
just package     # full current-platform package; local/ad-hoc, no auto-upload
```

`go build .` alone compiles but does **not** embed the frontend. Use MyGo's
build command for runnable distribution. `just generate` regenerates typed
`src/mygo.ts` from the bound Go Desktop facade; never hand-edit it.
`just scaffold-check` copies/replaces placeholders into a temporary scaffold
and executes frozen install, check, build and the fixture gate.

## Headless example

```sh
build/{{project_name}}ctl snapshot --json > envelope.json
# Snapshot input is the result.snapshot object, not the response envelope.
build/{{project_name}}ctl validate --input examples/snapshot.json --json
cat examples/snapshot.json | build/{{project_name}}ctl validate --json
build/{{project_name}}ctl apply --json # AUTHORITY_DENIED, exit 1; always read-only
build/{{project_name}}ctl schema
build/{{project_name}}ctl config schema
```

Desktop and CLI use one Go service. There is no mutation adapter or external
device discovery, and no silent fallback from live state to fixtures.
See [contracts](docs/contracts.md), [design/provenance](docs/design-system.md),
[agent gate](docs/agent-readiness.md) and [architecture](docs/adr/0001-substrate.md).

Configuration: flags > explicit `--config` TOML > `MYGO_APP_*` environment >
`./config.toml` > XDG global > defaults. Scalar overlay, no list merging.
Global: `$XDG_CONFIG_HOME/{{project_name}}/config.toml`, falling back to
`~/.config/{{project_name}}/config.toml` on every platform (not macOS Preferences).
Desktop first run creates a private default only when no file exists.
Headless inspection/help/generation is read-only; `config init` explicitly
creates defaults without overwriting. `config show --json` shows effective values.
Env keys: `MYGO_APP_THEME`, `MYGO_APP_RADIUS`, `MYGO_APP_OMARCHY_FILE`.
Missing explicit/invalid files fail, not silently ignored.

## Review boundaries

The starter is illustrative and **unapproved**. It inherits token mechanism,
not Oqto composition or Vial approval. Consumers need a product/surface/variant/
revision/approval record in [Design Studio](https://github.com/byteowlz/design-studio/blob/main/VISUAL-CONTRACT.md)
before production UI changes. No 3D engine is included; a consumer can opt in,
with a semantic 2D equivalent and untrusted-label safety.

`just agent-check` proves only the read-only synthetic fixture contract.
Compilation, generated bindings or report validation do not certify native
IPC, accessibility, platform behavior or hardware. Consumers replace this
fixture with domain tests, including applicable mutation/authority gates.
