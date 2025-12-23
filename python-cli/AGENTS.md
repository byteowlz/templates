# AGENTS.md

Guidance for coding agents working on this Python CLI template.

## Core Principles

- **Never publish** artifacts to public registries without explicit user approval.
- Prefer cohesive refactors over additive `FooV2` clones; update code in place.
- Target Windows 11, Linux, and macOS 14+ with identical behaviour.
- Keep file headers minimal—no author or timestamp banners.

## Python Workflow

- Use `uv` for dependency management (`uv add`, `uv sync`, `uv run`).
- Run `uv run ruff check .` and `uv run pytest` after significant edits; expand coverage only when the user requests broader suites.
- Favour standard library solutions before introducing third-party packages; ensure new dependencies keep the template lean.
- When touching concurrency or async flows, document invariants and honour cancellation signals.

## CLI Expectations

- Prefer subcommands for verbs and expose common global flags (`-q`, stacked `-v`, `--debug`, `--trace`).
- Support machine-readable modes via `--json/--yaml` and honour NO_COLOR/FORCE_COLOR.
- Offer `--dry-run`, `--yes/--force`, `--no-progress`, `--timeout`, and `--parallel` when operations warrant them.
- Keep `--help` responsive and generate completions from the Typer command tree.

## Configuration & Storage

- Follow XDG directory conventions when available: config at `$XDG_CONFIG_HOME/<app>/config.toml`, data at `$XDG_DATA_HOME/<app>`, state at `$XDG_STATE_HOME/<app>`.
- Expand `~` and environment variables in config paths.
- Ship a commented example under `examples/`, create a default config on first run, and load overrides via `PYTHON_CLI__*` environment variables.

## House Rules

- Do exactly what the user asks—no unsolicited files or docs.
- Keep README updates concise, emoji-free, and only when requested.
- Never commit secrets or sensitive paths; scrub logs before surfacing them.

## Justfile Commands

This project uses [just](https://github.com/casey/just) as a command runner. Run `just` to see available commands.

**Core commands:**
```bash
just              # Show available commands
just install      # Install as global tool (uv tool install)
just install-dev  # Install in development mode
just run          # Run the CLI
just test         # Run tests
just fmt          # Format code
just lint         # Run linter
just check        # Lint + test
```

**Development workflow:**
```bash
just test-cov     # Run tests with coverage
just fix          # Auto-fix lint issues
just typecheck    # Type check with mypy
just update       # Update dependencies
just build        # Build the package
```

Always run `just check` before committing significant changes.

## Issue Tracking (bd/beads)

Use `bd` for all issue tracking. Do NOT use markdown TODOs or external trackers.

```bash
bd ready --json                              # Find unblocked work
bd create "Title" -t task -p 2 --json        # Create issue (types: bug/feature/task/epic/chore)
bd update <id> --status in_progress --json   # Claim task
bd close <id> --reason "Done" --json         # Complete work
```

Priorities: 0=critical, 1=high, 2=medium (default), 3=low, 4=backlog

Always commit `.beads/issues.jsonl` with code changes.

## Memory System (byt/mmry)

Use `byt memory` to store and retrieve project knowledge. Memories auto-detect the current repo.

**Adding memories:**
```bash
byt memory add "Important decision or learning"              # Auto-detects current repo
byt memory add "Cross-repo architecture decision" --govnr    # Force govnr store
byt memory add "Specific insight" -c "architecture" -i 8     # With category and importance
```

**Searching memories:**
```bash
byt memory search "query"           # Search current repo's memories
byt memory search "query" --govnr   # Search cross-repo memories
byt memory search "query" --all     # Search ALL projects
```

**When to add memories:**
- Architecture decisions and their rationale
- Non-obvious solutions to tricky problems
- Integration patterns with other byteowlz repos
- Performance findings or benchmarks
- API contracts or breaking changes

**When to search memories:**
- Before starting work on a feature (check for prior decisions)
- When encountering unfamiliar code patterns
- When integrating with other repos (`byt memory search "query" --all`)

## Releases & Distribution

This project uses GitHub Actions for automated releases. See `.github/workflows/release.yml`.

**Creating a release:**
```bash
# Tag-based (automatic trigger)
git tag v1.0.0 && git push --tags

# Manual trigger via CLI
gh workflow run release.yml -f tag=v1.0.0
```

**What the workflow builds:**
- Python wheels (platform-independent)
- Source distribution (sdist)

Uses `uv build` for fast, reliable builds.

**Required secrets:**
- `PYPI_API_TOKEN` - API token for PyPI publishing (uncomment publish step in workflow)

**Installation methods (once published):**
```bash
# PyPI
pip install <package-name>
uv tool install <package-name>

# From GitHub release
uv pip install https://github.com/byteowlz/<repo>/releases/download/v1.0.0/<package>-1.0.0-py3-none-any.whl
```
