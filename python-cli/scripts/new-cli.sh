#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: new-cli.sh <name> [--path DIR] [--package MODULE]

Create a new CLI project by cloning the current template into DIR (defaults to <name>).

Options:
  -h, --help         Show this message
      --path DIR     Destination directory for the new project
      --package MOD  Python package name for the project (defaults to <name> with '-' -> '_')
USAGE
}

die() {
  echo "new-cli.sh: $*" >&2
  exit 1
}

NAME=""
DEST=""
PACKAGE=""

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
    --package)
      shift
      [[ $# -gt 0 ]] || die "--package requires an argument"
      PACKAGE="$1"
      ;;
    -* )
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

if [[ -z "$PACKAGE" ]]; then
  PACKAGE="${NAME//-/_}"
fi

if [[ ! "$PACKAGE" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]]; then
  die "package name must start with a letter or underscore and contain only letters, numbers, or '_'"
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

IGNORE = {'.git', '.ruff_cache', '.mypy_cache', '.pytest_cache', '.venv', 'dist', 'build', '__pycache__'}

DefIgnore = shutil.ignore_patterns(*IGNORE)

shutil.copytree(root, dest, ignore=DefIgnore)

for path in dest.rglob('*.sh'):
    if path.is_file():
        path.chmod(path.stat().st_mode | 0o111)
PY

python3 - "$NAME" "$PACKAGE" "$DEST" <<'PY'
import pathlib
import sys

name = sys.argv[1]
package = sys.argv[2]
dest = pathlib.Path(sys.argv[3])
upper = name.upper().replace('-', '_')

replacements = {
    'python-cli': name,
    'python_cli': package,
    'PYTHON_CLI': upper,
}

for path in dest.rglob('*'):
    if not path.is_file():
        continue
    try:
        text = path.read_text()
    except UnicodeDecodeError:
        continue
    original = text
    for old, new in replacements.items():
        if old in text:
            text = text.replace(old, new)
    if text != original:
        path.write_text(text)

package_dir = dest / 'python_cli'
if package_dir.exists():
    package_dir.rename(dest / package)
PY

echo "Created CLI project at $DEST"
echo "Next steps:"
echo "  1. cd $DEST"
echo "  2. uv sync"
echo "  3. uv run $NAME -- --help"
