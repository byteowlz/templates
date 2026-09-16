#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: new-cli.sh <name> [--path DIR]

Create a new CLI project by cloning this template into DIR (defaults to <name>).

Renames the crate, its code identifier, environment prefix, and every
project-identity reference (README, config paths, schema URLs, CI, justfile)
from the template defaults to <name>. Rust imports use the underscored form,
so this script rewrites both `rust-cli` and `rust_cli`.

Options:
  -h, --help      Show this message
      --path DIR  Destination directory for the new project
USAGE
}

die() {
  echo "new-cli.sh: $*" >&2
  exit 1
}

NAME=""
DEST=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help)
      usage
      exit 0
      ;;
    --path)
      shift
      [[ $# -gt 0 ]] || die "--path requires an argument"
      DEST="$1"
      ;;
    -*)
      die "unknown option: $1"
      ;;
    *)
      if [[ -z "$NAME" ]]; then
        NAME="$1"
      else
        die "unexpected argument: $1"
      fi
      ;;
  esac
  shift
done

[[ -n "$NAME" ]] || die "project name is required"

if [[ ! "$NAME" =~ ^[a-zA-Z][a-zA-Z0-9_-]*$ ]]; then
  die "project name must start with a letter and contain only letters, numbers, '_' or '-'"
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_ROOT="$(dirname "$SCRIPT_DIR")"

if [[ -z "$DEST" ]]; then
  DEST="$(dirname "$TEMPLATE_ROOT")/$NAME"
fi

case "$DEST" in
  /*) ;; # absolute path
  *) DEST="$PWD/$DEST" ;;
esac

if [[ -e "$DEST" ]]; then
  die "destination already exists: $DEST"
fi

mkdir -p "$(dirname "$DEST")"

python3 - "$TEMPLATE_ROOT" "$DEST" <<'PY'
import pathlib
import shutil
import sys

root = pathlib.Path(sys.argv[1])
dest = pathlib.Path(sys.argv[2])

def ignore(directory, contents):
    ignored = {'.git', 'target', '.DS_Store'}
    return ignored.intersection(contents)

shutil.copytree(root, dest, ignore=ignore)

for path in dest.rglob('new-cli.sh'):
    if path.is_file():
        path.chmod(path.stat().st_mode | 0o111)
PY

python3 - "$NAME" "$DEST" <<'PY'
import pathlib
import sys

name = sys.argv[1]
dest = pathlib.Path(sys.argv[2])

underscore = name.replace('-', '_')
upper = name.upper().replace('-', '_')

# Order matters: replace the exact package name and env prefix before the
# shorter hyphenated shorthand, and replace the underscored crate identifier
# (used in `use rust_cli::...` imports) that byt-era scaffolding missed.
replacements = [
    ("rust-cli", name),
    ("rust_cli", underscore),
    ("RUST_CLI", upper),
    ("{{project_name}}", name),
    ("your-binary-name", name),
]

def is_text(path: pathlib.Path) -> bool:
    try:
        path.read_bytes().decode('utf-8')
        return True
    except (UnicodeDecodeError, OSError):
        return False

def apply(path: pathlib.Path) -> None:
    if not path.is_file() or not is_text(path):
        return
    # Do not rewrite the scaffolding/tooling scripts themselves: their embedded
    # replacement tables and token lists must stay intact so they remain reusable
    # after copying.
    if path.name in ('new-cli.sh', 'new-cli.ps1', 'smoke-test.sh', 'drift-check.sh'):
        return
    text = path.read_text()
    original = text
    for old, new in replacements:
        text = text.replace(old, new)
    if text != original:
        path.write_text(text)

for path in dest.rglob('*'):
    apply(path)
PY

echo "Created CLI project at $DEST"
