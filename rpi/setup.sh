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

# Source common utilities
# shellcheck source=../utils.sh
source "$(dirname "${BASH_SOURCE[0]}")/../utils.sh"

###############################################################################
# Functions
###############################################################################

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
    # Check if Docker is already installed
    if command -v docker &>/dev/null; then
        log_info "Docker is already installed, skipping installation"
    else
        log_info "Installing Docker..."
        curl -sSL https://get.docker.com | sh || { log_error "Docker installation failed"; exit 1; }
        usermod -aG docker "$USER" || { log_error "Failed to add user to docker group"; exit 1; }
    fi
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
    if command -v piactl &>/dev/null; then
        log_info "PIA is already installed, skipping installation"
    else
        # Get latest PIA installer URL for arm64
        log_info "Finding latest PIA installer version..."
        local PIA_RELEASES_URL="https://api.github.com/repos/pia-foss/desktop/releases/latest"
        local PIA_INSTALLER_URL
        
        # Use curl and jq to parse GitHub API response and find the arm64 installer
        if ! command -v jq &>/dev/null; then
            log_info "Installing jq..."
            apt install -y jq || { log_error "Failed to install jq"; exit 1; }
        fi
        
        # Extract download URL for the arm64 installer from the latest release
        PIA_INSTALLER_URL=$(curl -s "$PIA_RELEASES_URL" | \
            jq -r '.assets[] | select(.name | test("pia-linux-arm64.*\\.run$")) | .browser_download_url')
        
        if [[ -z "$PIA_INSTALLER_URL" ]]; then
            log_error "Failed to find latest PIA installer for arm64"
            # Fallback to a known version
            PIA_INSTALLER_URL="https://installers.privateinternetaccess.com/download/pia-linux-arm64-3.6-08303.run"
            log_warn "Using fallback PIA installer version"
        else
            log_info "Found latest PIA installer: $(basename "$PIA_INSTALLER_URL")"
        fi
        local PIA_INSTALLER_PATH="/tmp/pia_installer.run"

        log_info "Downloading PIA installer..."
        curl -o "$PIA_INSTALLER_PATH" "$PIA_INSTALLER_URL" || { log_error "Failed to download PIA installer"; exit 1; }
        chmod +x "$PIA_INSTALLER_PATH" || { log_error "Failed to make PIA installer executable"; exit 1; }
        log_info "Running PIA installer..."
        "$PIA_INSTALLER_PATH" || { log_error "PIA installation failed"; exit 1; }
        rm "$PIA_INSTALLER_PATH" || log_warn "Failed to remove PIA installer"
        log_info "PIA installation completed"
        
        log_info "Enabling PIA background service..."
        piactl background enable || { log_error "Failed to enable PIA background service"; exit 1; }
        log_info "PIA background service enabled"
    fi
}

###############################################################################
# Main Script Execution
###############################################################################

main() {
    check_root
    log_info "Starting Raspberry Pi setup process..."
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