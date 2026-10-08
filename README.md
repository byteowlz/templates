# Byteowlz Project Templates

Project templates for use with `byt new`.

## Available Templates

| Template            | Description                                                                                                                       |
| ------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| `rust-cli`          | Rust CLI application with clap, config, XDG paths                                                                                 |
| `rust-workspace`    | Rust workspace with multiple crates                                                                                               |
| `rust-gpui`         | Native GPUI desktop app; byteowlz design system built in (one `gpui-kit` dep, `scheme -> closed role layer -> gpui theme set`)     |
| `python-cli`        | Python CLI with uv, typer, XDG paths                                                                                              |
| `go-cli`            | Go CLI with cobra, viper, XDG paths                                                                                               |
| `go-mygo`           | Fully native MyGo v0.3.2 GPU UI + canonical Base24/roles/R1 adapter, XDG TOML and shared semantic CLI; no HTML/webview |
| `webapp`            | Focused webapp: Vite + React + TS + Tailwind v4, byteowlz design-system tokens, optional omarchy theme follow + oqto app manifest |
| `experiment-review` | Private, mobile-first curated comparison gallery (pxltr visual grammar, synthetic samples only)                                   |

## Usage

```bash
# Create a new project from template
byt new myproject --template rust-cli

# Create and also init GitHub repo
byt new myproject --template rust-cli --github

# On byt versions without the named experiment-review template:
byt new myreview --from-git byteowlz/templates --subdir experiment-review
```

`go-mygo` uses Go-only native `ui.View`; the unselected system-webview prototype
is preserved on `feat/mygo-webview-prototype`, not as the default template.
Its fixture agent gate is evidence-scoped, not a consumer readiness certificate.
Native signing/release configuration is deliberately left to consumers.

## Template Structure

Each template includes:

- `justfile` - Standard commands (build, test, install, etc.)
- `AGENTS.md` - AI agent instructions
- `CONTEXT.md` - Project domain glossary starter
- `docs/adr/` - Concise architecture decision guidance and template
- `.github/workflows/release.yml` - Automated releases
- Language-specific project files

## Project Icon

Projects with a visual identity add a tool-facing mark at a fixed path (see the
byteowlz repository standard in the wiki):

```text
icon/icon_on_dark.svg    # light-coloured mark, for dark backgrounds
icon/icon_on_light.svg   # dark-coloured mark, for light backgrounds
```

Square, transparent, no text, legible at 16 px. Run `just icon` to render the
256x256 PNGs next to them, and commit all four files. Templates do not ship a
placeholder icon: without `icon/`, tools show nothing.

## Customization

Fork this repo and set your custom template repo in `~/.config/byt/config.toml`:

```toml
[templates]
repo = "your-org/templates"
```

## Template Variables

Templates use `{{project_name}}` as a placeholder. This is replaced during scaffolding with the actual project name.

## Conventions

- [OSC7501.md](OSC7501.md) — Adopt terminal program status where useful; headless apps use structured activity channels. AGENT_CTX supplies injected runtime attribution, not live state. Contract development: standalone `wismut/agent-ctx` on Forgejo (local `../agent-ctx`); v3 is draft.
