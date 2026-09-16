# AGENTS.md

Guidance for coding agents working on this Rust CLI template. Kept short: it
carries enforceable boundaries and workflows and points at machine-checked
configuration instead of duplicating any inventory.

## Source of truth

Static facts are machine-checked, not copied here:

- Crate/dependency/lint settings → `Cargo.toml`; authoritative enumeration is
  `cargo metadata --no-deps --format-version 1`.
- Task commands → `just` (run `just` to list).
- Issues → `trx` (see below).
- Drift guard → `scripts/drift-check.sh` (verifies commands exist, the
  JSON/TOML-only constraint, and that documented version claims match the real
  manifest). Run it after touching the manifest or any doc claim.

If `cargo metadata`, `just`, or `scripts/drift-check.sh` disagree with anything
here, the command is right and this file is wrong.

## Domain and architecture

- Read `CONTEXT.md` before domain work; keep it a glossary only.
- Read `docs/adr/` before architectural changes; add an ADR only for hard-to-reverse
  decisions with a real trade-off.
- Never publish to a public registry without explicit user approval.

## Strict lints

`[lints.clippy]` is maximum-strictness: `unsafe_code = "forbid"`;
`unwrap_used`, `expect_used`, `panic`, `todo`, `unimplemented`, `dbg_macro`,
`exit` = deny. Propagate errors with `?`, `anyhow::Result`, `.context("...")`.
Output macros are allowed for a CLI.

## Workflow

- CLI: subcommands for verbs; `{{project_name}}` global flags `-q`, `-v`,
  `--debug`, `--trace`, `--no-color`, `--dry-run`, `--yes`, `--no-progress`.
- Config structs: after editing `src/config.rs`, run `just generate-config` and
  `just test` (the `validate_examples_are_up_to_date` test enforces it).
- Before anything significant: `just check-all`.

## Application formats: JSON and TOML only

- Machine output and config are **JSON or TOML only** — `--json` plus TOML
  config files. Never add a YAML output mode, a YAML example, or a YAML crate.
- The `config` crate runs with `default-features = false` and only `json`/`toml`
  features. `scripts/drift-check.sh` enforces the absence of `serde_yaml`.

## Configuration & storage

- XDG paths with sensible fallbacks; expand `~` and env vars; ship a commented
  example under `examples/`; write a default config on first run; override via
  the `config` crate. Env prefix derives from `APP_NAME` — reference
  `env_prefix()`, don't hardcode a literal.

## Issue tracking (trx)

Use `trx` for all issue tracking — never markdown TODOs or `.beads`.

```bash
trx ready --json                                   # find unblocked work
trx create "Title" -t task -p 2 --json             # create (bug/feature/task/epic/chore)
trx update <id> --status in_progress --json        # claim
trx close <id> -r "reason" --json                  # complete with reason
```

Priorities: 0=critical, 1=high, 2=medium (default), 3=low, 4=backlog.
Issue state lives in `.trx/` (JSONL) — commit it with code changes.

## House rules

- Do exactly what the user asks — no unsolicited files.
- Keep README updates concise and emoji-free.
- Never commit secrets or sensitive paths; scrub logs.
- `Cargo.lock` is committed; bump manifest + lock together.
