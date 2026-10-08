#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
# The uninstantiated module name is intentionally invalid until byt replaces it.
# stdlib-only instantiator is run outside module mode; no global Go setting changes.
GOTOOLCHAIN=go1.27.1 GO111MODULE=off go run ./scripts/instantiate.go .
