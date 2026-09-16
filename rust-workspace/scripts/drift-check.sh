#!/usr/bin/env bash
#
# Drift check (tmpl-8adx): verify that the facts this repo documents
# (commands, crate/version references, JSON/TOML-only constraint) match the
# machine-checked configuration (the actual Cargo.toml manifests and PATH).
#
# Anything this script checks is "source of truth" in AGENTS.md — if it fails,
# either the code or the documentation drift must be fixed.

set -u
FAILED=0

echo "==> Required commands on PATH"
for cmd in cargo just trx python3; do
  if command -v "$cmd" >/dev/null 2>&1; then
    echo "  ok: $cmd -> $(command -v "$cmd")"
  else
    echo "  FAIL: required command not found on PATH: $cmd"
    FAILED=1
  fi
done

echo "==> No YAML application format (JSON/TOML only)"
if command -v rg >/dev/null 2>&1 && rg -l 'serde_yaml|serde-yaml|yaml' --glob 'Cargo.toml' . >/dev/null 2>&1; then
  echo "  FAIL: found 'serde_yaml'/'yaml' in a Cargo.toml"
  rg -n 'serde_yaml|serde-yaml' --glob 'Cargo.toml' . || true
  FAILED=1
else
  if grep -rl 'serde_yaml' --include='Cargo.toml' . >/dev/null 2>&1; then
    echo "  FAIL: found 'serde_yaml' in a Cargo.toml"
    grep -rn 'serde_yaml' --include='Cargo.toml' . || true
    FAILED=1
  else
    echo "  ok: no serde_yaml dependency in any manifest"
  fi
fi

echo "==> Workspace config crate: default-features=false, json/toml only"
python3 - "$FAILED" <<'PY'
import re, sys
from pathlib import Path

failed = int(sys.argv[1])
root = Path('.')

# Collect the two relevant Cargo.toml files: the workspace root and any crate
# manifest that lists a `config` dependency.
manifests = list(root.rglob('Cargo.toml'))
found_config = False
for m in manifests:
    text = m.read_text()
    if 'config' not in text:
        continue
    for line in text.splitlines():
        s = line.strip()
        if s.startswith('config') and ('version' in s or 'features' in s):
            found_config = True
            if 'default-features = false' not in s:
                print(f"  FAIL: {m}: config dep missing default-features = false")
                failed = 1
            if 'yaml' in s:
                print(f"  FAIL: {m}: config dep enables yaml")
                failed = 1
            if 'json' not in s or 'toml' not in s:
                print(f"  WARN: {m}: config features do not mention both json and toml: {s}")

if found_config and failed == 0:
    print("  ok: config crate is default-features=false with json/toml only")

sys.exit(failed)
PY
FAILED=$?

echo "==> rmcp version references consistent with Cargo.toml"
python3 - <<'PY'
import re, sys
from pathlib import Path

root = Path('.')
ws = root / 'Cargo.toml'
text = ws.read_text()
m = re.search(r'^rmcp\s*=\s*\{[^}]*version\s*=\s*"([^"]+)"', text, re.M)
if m:
    actual = m.group(1)
    major = actual.split('.')[0]
    print(f"  actual rmcp version in Cargo.toml: {actual}")
    bad = []
    for md in list(root.rglob('*.md')):
        for line in md.read_text().splitlines():
            for mm in re.finditer(r'rmcp\s+(\d+)\.\d+', line, re.I):
                if mm.group(1) != major:
                    bad.append(f"{md}: {line.strip()}")
    if bad:
        print("  FAIL: documented rmcp version contradicts Cargo.toml:")
        for b in bad:
            print("    " + b)
        sys.exit(1)
    else:
        print("  ok: no stale rmcp version claims in docs")
else:
    print("  io: no explicit rmcp dep found (nothing to check)")
sys.exit(0)
PY
RC=$?
FAILED=$((FAILED == 0 ? RC : FAILED))

echo "==> Workspace resolves (cargo metadata)"
if cargo metadata --no-deps --format-version 1 >/dev/null 2>&1; then
  echo "  ok: cargo metadata resolves"
else
  echo "  FAIL: cargo metadata --no-deps does not resolve"
  FAILED=1
fi

echo
if [[ "$FAILED" -eq 0 ]]; then
  echo "drift-check: PASS"
else
  echo "drift-check: FAIL"
  exit 1
fi
