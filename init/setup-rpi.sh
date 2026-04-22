#!/bin/bash

set -e
set -u
set -o pipefail

###############################################################################
# Raspberry Pi Setup Script
# Purpose: Initial system setup (VNC, WiFi, Docker, Go, K3s, PIA, DNS).
###############################################################################

source "$(dirname "${BASH_SOURCE[0]}")/../utils.sh"

check_root() {
  if [[ "${EUID}" -ne 0 ]]; then
    log_error "This script must be run as root"
    exit 1
  fi
}

update_system() {
  log_info "Updating system packages..."
  apt update -y
  apt upgrade -y 
}

enable_vnc() {
  log_info "Enabling VNC..."
  raspi-config nonint do_vnc 0
}

disable_wifi() {
  log_info "Disabling WiFi..."
  if ! grep -q "dtoverlay=disable-wifi" /boot/firmware/config.txt; then
    echo "dtoverlay=disable-wifi" | tee -a /boot/firmware/config.txt
  fi
}

install_tools() {
  log_info "Installing required tools..."
    apt install -y jq tmux
}

install_docker() {
  if command -v docker &>/dev/null; then
    log_info "Docker is already installed, skipping."
  else
    log_info "Installing Docker..."
    curl -sSL https://get.docker.com | sh
        usermod -aG docker "$USER"
  fi
}

install_go() {
  if command -v go &>/dev/null; then
    log_info "Go is already installed, skipping."
  else
    log_info "Installing Go..."
    local file_name="go1.25.0.linux-arm64.tar.gz"
    wget "https://go.dev/dl/${file_name}"
    tar -C /usr/local -xzf "${file_name}"

    # shellcheck disable=SC2016
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    # shellcheck disable=SC1091
    source "$HOME/.bashrc"
    rm "$fileName"
    
    log_info "Go installed successfully."
  fi
}

configure_k3s() {
  log_info "Configuring k3s..."
  if ! grep -q "cgroup_enable=memory" /boot/firmware/cmdline.txt; then
    echo ' cgroup_enable=memory' | tee -a /boot/firmware/cmdline.txt
  fi
}

install_pia() {
  log_info "Finding latest PIA installer version..."
  local pia_releases_url="https://api.github.com/repos/pia-foss/desktop/releases/latest"
  local pia_installer_url
  pia_installer_url=$(curl -sSL "${pia_releases_url}" | jq -r '.body | split("\r\n") | .[] | select(contains("linux_arm64")) | split(" - ")[1]')
  local pia_installer_path="/tmp/pia_installer.run"
  
  log_info "Downloading PIA installer..."
  curl -sSL -o "${pia_installer_path}" "${pia_installer_url}"
  chmod +x "${pia_installer_path}"
  
  log_info "Running PIA installer..."
  su x42 -c "$pia_installer_path"
  rm -f "$pia_installer_path"
}

configure_pia() {
  log_info "Configuring PIA VPN..."
  piactl login .piactl_login
  piactl background enable
  piactl -u applysettings '{"killswitch":"on"}'
  piactl set protocol wireguard

  log_info "Reconnecting piactl..."
  piactl disconnect
  while [[ "$(piactl get connectionstate)" != "Disconnected" ]]; do
    sleep 2
  done
  piactl connect
  log_info "PIA configured successfully."
}

setup_pia() {
  if command -v piactl &>/dev/null; then
    log_info "PIA is already installed, skipping installation."
  else
    install_pia
    configure_pia
  fi
}

configure_avahi() {
  if command -v go-avahi-cname &>/dev/null; then
    log_info "go-avahi-cname is already installed, skipping."
  else
    log_info "Installing go-avahi-cname..."
    local latest_avahi_version
    # Note: [1:] used in jq removes the leading 'v' from the version string
    latest_avahi_version=$(curl -sSL https://api.github.com/repos/grishy/go-avahi-cname/releases/latest | jq -r ".tag_name | .[1:]")
    
    if [[ -z "${latest_avahi_version}" ]]; then
      log_error "Failed to retrieve the latest go-avahi-cname version"
      exit 1
    fi
    
    local avahi_url="https://github.com/grishy/go-avahi-cname/releases/download/v${latest_avahi_version}/go-avahi-cname_${latest_avahi_version}_linux_arm64.tar.gz"
    curl -sSL -o go-avahi-cname.tar.gz "${avahi_url}"
    tar -xzf go-avahi-cname.tar.gz
    
    sudo mv go-avahi-cname /usr/local/bin/go-avahi-cname
    sudo chmod +x /usr/local/bin/go-avahi-cname
    rm go-avahi-cname.tar.gz LICENSE README.md
    
    log_info "Creating systemd service for go-avahi-cname..."
    sudo cp go-avahi-cname.service /etc/systemd/system/go-avahi-cname.service

    sudo systemctl daemon-reload
    sudo systemctl enable go-avahi-cname.service
    sudo systemctl start go-avahi-cname.service
    log_info "go-avahi-cname installed and started."
  fi
}

configure_local_dns_resolution() {
  log_info "Configuring local DNS resolution..."
  if grep -q "pi.local" /etc/hosts; then
    log_info "pi.local entries already exist in /etc/hosts, skipping."
  else
    cat << EOF >> /etc/hosts
127.0.0.1       pi.local
127.0.0.1       qbittorrent.pi.local
127.0.0.1       flaresolverr.pi.local
127.0.0.1       prowlarr.pi.local
127.0.0.1       jellyfin.pi.local
127.0.0.1       argocd.pi.local
127.0.0.1       grafana.pi.local
EOF
    log_info "pi.local entries added to /etc/hosts."
  fi
}

main() {
  log_info "Starting Raspberry Pi setup process..."
  check_root
  update_system
  enable_vnc
  disable_wifi
  install_tools
  install_docker
  install_go
  configure_k3s
  setup_pia
  configure_avahi
  configure_local_dns_resolution
  
  log_info "Setup complete. Rebooting system..."
  reboot
}

main