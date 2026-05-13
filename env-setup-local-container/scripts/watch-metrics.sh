#!/usr/bin/env bash
# watch-metrics.sh — poll nri-haproxy every N seconds and show a live metric
# summary table. Good for watching metrics change as traffic flows.
# Usage: ./env-setup-local-container/scripts/watch-metrics.sh [interval_seconds]

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
BIN="${BIN:-$REPO_ROOT/bin/nri-haproxy}"
STATS_URL="${STATS_URL:-http://localhost:8404/stats}"
CLUSTER_NAME="${CLUSTER_NAME:-demo-cluster}"
INTERVAL="${1:-5}"

if [[ ! -x "$BIN" ]]; then
  echo "ERROR: binary not found at $BIN — run 'make compile' first."
  exit 1
fi

echo "Watching nri-haproxy every ${INTERVAL}s  (Ctrl+C to stop)"
echo ""

while true; do
  JSON=$("$BIN" \
    --stats_url "$STATS_URL" \
    --ha_proxy_cluster_name "$CLUSTER_NAME" \
    2>/dev/null)

  clear
  echo "================================================================"
  echo " nri-haproxy live snapshot   $(date)"
  echo "================================================================"

  printf "\n%-48s %8s %8s %8s %8s\n" "ENTITY" "req/s" "200/s" "500/s" "sessions"
  printf "%-48s %8s %8s %8s %8s\n"   "------" "-----" "-----" "-----" "--------"

  # Frontends — keys are frontend.*
  echo "$JSON" | jq -r '
    .data[]
    | select(.entity.type == "ha-frontend")
    | [
        .entity.name,
        (.metrics[0]["frontend.httpRequestsPerSecond"]   // 0 | tostring),
        (.metrics[0]["frontend.http200ResponsesPerSecond"] // 0 | tostring),
        (.metrics[0]["frontend.http500ResponsesPerSecond"] // 0 | tostring),
        (.metrics[0]["frontend.currentSessions"]          // 0 | tostring)
      ]
    | @tsv
  ' | awk -F'\t' '{printf "%-48s %8s %8s %8s %8s\n", $1, $2, $3, $4, $5}'

  # Backends — keys are backend.*
  echo "$JSON" | jq -r '
    .data[]
    | select(.entity.type == "ha-backend")
    | [
        .entity.name,
        (.metrics[0]["backend.httpRequestsPerSecond"]    // 0 | tostring),
        (.metrics[0]["backend.http200ResponsesPerSecond"] // 0 | tostring),
        (.metrics[0]["backend.http500ResponsesPerSecond"] // 0 | tostring),
        (.metrics[0]["backend.currentSessions"]           // 0 | tostring)
      ]
    | @tsv
  ' | awk -F'\t' '{printf "%-48s %8s %8s %8s %8s\n", $1, $2, $3, $4, $5}'

  # Servers — keys are server.*
  echo "$JSON" | jq -r '
    .data[]
    | select(.entity.type == "ha-server")
    | [
        .entity.name,
        (.metrics[0]["server.sessionsPerSecond"]          // 0 | tostring),
        (.metrics[0]["server.http200ResponsesPerSecond"]  // 0 | tostring),
        (.metrics[0]["server.http500ResponsesPerSecond"]  // 0 | tostring),
        (.metrics[0]["server.currentSessions"]            // 0 | tostring)
      ]
    | @tsv
  ' | awk -F'\t' '{printf "%-48s %8s %8s %8s %8s\n", $1, $2, $3, $4, $5}'

  echo ""
  printf "%-48s %8s %8s %8s\n" "ENTITY" "avg_rsp_ms" "avg_q_ms" "status"
  printf "%-48s %8s %8s %8s\n" "------" "----------" "--------" "------"

  # Backend latency
  echo "$JSON" | jq -r '
    .data[]
    | select(.entity.type == "ha-backend")
    | [
        .entity.name,
        ((.metrics[0]["backend.averageResponseTimeInSeconds"]        // 0) * 1000 | round | tostring),
        ((.metrics[0]["backend.averageTotalSessionTimeInSeconds"]    // 0) * 1000 | round | tostring),
        (.metrics[0]["backend.status"] // "?")
      ]
    | @tsv
  ' | awk -F'\t' '{printf "%-48s %8s %8s %8s\n", $1, $2, $3, $4}'

  # Server health
  echo "$JSON" | jq -r '
    .data[]
    | select(.entity.type == "ha-server")
    | [
        .entity.name,
        ((.metrics[0]["server.averageResponseTimeInSeconds"]      // 0) * 1000 | round | tostring),
        ((.metrics[0]["server.averageConnectTimeInSeconds"]       // 0) * 1000 | round | tostring),
        (.metrics[0]["server.status"] // "?")
      ]
    | @tsv
  ' | awk -F'\t' '{printf "%-48s %8s %8s %8s\n", $1, $2, $3, $4}'

  echo ""
  echo "(refreshing in ${INTERVAL}s...)"
  sleep "$INTERVAL"
done
