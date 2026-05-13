#!/usr/bin/env bash
# Continuous traffic generator — sends requests to HAProxy every 200ms.
# Ctrl+C to stop.
set -e

HAPROXY_URL="${HAPROXY_URL:-http://127.0.0.1:8080}"

echo "Traffic generator started → $HAPROXY_URL (every 200ms)"
echo "Press Ctrl+C to stop."

while true; do
  curl -sf "$HAPROXY_URL/" > /dev/null 2>&1 || true
  curl -sf "$HAPROXY_URL/api/" > /dev/null 2>&1 || true
  curl -sf "$HAPROXY_URL/health" > /dev/null 2>&1 || true
  sleep 0.2
done
