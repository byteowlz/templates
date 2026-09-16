# Rust Workspace Template

A batteries-included Rust workspace template with CLI, TUI, MCP server, and HTTP API crates sharing a common core library.

## Quick Start

Install the latest stable Rust toolchain (`rustup default stable`), then:

```bash
cargo build
cargo test
```

Run individual binaries:

```bash
cargo run -p {{project_name}}-cli -- run
cargo run -p {{project_name}}-tui
cargo run -p {{project_name}}-api -- --port 3000
cargo run -p {{project_name}}-mcp
```

Scaffold a new project:

```bash
scripts/new-cli.sh my-app
```

```powershell
pwsh scripts/new-cli.ps1 my-app
```

This creates a new workspace with all crates renamed (e.g., `my-app-core`, `my-app-cli`, etc.).

## Workspace Structure

```
crates/
  {{project_name}}-core/    # Shared library: config, paths, error types
  {{project_name}}-cli/     # Command-line interface
  {{project_name}}-tui/     # Terminal user interface (ratatui)
  {{project_name}}-mcp/     # Model Context Protocol server
  {{project_name}}-api/     # HTTP API server (axum)
examples/
  config.toml   # Example configuration
scripts/
  new-cli.sh    # Unix scaffolding script
  new-cli.ps1   # PowerShell scaffolding script
```

## Crates

### {{project_name}}-core

Shared library providing:
- `AppConfig` - Configuration loading via `config` crate
- `AppPaths` - XDG-compliant path resolution
- Error types and common utilities

### {{project_name}}-cli

Command-line interface with:
- Subcommands: `run`, `init`, `config`, `completions`
- Global flags: `-q`, `-v`, `--debug`, `--trace`, `--json`, `--no-color`, `--dry-run`, `--yes`
- Shell completion generation

```bash
cargo run -p {{project_name}}-cli -- --help
cargo run -p {{project_name}}-cli -- completions bash > target/{{project_name}}-cli.bash
```

### {{project_name}}-tui

Terminal UI built with ratatui featuring:
- Three-pane layout (navigation, list, details)
- Vim-style navigation (j/k/h/l)
- Modal help system

```bash
cargo run -p {{project_name}}-tui
```

### {{project_name}}-mcp

MCP (Model Context Protocol) server exposing tools:
- `get_profile` - Current configuration profile
- `echo` - Echo messages
- `get_runtime_config` - Runtime configuration

```bash
cargo run -p {{project_name}}-mcp
```

### {{project_name}}-api

HTTP API server (axum) with non-sensitive endpoints:
- `GET /` - Service info
- `GET /health` - Health check

The default scaffold exposes **no configuration endpoint** and configures **no
CORS layer**; configuration (which may later hold secrets) is never served over
the wire. If a browser client needs cross-origin access, add an explicit,
exact-origin `CorsLayer` and keep it scoped to required origins.

```bash
cargo run -p {{project_name}}-api -- --port 3000
curl http://localhost:3000/health
```

## Configuration

Default config path: `$XDG_CONFIG_HOME/rust-workspace/config.toml`

Override with `--config <path>` or environment variables using the `RUST_WORKSPACE__` prefix:

```bash
RUST_WORKSPACE__LOGGING__LEVEL=debug cargo run -p {{project_name}}-cli -- run
```

See `examples/config.toml` for all options.

## Development

```bash
cargo fmt                                    # Format code
cargo clippy --all-targets --all-features   # Lint
sg scan --config .ast-grep/sgconfig.yml     # ast-grep Rust guardrails
cargo test                                   # Run tests
cargo build --release                        # Release build
```

## Scaffolding

The `scripts/new-cli.sh` (Unix) and `scripts/new-cli.ps1` (PowerShell) scripts create a new project from this template:

```bash
scripts/new-cli.sh my-app --path ~/projects/my-app
```

This will:
1. Copy the template to the destination
2. Rename all crates from `rust-*` to `my-app-*`
3. Update all references in Cargo.toml, source files, and documentation
4. Rename crate directories accordingly

Requirements: `python3` for the shell script, PowerShell 7 for the Windows script.
