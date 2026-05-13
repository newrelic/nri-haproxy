#!/usr/bin/env bash
# provision.sh — one-time setup on a fresh EC2 instance.
# Supports Amazon Linux 2023 / Amazon Linux 2 / Ubuntu 22.04.
# Run once, then log out and back in for the docker group to take effect.
#
# Usage:
#   bash env-setup-local-container/ec2/provision.sh

set -euo pipefail

COMPOSE_VERSION="v2.27.0"

detect_os() {
  if [ -f /etc/os-release ]; then
    . /etc/os-release
    echo "${ID}"
  else
    echo "unknown"
  fi
}

install_amazon_linux() {
  echo "Detected Amazon Linux — installing Docker via dnf/yum..."
  sudo dnf install -y docker git 2>/dev/null || sudo yum install -y docker git
  sudo systemctl enable docker
  sudo systemctl start docker
  sudo usermod -aG docker "${USER}"
}

install_ubuntu() {
  echo "Detected Ubuntu — installing Docker via apt..."
  sudo apt-get update -y
  sudo apt-get install -y ca-certificates curl gnupg git
  sudo install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg \
    | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
  sudo chmod a+r /etc/apt/keyrings/docker.gpg
  echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
    https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" \
    | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
  sudo apt-get update -y
  sudo apt-get install -y docker-ce docker-ce-cli containerd.io
  sudo systemctl enable docker
  sudo systemctl start docker
  sudo usermod -aG docker "${USER}"
}

install_docker_compose() {
  echo "Installing Docker Compose ${COMPOSE_VERSION}..."
  PLUGIN_DIR="/usr/local/lib/docker/cli-plugins"
  sudo mkdir -p "${PLUGIN_DIR}"
  sudo curl -fsSL \
    "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-linux-$(uname -m)" \
    -o "${PLUGIN_DIR}/docker-compose"
  sudo chmod +x "${PLUGIN_DIR}/docker-compose"
}

# ── Main ──────────────────────────────────────────────────────────────────────

OS=$(detect_os)

case "${OS}" in
  amzn)
    install_amazon_linux
    ;;
  ubuntu)
    install_ubuntu
    ;;
  *)
    echo "ERROR: Unsupported OS '${OS}'. Install Docker manually."
    exit 1
    ;;
esac

install_docker_compose

echo ""
echo "================================================================"
echo " Provision complete"
echo "================================================================"
echo ""
echo " Docker version   : $(docker --version)"
echo " Compose version  : $(docker compose version)"
echo ""
echo " IMPORTANT: Log out and back in for the docker group to take effect."
echo "            (Otherwise you will need to prefix docker commands with sudo.)"
echo ""
echo " Next steps:"
echo "   exit"
echo "   ssh -i your-key.pem \$(whoami)@<instance-ip>"
echo "   cd nri-haproxy"
echo "   cp env-setup-local-container/.env.example env-setup-local-container/.env"
echo "   nano env-setup-local-container/.env   # set NEW_RELIC_LICENSE_KEY"
echo "   ./env-setup-local-container/start.sh"
echo "================================================================"
