#!/usr/bin/env bash
# Register this webapp in the byteowlz service registry so it is discoverable
# fleet-wide. Call when the UI starts serving; unregister on shutdown.
# Env overrides: NAME, PORT, GROUP, DESC, AGENT, REGISTRY
set -euo pipefail

REGISTRY_URL="${REGISTRY:-http://100.64.0.12:8774}"
NAME="${NAME:-{{project_name}}}"
PORT="${PORT:-8790}"
GROUP="${GROUP:-misc}"   # pixelgym | voice | ops | infra | gpu | misc
DESC="${DESC:-{{project_name}} web UI}"
AGENT="${AGENT:-governor}"

IP="$(tailscale ip -4 2>/dev/null | head -1 || true)"
[ -z "${IP}" ] && IP="$(hostname -I 2>/dev/null | awk '{print $1}' || true)"
[ -z "${IP}" ] && IP="localhost"
URL="http://${IP}:${PORT}"

curl -fsS -X POST "${REGISTRY_URL}/register" \
	-H 'content-type: application/json' \
	-d "{\"name\":\"${NAME}\",\"url\":\"${URL}\",\"group\":\"${GROUP}\",\"desc\":\"${DESC}\",\"agent\":\"${AGENT}\"}"
echo
