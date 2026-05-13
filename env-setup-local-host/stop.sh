#!/usr/bin/env bash
# stop.sh — Tear down the local HAProxy + OTel stack.
# Kills all processes tracked in .pids/ and cleans up.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PID_DIR="$SCRIPT_DIR/.pids"

if [[ ! -d "$PID_DIR" ]]; then
  echo "No running stack found (no .pids directory)."
  exit 0
fi

echo "Stopping local HAProxy + OTel stack..."

# Kill each tracked process
for pidfile in "$PID_DIR"/*.pid; do
  [[ -f "$pidfile" ]] || continue
  name="$(basename "$pidfile" .pid)"
  pid="$(cat "$pidfile")"
  if kill -0 "$pid" 2>/dev/null; then
    kill "$pid" 2>/dev/null || true
    echo "  Stopped $name (pid $pid)"
  else
    echo "  $name (pid $pid) already stopped"
  fi
done

# Clean up runtime directories
rm -rf "$PID_DIR"

echo "All processes stopped. Logs preserved in .logs/ (delete manually if desired)."
