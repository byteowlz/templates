# Python CLI Template

This repository provides a batteries-included starting point for building cross-platform Python CLIs using Typer, Rich, and uv for tooling.

## Quick Start

- Install [uv](https://github.com/astral-sh/uv) and ensure Python 3.12+ is available.
- Bootstrap the environment and run the test matrix:

  ```bash
  uv sync
  uv run pytest
  ```

- Run the CLI in place:

  ```bash
  uv run python-cli -- run
  ```

- Scaffold a fresh project from this template:

  ```bash
  scripts/new-cli.sh my-cli
  ```

  ```powershell
  pwsh scripts/new-cli.ps1 -Name my-cli
  ```

## Features

- Typer-powered command interface with shared global flags (`-q`, `-v`, `--debug`, `--trace`, `--json`, `--yaml`, `--no-color`, `--dry-run`, `--yes`).
- Structured configuration loader that creates `$XDG_CONFIG_HOME/python-cli/config.toml` (or platform equivalents) on first run.
- Environment overrides using the `PYTHON_CLI__*` prefix; e.g. `PYTHON_CLI__LOGGING__LEVEL=debug`.
- Configurable data and state directories that honor XDG locations on Unix and the appropriate directories on Windows.
- Shell completion generation via `uv run python-cli -- completions <shell>`.
- Rich-powered logging and console output with JSON/YAML renderers for automation scenarios.
- `scripts/new-cli.sh` to clone the template with a new package name and CLI binary.

## CLI Overview

```bash
uv run python-cli -- --help
```

Key subcommands:

- `run [TASK]` – executes the primary workflow with optional profile overrides.
- `init` – creates or refreshes the config file (use `--force` or `--yes` to overwrite).
- `config show|path|reset` – inspects the effective configuration.
- `completions <shell>` – emits shell completions to stdout (`bash`, `zsh`, `fish`, `powershell`, `elvish`).

Global flags apply to every subcommand, enabling quiet mode, stacked verbosity (`-vv`), trace logging, dry runs, JSON/YAML output, color control, progress suppression, and timeouts.

## Configuration

- Default config path: `$XDG_CONFIG_HOME/python-cli/config.toml` (or `%APPDATA%\python-cli\config.toml` on Windows). Override with `--config <path>`.
- Sample configuration with inline comments is available at `examples/config.toml`.
- Data and state directories default to `$XDG_DATA_HOME/python-cli` and `$XDG_STATE_HOME/python-cli` (falling back to `~/.local/share` and `~/.local/state` when unset). Override inside the config file.
- Values support `~` expansion and environment variables via the `PYTHON_CLI__*` prefix.

## Development Workflow

- Format and lint the codebase:

  ```bash
  uv run ruff check .
  uv run ruff format .
  ```

- Run the test suite:

  ```bash
  uv run pytest
  ```

- Optional static type checking:

  ```bash
  uv run mypy python_cli
  ```

- Generate completions for your shell:

  ```bash
  uv run python-cli -- completions bash > dist/python-cli.bash
  ```

## Scaffold New Projects

- Run `scripts/new-cli.sh my-cli` (Unix shells) or `pwsh scripts/new-cli.ps1 -Name my-cli` (Windows/PowerShell) to copy the template into `./my-cli` with all configuration files updated to the new package name.
- Provide `--path /some/where` (or `-Path C:\work\my-cli`) to choose a different destination directory, and `--package my_cli` (or `-Package my_cli`) to specify a custom import package.
- Requirements: uv for dependency management and Python 3.12+.

## Project Structure

- `python_cli/main.py` – Typer application and command wiring.
- `python_cli/config.py` – configuration defaults, loaders, and environment overrides.
- `python_cli/runtime.py` – runtime context, logging setup, and shared helpers.
- `python_cli/handlers.py` – command handlers for `run`, `init`, and `config` operations.
- `examples/config.toml` – commented configuration template.
- `pyproject.toml` – project metadata and dependency declarations for uv.

Feel free to fork this template and tailor the commands, config schema, or runtime behavior to your project's needs.
