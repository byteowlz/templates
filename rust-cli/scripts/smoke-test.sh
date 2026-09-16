#!/usr/bin/env bash
#
# Smoke test for the scaffold (tmpl-d9wd).
#
# Renders the template to a non-default name, asserts that every project-identity
# reference (crate name, underscore import, env prefix, config path, schema URL,
# scaffolding placeholders) was substituted consistently, then compiles and tests
# the generated project so a renamed scaffold is proven to build.
#
# A previous bug left `use rust_cli::...` imports and `rust-cli` identity
# references in place, so renamed projects failed to compile. This test guards
# against a regression in scripts/new-cli.sh.

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
    'rust-cli|rust-workspace|rust_core|rust_cli|RUST_CLI|RUST_WORKSPACE|\{\{project_name\}\}|your-binary-name' \
    "$DEST" \
    --exclude=new-cli.sh --exclude=new-cli.ps1 \
    --exclude=smoke-test.sh --exclude=drift-check.sh || true
)
if [[ -n "$STALE" ]]; then
  echo "ERROR: stale template identity found in generated project:" >&2
  echo "$STALE" >&2
  exit 1
fi

echo "==> Assert crates/package identity actually present"
# Sanity: the new name must appear in the manifest and (for the workspace) the
# underscored import must reference the renamed core crate.
grep -q "$NAME" "$DEST/Cargo.toml" || { echo "ERROR: new name missing from Cargo.toml" >&2; exit 1; }

echo "==> Compile generated project"
(
  cd "$DEST"
  cargo check --workspace
)

echo "==> Run generated project tests"
(
  cd "$DEST"
  cargo test --workspace --lib --bins --examples
)

echo "PASS: scaffold '$NAME' substitutes identity and builds/tests cleanly"
