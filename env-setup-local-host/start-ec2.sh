#!/usr/bin/env bash
# start-ec2.sh — Launch the full HAProxy + OTel stack on an EC2 Linux host (no Docker).
#
# Supports: Amazon Linux 2023 (dnf), Ubuntu 22.04+ (apt)
#
# Prerequisites (installed automatically):
#   - haproxy (via dnf or apt)
#   - python3 (via dnf or apt)
#   - otelcol-contrib (downloaded from GitHub releases)
#   - curl, tar, gzip (standard on all AMIs)
#
# What it starts (ALL on host, no containers):
#   1. Five Python HTTP backend servers (web1-3 on 9001-9003, api1-2 on 9011-9012)
#   2. HAProxy on ports 8080 (frontend) and 8404 (stats)
#   3. OTel Collector with haproxyreceiver (scrapes stats, ships to NR staging)
#   4. A traffic generator sending requests every 200ms
#
# All processes are launched in the background. Use ./stop.sh to tear down.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PID_DIR="$SCRIPT_DIR/.pids"
LOG_DIR="$SCRIPT_DIR/.logs"
BIN_DIR="$SCRIPT_DIR/.bin"
ENV_FILE="$SCRIPT_DIR/.env"

OTEL_VERSION="0.121.0"

# ── Load .env ─────────────────────────────────────────────────────────────────
if [[ ! -f "$ENV_FILE" ]]; then
  echo "ERROR: $ENV_FILE not found." >&2
  echo "  cp env-setup-local-host/.env.example env-setup-local-host/.env" >&2
  echo "  Then set NEW_RELIC_LICENSE_KEY in that file." >&2
  exit 1
fi
set -a
source "$ENV_FILE"
set +a

if [[ -z "${NEW_RELIC_LICENSE_KEY:-}" || "$NEW_RELIC_LICENSE_KEY" == "YOUR_KEY_HERE" ]]; then
  echo "ERROR: NEW_RELIC_LICENSE_KEY is not set in $ENV_FILE" >&2
  exit 1
fi

# ── Cleanup any previous run ─────────────────────────────────────────────────
if [[ -d "$PID_DIR" ]]; then
  echo "Stopping previous run..."
  "$SCRIPT_DIR/stop.sh" 2>/dev/null || true
fi
mkdir -p "$PID_DIR" "$LOG_DIR" "$BIN_DIR"

# ── Install packages via dnf or apt ──────────────────────────────────────────
install_pkg() {
  local pkg="$1"
  if command -v "$pkg" &>/dev/null; then return 0; fi
  echo "Installing $pkg..."
  if command -v dnf &>/dev/null; then
    sudo dnf install -y "$pkg"
  elif command -v apt-get &>/dev/null; then
    sudo apt-get update -qq && sudo apt-get install -y "$pkg"
  else
    echo "ERROR: No supported package manager (dnf or apt) found." >&2
    exit 1
  fi
}

install_pkg haproxy
install_pkg python3
install_pkg curl

# ── Download otelcol-contrib binary if not already present ───────────────────
OTEL_BIN="$BIN_DIR/otelcol-contrib"
if [[ ! -x "$OTEL_BIN" ]]; then
  ARCH="$(uname -m)"
  case "$ARCH" in
    aarch64|arm64) OTEL_ARCH="arm64" ;;
    x86_64)        OTEL_ARCH="amd64" ;;
    *)             echo "ERROR: Unsupported architecture: $ARCH" >&2; exit 1 ;;
  esac
  OTEL_TAR="otelcol-contrib_${OTEL_VERSION}_linux_${OTEL_ARCH}.tar.gz"
  OTEL_URL="https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v${OTEL_VERSION}/${OTEL_TAR}"

  echo "Downloading otelcol-contrib v${OTEL_VERSION} (linux/${OTEL_ARCH})..."
  curl -fSL "$OTEL_URL" -o "$BIN_DIR/$OTEL_TAR"
  tar -xzf "$BIN_DIR/$OTEL_TAR" -C "$BIN_DIR"
  rm -f "$BIN_DIR/$OTEL_TAR"
  chmod +x "$OTEL_BIN"
  echo "  Installed: $OTEL_BIN"
fi

# ── Start backend servers ────────────────────────────────────────────────────
echo "Starting backend servers..."

start_backend() {
  local port="$1" name="$2"
  python3 "$SCRIPT_DIR/backends/server.py" "$port" "$name" &
  echo $! > "$PID_DIR/$name.pid"
}

start_backend 9001 web1
start_backend 9002 web2
start_backend 9003 web3
start_backend 9011 api1
start_backend 9012 api2

# Give backends a moment to bind
sleep 0.5

# ── Start HAProxy ────────────────────────────────────────────────────────────
echo "Starting HAProxy..."
haproxy -f "$SCRIPT_DIR/haproxy/haproxy.cfg" -D -p "$PID_DIR/haproxy.pid"
echo "  HAProxy frontend: http://127.0.0.1:8080"
echo "  HAProxy stats:    http://127.0.0.1:8404/stats"

# ── Start OTel Collector ─────────────────────────────────────────────────────
echo "Starting OTel Collector (haproxyreceiver → NR staging)..."
"$OTEL_BIN" --config "$SCRIPT_DIR/otel/config.yaml" \
  > "$LOG_DIR/otel-collector.log" 2>&1 &
echo $! > "$PID_DIR/otel-collector.pid"
echo "  OTel logs: $LOG_DIR/otel-collector.log"

# ── Start traffic generator ──────────────────────────────────────────────────
echo "Starting traffic generator..."
bash "$SCRIPT_DIR/scripts/traffic-gen.sh" > /dev/null 2>&1 &
echo $! > "$PID_DIR/traffic-gen.pid"

# ── Summary ──────────────────────────────────────────────────────────────────
# Try to get public IP (works on EC2, fails gracefully elsewhere)
PUBLIC_IP="$(curl -sf --connect-timeout 2 http://169.254.169.254/latest/meta-data/public-ipv4 2>/dev/null || echo '<EC2-PUBLIC-IP>')"

echo ""
echo "═══════════════════════════════════════════════════════════════════"
echo " HAProxy + OTel stack is running on EC2! (all on host, no Docker)"
echo "═══════════════════════════════════════════════════════════════════"
echo ""
echo " HAProxy frontend:   http://127.0.0.1:8080"
echo " HAProxy stats page: http://127.0.0.1:8404/stats"
echo " HAProxy stats CSV:  http://127.0.0.1:8404/stats;csv"
echo ""
echo " External access (if security group allows):"
echo "   Stats: http://${PUBLIC_IP}:8404/stats"
echo ""
echo " Backends (host):    web1(:9001) web2(:9002) web3(:9003)"
echo "                     api1(:9011) api2(:9012)"
echo ""
echo " OTel Collector:     Scraping stats every 15s → NR STAGING"
echo " OTel logs:          tail -f $LOG_DIR/otel-collector.log"
echo ""
echo " Data destination:   https://staging-otlp.nr-data.net (STAGING)"
echo " License key:        loaded from env-setup-local-host/.env"
echo ""
echo " Run nri-haproxy against the same HAProxy:"
echo "   go run ./src/... -stats_url http://127.0.0.1:8404/stats -cluster_name ec2-test"
echo ""
echo " Stop everything:    ./env-setup-local-host/stop.sh"
echo ""
