#!/usr/bin/env bash

###############################################################################
# Raspberry Pi Setup Script
#
# This script automates the initial setup of a Raspberry Pi including:
# - System updates
# - Docker installation
# - K3s preconfiguration
# - WiFi disable
# - VNC enablement
# - PIA VPN installation
###############################################################################

set -e  # Exit on fail
set -u  # Treat unset variables as an error
set -o pipefail  # Fail if any command in a pipeline fails

# Color definitions
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly NC='\033[0m' # No Color

###############################################################################
# Functions
###############################################################################

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root"
        exit 1
    fi
}

update_system() {
    log_info "Updating system packages..."
    apt update || { log_error "Failed to update package list"; exit 1; }
    apt upgrade -y || { log_error "Failed to upgrade packages"; exit 1; }
}

install_docker() {
    log_info "Installing Docker..."
    curl -sSL https://get.docker.com | sh || { log_error "Docker installation failed"; exit 1; }
    usermod -aG docker "$USER" || { log_error "Failed to add user to docker group"; exit 1; }
}

configure_k3s() {
    log_info "Configuring k3s..."
    if ! grep -q "cgroup_enable=memory" /boot/firmware/cmdline.txt; then
        echo ' cgroup_enable=memory' | tee -a /boot/firmware/cmdline.txt
    fi
}

disable_wifi() {
    log_info "Disabling WiFi..."
    if ! grep -q "dtoverlay=disable-wifi" /boot/firmware/config.txt; then
        echo "dtoverlay=disable-wifi" | tee -a /boot/firmware/config.txt
    fi
}

enable_vnc() {
    log_info "Enabling VNC..."
    raspi-config nonint do_vnc 0 || { log_error "Failed to enable VNC"; exit 1; }
    log_info "VNC has been enabled"
}

install_pia() {
    readonly PIA_URL="https://www.privateinternetaccess.com/installer/x/download_installer_linux_arm/arm64"
    readonly INSTALLER_PATH="/tmp/pia_installer.run"

    log_info "Downloading PIA installer..."
    curl -o "$INSTALLER_PATH" "$PIA_URL" || { log_error "Failed to download PIA installer"; exit 1; }
    chmod +x "$INSTALLER_PATH" || { log_error "Failed to make PIA installer executable"; exit 1; }
    log_info "Running PIA installer..."
    "$INSTALLER_PATH" || { log_error "PIA installation failed"; exit 1; }
    rm "$INSTALLER_PATH" || log_warn "Failed to remove PIA installer"
    log_info "PIA installation completed"
}

###############################################################################
# Main Script Execution
###############################################################################

main() {
    check_root
    log_info "Starting Raspberry Pi bootstrap process..."
    update_system
    install_docker
    configure_k3s
    disable_wifi
    enable_vnc
    install_pia
    log_info "Rebooting system..."
    reboot
}

# Execute main function
main