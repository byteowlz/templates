#!/usr/bin/env bash
#
# Smoke test for the workspace scaffold (tmpl-d9wd).
#
# Renders the template to a non-default name, asserts that every project-identity
# reference (workspace name, per-crate hyphenated names, underscore `use`
# imports, env prefix, config path, schema URL, scaffolding placeholders) was
# substituted consistently, then compiles and tests the generated workspace so a
# renamed scaffold is proven to build.
#
# A previous bug left crate directories and `use {{project_name}}_core::...` imports
# un-renamed while renaming the rest, so generated projects failed to compile.
# This test guards against a regression in scripts/new-cli.sh.

set -euo pipefail

NAME="${SMOKE_NAME:-smokey-app}"
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_ROOT="$(dirname "$SCRIPT_DIR")"
DEST="$TMPDIR/$NAME"

echo "==> Render template as '$NAME'"
"$SCRIPT_DIR/new-cli.sh" "$NAME" --path "$DEST"

echo "==> Assert identity substitution (no stale template tokens)"
# The scaffolding scripts intentionally retain the template's literal crate
# names in their own replacement tables, so exclude them from the scan.
STALE=$(
  grep -rInE \
    'rust-workspace|{{project_name}}_cli|{{project_name}}_core|{{project_name}}_tui|{{project_name}}_mcp|{{project_name}}_api|RUST_WORKSPACE|\{\{project_name\}\}|your-binary-name' \
    "$DEST" \
    --exclude=new-cli.sh --exclude=new-cli.ps1 \
    --exclude=smoke-test.sh --exclude=drift-check.sh --exclude=privilege-boundary-test.sh || true
)
if [[ -n "$STALE" ]]; then
  echo "ERROR: stale template identity found in generated workspace:" >&2
  echo "$STALE" >&2
  exit 1
fi

echo "==> Assert crate layout + renamed imports actually present"
grep -q "$NAME" "$DEST/Cargo.toml" || { echo "ERROR: new name missing from workspace Cargo.toml" >&2; exit 1; }
UNDERSCORE="${NAME//-/_}"
if [[ ! -d "$DEST/crates/${UNDERSCORE}-core" && ! -d "$DEST/crates/$NAME-core" ]]; then
  echo "ERROR: core crate directory not renamed to '<name>-core'" >&2
  ls -1 "$DEST/crates" >&2
  exit 1
fi
if ! grep -rq "use ${UNDERSCORE}_core::" "$DEST/crates"; then
  echo "ERROR: expected renamed import 'use ${UNDERSCORE}_core::' not found" >&2
  exit 1
fi

echo "==> Compile generated workspace"
(
  cd "$DEST"
  cargo check --workspace
)

echo "==> Run generated workspace tests"
(
  cd "$DEST"
  cargo test --workspace --lib --bins --examples
)

echo "PASS: scaffold '$NAME' substitutes identity and builds/tests cleanly"
