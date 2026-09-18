#!/usr/bin/env bash
# Unregister this webapp from the byteowlz service registry. Call on shutdown.
# Env overrides: NAME, PORT, REGISTRY
set -euo pipefail

REGISTRY_URL="${REGISTRY:-http://100.64.0.12:8774}"
NAME="${NAME:-{{project_name}}}"
PORT="${PORT:-8790}"

IP="$(tailscale ip -4 2>/dev/null | head -1 || true)"
[ -z "${IP}" ] && IP="$(hostname -I 2>/dev/null | awk '{print $1}' || true)"
[ -z "${IP}" ] && IP="localhost"
URL="http://${IP}:${PORT}"

curl -fsS -X POST "${REGISTRY_URL}/unregister" \
	-H 'content-type: application/json' \
	-d "{\"name\":\"${NAME}\",\"url\":\"${URL}\"}" \
	|| echo "unregister failed (registry unreachable?)"
echo
