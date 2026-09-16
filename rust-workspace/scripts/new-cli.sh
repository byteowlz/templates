#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: new-cli.sh <name> [--path DIR]

Create a new workspace project by cloning this template into DIR (defaults to <name>).

Renames all crates from rust-* to <name>-*, rewrites their code identifiers
({{project_name}}_core, {{project_name}}_cli, ... used in `use` statements), the environment prefix,
the workspace name, and every project-identity reference (README, config
paths, schema URLs, CI, justfile) from the template defaults to <name>.

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

# Order matters:
# 1. Replace the workspace name and env prefix first.
# 2. Rewrite each crate's hyphenated name (`{{project_name}}-core` -> `<name>-core`) and
#    its underscored code identifier (`{{project_name}}_core`) used in `use rust_*::...`
#    imports — these aliases are what byt-era scaffolding missed and left
#    uncompilable. Exact names only: `rust-magic-linter` (a lint preset) and
#    `dtolnay/rust-toolchain` (a CI action) must stay untouched.
# 3. Rewrite scaffolding placeholders that only `byt` understood.
replacements = [
    ("rust-workspace", name),
    ("RUST_WORKSPACE", upper),
    ("{{project_name}}-core", f"{name}-core"),
    ("{{project_name}}-cli", f"{name}-cli"),
    ("{{project_name}}-tui", f"{name}-tui"),
    ("{{project_name}}-mcp", f"{name}-mcp"),
    ("{{project_name}}-api", f"{name}-api"),
    ("{{project_name}}_core", f"{underscore}_core"),
    ("{{project_name}}_cli", f"{underscore}_cli"),
    ("{{project_name}}_tui", f"{underscore}_tui"),
    ("{{project_name}}_mcp", f"{underscore}_mcp"),
    ("{{project_name}}_api", f"{underscore}_api"),
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
    if path.name in ('new-cli.sh', 'new-cli.ps1', 'smoke-test.sh', 'drift-check.sh', 'privilege-boundary-test.sh'):
        return
    text = path.read_text()
    original = text
    for old, new in replacements:
        text = text.replace(old, new)
    if text != original:
        path.write_text(text)

def rename_directories(base: pathlib.Path) -> None:
    crates_dir = base / "crates"
    if not crates_dir.exists():
        return
    for crate_dir in sorted(crates_dir.iterdir(), reverse=True):
        if not crate_dir.is_dir():
            continue
        # New templates carry literal `{{project_name}}-*` crate dirs; older ones
        # carried `rust-*`. Rename either to `<name>-*`.
        if "{{project_name}}" in crate_dir.name:
            new_name = crate_dir.name.replace("{{project_name}}", name)
        elif crate_dir.name.startswith("rust-"):
            new_name = f"{name}-" + crate_dir.name[len("rust-"):]
        else:
            continue
        crate_dir.rename(crate_dir.parent / new_name)

for path in dest.rglob('*'):
    apply(path)

rename_directories(dest)
PY

echo "Created workspace project at $DEST"
echo "Crates renamed to: ${NAME}-core, ${NAME}-cli, ${NAME}-tui, ${NAME}-mcp, ${NAME}-api"
