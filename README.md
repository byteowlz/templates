# Byteowlz Project Templates

Project templates for use with `byt new`.

## Available Templates

| Template            | Description                                                                                                                       |
| ------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| `rust-cli`          | Rust CLI application with clap, config, XDG paths                                                                                 |
| `rust-workspace`    | Rust workspace with multiple crates                                                                                               |
| `python-cli`        | Python CLI with uv, typer, XDG paths                                                                                              |
| `go-cli`            | Go CLI with cobra, viper, XDG paths                                                                                               |
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

## Template Structure

Each template includes:

- `justfile` - Standard commands (build, test, install, etc.)
- `AGENTS.md` - AI agent instructions
- `CONTEXT.md` - Project domain glossary starter
- `docs/adr/` - Concise architecture decision guidance and template
- `.github/workflows/release.yml` - Automated releases
- Language-specific project files

## Customization

Fork this repo and set your custom template repo in `~/.config/byt/config.toml`:

```toml
[templates]
repo = "your-org/templates"
```

## Template Variables

Templates use `{{project_name}}` as a placeholder. This is replaced during scaffolding with the actual project name.
