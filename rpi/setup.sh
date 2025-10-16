#!/usr/bin/env bash

set -e
set -u
set -o pipefail

###############################################################################
# Raspberry Pi Setup Script
#
# This script automates the initial setup of a Raspberry Pi including:
# - System updates
# - VNC enablement
# - WiFi disable
# - Installation of required tools
# - Docker installation
# - Go installation
# - K3s preconfiguration
# - PIA VPN installation and configuration
# - mDNS setup with go-avahi-cname
# - Local DNS resolution configuration
###############################################################################

# Source common utilities
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
    apt update -y
    apt upgrade -y 
}

enable_vnc() {
    log_info "Enabling VNC..."
    raspi-config nonint do_vnc 0
    log_info "VNC has been enabled"
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
        log_info "Docker is already installed, skipping installation"
    else
        log_info "Installing Docker..."
        curl -sSL https://get.docker.com | sh
        usermod -aG docker "$USER"
    fi
}

install_go() {
    if command -v go &>/dev/null; then
        log_info "Go is already installed, skipping installation"
    else
        log_info "Installing go..."
        fileName="go1.25.0.linux-arm64.tar.gz"
        wget "https://go.dev/dl/$fileName"
        sudo tar -C /usr/local -xzf "$fileName"

        # shellcheck disable=SC2016
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        # shellcheck disable=SC1091
        source "$HOME/.bashrc"
        rm "$fileName"
        
        log_info "Added Go to PATH in ~/.bashrc"
    fi
}

configure_k3s() {
    log_info "Configuring k3s..."
    if ! grep -q "cgroup_enable=memory" /boot/firmware/cmdline.txt; then
        echo ' cgroup_enable=memory' | tee -a /boot/firmware/cmdline.txt
    fi
}

install_pia() {
    log_info 'Installing PIA...'
    log_info "Finding latest PIA installer version..."

    local PIA_RELEASES_URL="https://api.github.com/repos/pia-foss/desktop/releases/latest"
    local PIA_INSTALLER_URL
    PIA_INSTALLER_URL=$(curl -sSL "${PIA_RELEASES_URL}" | \
        jq -r '.body | split("\r\n") | .[] | select(contains("linux_arm64")) | split(" - ")[1]')
    local PIA_INSTALLER_PATH="/tmp/pia_installer.run"
    log_info "Latest PIA installer URL: $PIA_INSTALLER_URL"

    log_info "Downloading PIA installer..."
    rm "$PIA_INSTALLER_PATH"
    curl -sSL -o "$PIA_INSTALLER_PATH" "$PIA_INSTALLER_URL"
    chmod +x "$PIA_INSTALLER_PATH"
    log_info "Running PIA installer..."
    su x42 -c "$PIA_INSTALLER_PATH"
    rm "$PIA_INSTALLER_PATH"
    log_info "PIA installation completed"
}

configure_pia() {
    log_info 'Configuring PIA...'
    piactl login .piactl_login
    piactl background enable
    piactl -u applysettings '{"killswitch":"on"}'
    piactl set protocol wireguard

    log_info 'Reconnecting piactl...'
    piactl disconnect
    while [ "$(piactl get connectionstate)" != "Disconnected" ]; do
        log_info "Waiting for VPN to disconnect..."
        sleep 2
    done
    piactl connect
    log_info 'PIA configured'
}

setup_pia() {
    if command -v piactl &>/dev/null; then
        log_info "PIA is already installed, skipping installation"
    else
        install_pia
        configure_pia
    fi
}

configure_avahi() {
    log_info "Checking if go-avahi-cname is already installed"

    if command -v go-avahi-cname &>/dev/null; then
        log_info "go-avahi-cname is already installed, skipping installation"
    else
        log_info "Installing go-avahi-cname..."
        log_info "Finding latest go-avahi-cname release..."
        # Note: [1:] used in jq removes the leading 'v' from the version string
        local latestAvahiVersion avahiUrl
        latestAvahiVersion=$(curl -sSL https://api.github.com/repos/grishy/go-avahi-cname/releases/latest | jq -r ".tag_name | .[1:]")
        if [[ -z "$latestAvahiVersion" ]]; then
            log_error "Failed to retrieve the latest go-avahi-cname version"
            exit 1
        fi
        log_info "Latest go-avahi-cname release found: ${latestAvahiVersion}"
        
        log_info "Download the latest release of go-avahi-cname..."
        avahiUrl="https://github.com/grishy/go-avahi-cname/releases/download/v${latestAvahiVersion}/go-avahi-cname_${latestAvahiVersion}_linux_arm64.tar.gz"
        curl -sSL -o go-avahi-cname.tar.gz "${avahiUrl}" || { log_error "Failed to download go-avahi-cname"; exit 1; }
        tar -xzf go-avahi-cname.tar.gz
        # Clean up the tarball and other files
        rm go-avahi-cname.tar.gz LICENSE README.md

        log_info "Moving go-avahi-cname to /usr/local/bin..."
        sudo mv go-avahi-cname /usr/local/bin/go-avahi-cname

        log_info "Setting permissions for go-avahi-cname..."
        sudo chmod +x /usr/local/bin/go-avahi-cname
    fi

    log_info "Creating systemd service for go-avahi-cname..."
    sudo cp go-avahi-cname.service /etc/systemd/system/go-avahi-cname.service

    log_info "Enabling and starting go-avahi-cname service..."
    sudo systemctl daemon-reload
    sudo systemctl enable go-avahi-cname.service
    sudo systemctl start go-avahi-cname.service
    log_info "go-avahi-cname systemd service started successfully"

    log_info "go-avahi-cname installed successfully"
}

configure_local_dns_resolution() {
    log_info "Configuring local DNS resolution..."
    
    # Check if /etc/hosts already contains pi.local entries
    if grep -q "pi.local" /etc/hosts; then
        log_info "pi.local entries already exist in /etc/hosts, skipping"
    else
        log_info "Adding pi.local entries to /etc/hosts"
        cat << EOF >> /etc/hosts
127.0.0.1       pi.local
127.0.0.1       jellyfin.pi.local
127.0.0.1       qbittorrent.pi.local
127.0.0.1       prowlarr.pi.local
127.0.0.1       argocd.pi.local
127.0.0.1       grafana.pi.local
EOF
        log_info "pi.local entries added to /etc/hosts successfully"
    fi
}

###############################################################################
# Main Script Execution
###############################################################################

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
    log_info "Rebooting system..."
    reboot
}

# Execute main function
main