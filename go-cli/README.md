# Go CLI Template

This repository provides a batteries-included starting point for building cross-platform Go CLIs. It mirrors the ergonomics of the Rust template while leaning on Cobra and Viper for command parsing and configuration.

## Quick Start

- Install Go 1.25 or newer.
- Fetch dependencies and verify the build:

  ```bash
  go test ./...
  ```

- Run the CLI in place:

  ```bash
  go run . -- run
  ```

- Scaffold a fresh project from this template:

  ```bash
  scripts/new-cli.sh my-cli
  ```

  ```powershell
  pwsh scripts/new-cli.ps1 my-cli
  ```

## Features

- Cobra-powered command interface with shared global flags (`-q`, `-v`, `--debug`, `--trace`, `--json`, `--yaml`, `--no-color`, `--dry-run`, `--yes`).
- Viper-based configuration loader that creates `$XDG_CONFIG_HOME/go-cli/config.toml` (or platform equivalents) on first run.
- Environment overrides using the `GO_CLI__*` prefix; e.g. `GO_CLI__LOGGING__LEVEL=debug`.
- Configurable data and state directories that honor XDG locations on Unix and the appropriate directories on Windows.
- Shell completion generation via `go run . -- completions <shell>`.
- Lightweight structured logging with color-aware console output and optional log file mirroring.
- `scripts/new-cli.sh` to clone the template with a new module name and paths.

## CLI Overview

```bash
go run . -- --help
```

Key subcommands:

- `run [TASK]` – executes the primary workflow with optional profile overrides.
- `init` – creates or refreshes the config file (use `--force` or `--yes` to overwrite).
- `config show|path|reset` – inspects the effective configuration.
- `completions <shell>` – emits shell completions to stdout (`bash`, `zsh`, `fish`, `powershell`).

Global flags apply to every subcommand, enabling quiet mode, stacked verbosity (`-vv`), trace logging, dry runs, JSON/YAML output, color control, progress suppression, and timeouts.

## Configuration

- Default config path: `$XDG_CONFIG_HOME/go-cli/config.toml` (or `%APPDATA%\go-cli\config.toml` on Windows). Override with `--config <path>`.
- Sample configuration with inline comments is available at `examples/config.toml`.
- Data and state directories default to `$XDG_DATA_HOME/go-cli` and `$XDG_STATE_HOME/go-cli` (falling back to `~/.local/share` and `~/.local/state` when unset). Override inside the config file.
- Values support `~` expansion and environment variables (e.g. `$HOME/logs/app.log`).

## Development Workflow

- Format the codebase:

  ```bash
  gofmt -w .
  ```

- Run the test suite:

  ```bash
  go test ./...
  ```

- Recommended lint pass during active development:

  ```bash
  golangci-lint run ./...
  ```

- Generate completions for your shell:

  ```bash
  go run . -- completions bash > dist/go-cli.bash
  ```

## Scaffold New Projects

- Run `scripts/new-cli.sh my-cli` (Unix shells) or `pwsh scripts/new-cli.ps1 my-cli` (Windows/PowerShell) to copy the template into `./my-cli` with all configuration files updated to the new module name.
- Provide `--path /some/where` (or `-Path C:\work\my-cli`) to choose a different destination directory.
- Requirements: `python3` for the shell script, PowerShell 7 (`pwsh`) for the Windows script.

## Project Structure

- `cmd/` – Cobra commands and CLI wiring.
- `internal/app/` – runtime context, configuration loaders, and command handlers.
- `examples/config.toml` – commented configuration template.
- `go.mod` – dependencies and metadata for the template module.

Feel free to fork this template and tailor the commands, config schema, or runtime behavior to your project's needs.
