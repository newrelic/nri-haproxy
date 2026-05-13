#!/usr/bin/env bash
# start.sh — Auto-detects OS and runs the appropriate start script.
#
# macOS  → start-mac.sh  (uses Homebrew, downloads darwin binary)
# Linux  → start-ec2.sh  (uses dnf/apt, downloads linux binary)
#
# You can also call the platform-specific script directly if you prefer.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

case "$(uname -s)" in
  Darwin)
    echo "Detected macOS — running start-mac.sh"
    exec "$SCRIPT_DIR/start-mac.sh"
    ;;
  Linux)
    echo "Detected Linux — running start-ec2.sh"
    exec "$SCRIPT_DIR/start-ec2.sh"
    ;;
  *)
    echo "ERROR: Unsupported OS: $(uname -s). Only macOS and Linux are supported." >&2
    exit 1
    ;;
esac
