# AGENTS.md

Guidance for coding agents working on this Rust workspace template. This file is
deliberately short: it keeps the enforceable boundaries and workflows, and points
at machine-checked configuration instead of duplicating any inventory.

## Source of truth

Static facts about this repo are machine-checked, **not** copied here:

- Crate list, dependency versions, lint settings → the workspace `Cargo.toml`.
  The authoritative enumeration is `cargo metadata --no-deps --format-version 1`.
- Task commands → `just` (run `just` to list).
- Issues → `trx` (see "Issue tracking" below).
- Drift guard → `scripts/drift-check.sh` (verifies commands exist and that
  documented crate/version references and the JSON/TOML-only constraint match
  the real manifests). Run it after touching any manifest or doc claim.

If `cargo metadata`, `just list`, or `scripts/drift-check.sh` disagree with
anything written below, the command is right and this file is wrong.

## Domain and architecture

- Read `CONTEXT.md` before domain work; keep it a glossary only.
- Read `docs/adr/` before architectural changes; add an ADR only for hard-to-reverse
  decisions with a real trade-off (`docs/adr/README.md`).
- Layout: `crates/{{project_name}}-core` is the only library crate and the
  dependency root; the four binaries (`-cli`, `-tui`, `-mcp`, `-api`) depend on
  it but never on each other (`cargo metadata` is the map).
- Never publish to a public registry without explicit user approval.

## Strict lints

`[workspace.lints.clippy]` is maximum-strictness. Key constraints:

- `unsafe_code = "forbid"`; `unwrap_used`, `expect_used`, `panic`, `todo`,
  `unimplemented`, `dbg_macro`, `exit` = deny — propagate with `?`,
  `anyhow::Result`, `.context("...")`.
- Output macros (`print_*`) are allowed for CLIs/TUIs/APIs.

## Workflow

- Domain terms: use `CONTEXT.md`.
- Add code: follow the crate's dominant pattern (CLI subcommands, MCP tools,
  API handlers, TUI screens).
- Config structs: after editing `crates/{{project_name}}-core/src/config.rs`, run
  `just generate-config` (regenerates `examples/config.toml` + `.schema.json`);
  the `validate_examples_are_up_to_date` test fails if you forget.
- Before committing anything significant: `just check-all` (fmt + clippy + test).
- Run `scripts/drift-check.sh` after touching manifests or versioned claims.

## Application formats: JSON and TOML only

- Configuration and machine output are **JSON or TOML only**. Never add a YAML
  output mode, a YAML dependency, or a YAML example.
- Root `config` crate runs with `default-features = false` and only the
  `json`/`toml` features — keep it that way (no `yaml` feature).
- Do **not** introduce `serde_yaml` or any YAML crate. `scripts/drift-check.sh`
  enforces this in CI-style runs.

## HTTP API service ({{project_name}}-api)

- The default scaffold exposes only non-sensitive endpoints (`GET /` and
  `GET /health`) and configures **no CORS layer** and **no config endpoint**.
- Never expose configuration (which may hold secrets) over the wire. If CORS is
  genuinely required, opt into an exact-origin `CorsLayer`, never `Any`.
- Privilege boundaries and secret migration follow the rules in
  `scripts/privilege-boundary-test.sh` (separate executors, explicit modes,
  ancestor ACLs, credential rotation) and its guardrail section below.

### Privilege-boundary guardrails (services)

When a service creates files, holds credentials, or migrates secrets:

1. **Separate executors**: privilege-lowering or -raising work runs in a distinct
   executor/process whose ambient privileges are explicit — never silently in the
   calling request thread.
2. **Ownership/ACL, including ancestor dirs**: set and assert ownership and ACL
   on the target *and every ancestor* it creates; a writable-by-others parent
   defeats a locked file.
3. **Explicit file modes, independent of umask**: create with an explicit mode
   (`open(..., 0o600)`, `std::os::unix::fs::OpenOptionsExt::mode`, or
   `set_permissions`) so a permissive umask cannot widen it. Never rely on the
   process umask.
4. **Credential rotation during privilege migration**: when a secret's location
   or owner changes, write the new value under a lock before removing the old,
   and destroy the old copy (don't just move the filename).
5. **Negative cross-user tests**: assert that an unprivileged/other user *cannot*
   read or write the protected resource (403/EACCES), not just that the owner can.
6. **Distinguish unit/API checks from OS enforcement**: in-process checks (role,
   route) are unit-level; ownership/ACL/mode are OS-level and must be asserted
   with real filesystem integration tests, not mocks.

`scripts/privilege-boundary-test.sh` is the reusable acceptance harness for (2)–(5).

## MCP ({{project_name}}-mcp)

- Uses `rmcp`; version is **pinned and documented in the root `Cargo.toml`**
  (which also pins `rmcp-macros` in lockstep — bump both together). Read the
  actual version there; never trust a version written in prose.
- `#[tool_router]` on the impl auto-registers tools; `ServerInfo` is
  non-exhaustive (mutate fields after `default()`).

## Configuration & storage

- XDG paths; expand `~` and env vars; ship a commented example under `examples/`,
  write a default on first run, override via the `config` crate.
- Env override prefix is derived from `APP_NAME` in `{{project_name}}-core` as
  `upper(PKG_NAME).replace('-', "_")`. Reference `APP_NAME`/`env_prefix()`
  rather than hardcoding a literal string.

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
- `Cargo.lock` is committed; bump manifests + lock together.
