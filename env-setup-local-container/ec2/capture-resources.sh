#!/usr/bin/env bash
# capture-resources.sh — snapshot CPU and RAM for the two collector containers.
# Run once at steady state (~15 min after start) and again after the 12-hour run.
#
# Usage (from repo root):
#   bash env-setup-local-container/ec2/capture-resources.sh

set -euo pipefail
cd "$(dirname "$0")/../.."   # repo root

COMPOSE_FILE="env-setup-local-container/docker-compose.yml"

# Container names as defined in docker-compose.yml
NRI_CONTAINER="demo-newrelic-infra-1"
OTEL_CONTAINER="demo-otel-collector-1"

echo ""
echo "================================================================"
echo " Resource snapshot  — $(date -u '+%Y-%m-%d %H:%M:%S UTC')"
echo "================================================================"

echo ""
echo "All containers (docker stats, single sample):"
echo ""
docker stats --no-stream --format \
  "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}" \
  $(docker compose -f "${COMPOSE_FILE}" ps -q)

echo ""
echo "================================================================"
echo " Collector summary"
echo "================================================================"
echo ""

for CONTAINER in "${NRI_CONTAINER}" "${OTEL_CONTAINER}"; do
  if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER}$"; then
    STATS=$(docker stats --no-stream --format \
      "{{.CPUPerc}} {{.MemUsage}} {{.MemPerc}}" "${CONTAINER}")
    CPU=$(echo "${STATS}" | awk '{print $1}')
    MEM=$(echo "${STATS}" | awk '{print $2, $3, $4}')
    MEMP=$(echo "${STATS}" | awk '{print $5}')
    echo "  ${CONTAINER}"
    echo "    CPU : ${CPU}"
    echo "    RAM : ${MEM} (${MEMP})"
    echo ""
  else
    echo "  WARNING: ${CONTAINER} is not running. Check: docker compose -f ${COMPOSE_FILE} ps"
    echo ""
  fi
done

echo "================================================================"
echo " Paste the above into the Confluence page's Resource Comparison"
echo " section, or run: bash env-setup-local-container/ec2/capture-resources.sh >> resources.log"
echo "================================================================"
