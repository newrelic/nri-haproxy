#!/usr/bin/env bash
# start.sh — always-safe way to start or restart the demo stack.
# Run this every time, regardless of current state.
# It handles: first run, already running, config changes, source changes.
#
# Usage:
#   ./env-setup-local-container/start.sh                # nri-haproxy + OTel default (18 metrics)
#   ./env-setup-local-container/start.sh --otel-full    # nri-haproxy + all OTel optional metrics enabled

set -euo pipefail
cd "$(dirname "$0")"

# ── Parse arguments ───────────────────────────────────────────────────────────

OTEL_FULL=false
for arg in "$@"; do
  case "${arg}" in
    --otel-full) OTEL_FULL=true ;;
    *) echo "ERROR: Unknown argument: ${arg}"; echo "Usage: $0 [--otel-full]"; exit 1 ;;
  esac
done

# ── Set mode-specific env vars (picked up by docker-compose.yml) ──────────────

if [[ "${OTEL_FULL}" == "true" ]]; then
  export CLUSTER_NAME="demo-otel-full"
  export OTEL_CONFIG_FILE="./otel/config.full.yaml"
  OTEL_MODE="full (all 34 optional metrics enabled)"
else
  export CLUSTER_NAME="demo-otel-default"
  export OTEL_CONFIG_FILE="./otel/config.yaml"
  OTEL_MODE="default (18 metrics)"
fi

# ── Preflight checks ──────────────────────────────────────────────────────────

if [[ ! -f .env ]]; then
  echo "ERROR: .env file not found."
  echo "Run: cp .env.example .env  and fill in your NEW_RELIC_LICENSE_KEY"
  exit 1
fi

if ! grep -q "NEW_RELIC_LICENSE_KEY=." .env; then
  echo "ERROR: NEW_RELIC_LICENSE_KEY is not set in .env"
  exit 1
fi

if ! docker info > /dev/null 2>&1; then
  echo "ERROR: Docker is not running. Start Docker Desktop first."
  exit 1
fi

# ── Always build + start ──────────────────────────────────────────────────────
# --build   : rebuilds images if Dockerfile or source changed (no-op if nothing changed)
# --remove-orphans : cleans up containers for removed services
# -d        : detached mode

echo "Starting demo stack (OTel mode: ${OTEL_MODE})..."
docker compose up -d --build --remove-orphans

# ── Wait for HAProxy to be ready ──────────────────────────────────────────────

echo "Waiting for HAProxy stats page..."
for i in $(seq 1 15); do
  if curl -sf "http://localhost:8404/stats;csv;norefresh" > /dev/null 2>&1; then
    echo "HAProxy is up."
    break
  fi
  if [[ $i -eq 15 ]]; then
    echo "ERROR: HAProxy did not become ready in time."
    echo "Check logs: docker compose logs haproxy"
    exit 1
  fi
  sleep 2
done

# ── Summary ───────────────────────────────────────────────────────────────────

echo ""
echo "================================================================"
echo " Demo stack is running"
echo "================================================================"
echo ""
docker compose ps --format "table {{.Name}}\t{{.Status}}"
echo ""
echo " Mode               : OTel ${OTEL_MODE}"
echo " Cluster name       : ${CLUSTER_NAME}"
echo " OTel config        : ${OTEL_CONFIG_FILE}"
echo ""
echo " HAProxy stats page : http://localhost:8404/stats"
echo " HAProxy frontend   : http://localhost:8080"
echo ""
echo " Check infra agent  : docker compose logs -f newrelic-infra"
echo " Live metrics       : ./scripts/watch-metrics.sh"
echo " One snapshot       : ./scripts/run-nri.sh"
echo ""
echo " Stop stack         : docker compose down"
echo "================================================================"
