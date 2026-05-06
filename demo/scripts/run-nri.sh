#!/usr/bin/env bash
# run-nri.sh — collect one snapshot of HAProxy metrics via nri-haproxy
# Usage: ./demo/scripts/run-nri.sh [--raw]
#   --raw  print the full JSON without filtering

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
BIN="${BIN:-$REPO_ROOT/bin/nri-haproxy}"
STATS_URL="${STATS_URL:-http://localhost:8404/stats}"
CLUSTER_NAME="${CLUSTER_NAME:-demo-cluster}"
RAW="${1:-}"

if [[ ! -x "$BIN" ]]; then
  echo "ERROR: binary not found at $BIN"
  echo "Run 'make compile' from the repo root first."
  exit 1
fi

echo "================================================================"
echo " nri-haproxy metrics snapshot"
echo " Stats URL : $STATS_URL"
echo " Cluster   : $CLUSTER_NAME"
echo " Time      : $(date)"
echo "================================================================"
echo ""

JSON=$("$BIN" \
  --stats_url "$STATS_URL" \
  --ha_proxy_cluster_name "$CLUSTER_NAME" \
  2>/dev/null)

if [[ "$RAW" == "--raw" ]]; then
  echo "$JSON" | jq .
  exit 0
fi

# ── Pretty-print by entity type ──────────────────────────────────────────────

echo "=== FRONTENDS (HAProxyFrontendSample) ==="
echo "$JSON" | jq -r '
  .data[]
  | select(.entity.type == "ha-frontend")
  | "Entity: \(.entity.name)",
    (.metrics[0] | to_entries
      | map("  \(.key): \(.value)")
      | join("\n")),
    ""
'

echo ""
echo "=== BACKENDS (HAProxyBackendSample) ==="
echo "$JSON" | jq -r '
  .data[]
  | select(.entity.type == "ha-backend")
  | "Entity: \(.entity.name)",
    (.metrics[0] | to_entries
      | map("  \(.key): \(.value)")
      | join("\n")),
    ""
'

echo ""
echo "=== SERVERS (HAProxyServerSample) ==="
echo "$JSON" | jq -r '
  .data[]
  | select(.entity.type == "ha-server")
  | "Entity: \(.entity.name)",
    (.metrics[0] | to_entries
      | map("  \(.key): \(.value)")
      | join("\n")),
    ""
'

echo ""
echo "=== INVENTORY ==="
echo "$JSON" | jq -r '
  .data[]
  | select(.inventory | length > 0)
  | "Entity: \(.entity.name)",
    (.inventory | to_entries
      | map("  \(.key): \(.value | tostring)")
      | join("\n")),
    ""
'
