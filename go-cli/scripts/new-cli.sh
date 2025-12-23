#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: new-cli.sh <name> [--path DIR] [--module MODULE]

Create a new CLI project by cloning the current template into DIR (defaults to <name>).

Options:
  -h, --help        Show this message
      --path DIR    Destination directory for the new project
      --module MOD  Go module path for the new project (defaults to <name>)
USAGE
}

die() {
  echo "new-cli.sh: $*" >&2
  exit 1
}

NAME=""
DEST=""
MODULE=""

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
    --module)
      shift
      [[ $# -gt 0 ]] || die "--module requires an argument"
      MODULE="$1"
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

if [[ -z "$MODULE" ]]; then
  MODULE="$NAME"
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

IGNORE = {'.git', 'dist', 'target', '.DS_Store', '__pycache__'}

def ignore(directory, contents):
    return IGNORE.intersection(contents)

shutil.copytree(root, dest, ignore=ignore)

for path in dest.rglob('new-cli.sh'):
    if path.is_file():
        path.chmod(path.stat().st_mode | 0o111)
PY

python3 - "$NAME" "$MODULE" "$DEST" <<'PY'
import pathlib
import sys

name = sys.argv[1]
module = sys.argv[2]
dest = pathlib.Path(sys.argv[3])

replacements = {
    "go-cli": name,
    "GO_CLI": name.upper().replace('-', '_'),
    "gitlab.cc-asp.fraunhofer.de/templates/go-cli": module,
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
PY

echo "Created CLI project at $DEST"
echo "Next steps:"
echo "  1. cd $DEST"
echo "  2. go mod tidy"
echo "  3. go run . -- --help"
