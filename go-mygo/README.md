# {{project_name}}

Fully native MyGo **v0.3.2 GPU UI** (`ui.View`) + shared Go service.
**No HTML, webview, JavaScript, Bun/TypeScript frontend or IPC bindings.**
Native views call Go directly. The earlier system-webview prototype is
unselected and checkpointed separately; it is not this template.

## Start

Go 1.27.1 (automatic toolchain download allowed) and just are required.
macOS uses Metal; Linux needs GTK3 (not WebKitGTK); Windows uses native
Direct3D11/WARP (no WebView2). Platform runtime proof is separate from compile.

```sh
just install
just check
just agent-check
just dev
just build       # semantic CLI + native host package; macOS DMG skipped
just package     # full current-platform native package; no upload
just native-proof # synthetic real window, content captures and GPU-path log
```

`just scaffold-check` copies and replaces all project placeholders into a
temporary directory and runs real check/build/fixture proof. There is no
frontend asset build to embed; scheme/role/schema data are Go-embedded.

## Semantic headless example

```sh
build/{{project_name}}ctl snapshot --json
build/{{project_name}}ctl validate --input examples/snapshot.json --json
cat examples/snapshot.json | build/{{project_name}}ctl validate --json
build/{{project_name}}ctl apply --json # always AUTHORITY_DENIED; exit 1
build/{{project_name}}ctl schema
build/{{project_name}}ctl config schema
```

Desktop and CLI share one service. No live discovery/write backend or silent
fixture fallback exists. See [contracts](docs/contracts.md),
[design/provenance](docs/design-system.md), [required agent gate](docs/agent-readiness.md)
and [recorded proof](docs/validation.md).

Config: flags > explicit `--config` TOML > `MYGO_APP_*` environment >
`./config.toml` > XDG global > defaults. Scalars overlay; no list merging.
Global: `$XDG_CONFIG_HOME/{{project_name}}/config.toml` or
`~/.config/{{project_name}}/config.toml` on every platform.
Desktop first run writes a 0600 default only if no config file exists.
Headless inspection/help are read-only; `config init` explicitly creates
defaults without overwrite, `config show --json` reports effective values.
Env: `MYGO_APP_THEME`, `MYGO_APP_RADIUS`, `MYGO_APP_OMARCHY_FILE`.
Missing explicit/invalid config or selected theme fails, not silently ignored.

## Boundaries

Starter composition is illustrative/unapproved. Canonical color/role/R1
mechanism is snapshotted from design-system at a pinned commit and verified
byte-for-byte; no absolute checkout dependency. Product composition needs its
own [Studio](https://github.com/byteowlz/design-studio/blob/main/VISUAL-CONTRACT.md)
variant/revision/approval record. Vial's selected layout does not approve this
starter for other products.

No 3D engine is included. Consumers may add a native scene behind an operable
2D/table equivalent. Imported labels are untrusted data, never executable markup.
Readiness, build success or schema validation grants no mutation authority.
