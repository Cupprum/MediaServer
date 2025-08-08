#!/usr/bin/env bash

###############################################################################
# Raspberry Pi Setup Script
#
# This script automates the initial setup of a Raspberry Pi including:
# - System updates
# - Installation of required tools
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

install_tools() {
    log_info "Installing required tools..."
    apt install -y jq tmux || { log_error "Failed to install required tools"; exit 1; }
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

configure_avahi() {
    log_info "Installing go-avahi-cname..."
    log_info "Finding latest go-avahi-cname release..."
    # [1:] used in jq removes the leading 'v' from the version string
    local latestAvahiVersion avahiUrl
    latestAvahiVersion=$(curl -sL https://api.github.com/repos/grishy/go-avahi-cname/releases/latest | jq -r ".tag_name | .[1:]")
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
    sudo mv go-avahi-cname /usr/local/bin/go-avahi-cname || { log_error "Failed to move go-avahi-cname to /usr/local/bin"; exit 1; }

    log_info "Setting permissions for go-avahi-cname..."
    chmod +x /usr/local/bin/go-avahi-cname

    log_info "Creating systemd service for go-avahi-cname..."
    sudo cp go-avahi-cname.service /etc/systemd/system/go-avahi-cname.service || { log_error "Failed to copy service file"; exit 1; }

    log_info "Enabling and starting go-avahi-cname service..."
    sudo systemctl daemon-reload
    sudo systemctl enable go-avahi-cname.service
    sudo systemctl start go-avahi-cname.service
    log_info "go-avahi-cname systemd service started successfully"

    log_info "go-avahi-cname installed successfully"
}

configure_usb_errors_hook() {
    log_info "Configuring USB Error hook"

    log_info "Moving 'check_usb_errors.sh' and 'restart_media_accessors' to /usr/local/bin..."
    sudo cp check_usb_errors.sh /usr/local/bin/check_usb_errors
    sudo cp restart_media_accessors.sh /usr/local/bin/restart_media_accessors

    log_info "Setting permissions for 'check_usb_errors' and 'restart_media_accessors'..."
    chmod +x /usr/local/bin/check_usb_errors
    chmod +x /usr/local/bin/restart_media_accessors

    log_info "Creating systemd service and timer for checking for usb errors..."
    sudo cp usb-error.service /etc/systemd/system/usb-error.service
    sudo cp usb-error.timer /etc/systemd/system/usb-error.timer

    log_info "Enabling and starting timer..."
    sudo systemctl daemon-reload
    sudo systemctl enable usb-error.timer
    sudo systemctl start usb-error.timer
    log_info "systemd timer for verifying usb errors started successfully"

    log_info "USB Error hook configured successfully"
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
        log_info "Finding latest PIA installer version..."

        local PIA_RELEASES_URL="https://api.github.com/repos/pia-foss/desktop/releases/latest"
        local PIA_INSTALLER_URL
        PIA_INSTALLER_URL=$(curl -s "${PIA_RELEASES_URL}" | \
            jq '.body | split("\r\n") | .[] | select(contains("linux_arm64")) | split(" - ")[1]')        
        local PIA_INSTALLER_PATH="/tmp/pia_installer.run"
        echo "Latest PIA installer URL: $PIA_INSTALLER_URL"

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
    install_tools
    install_docker
    configure_avahi
    configure_usb_errors_hook
    configure_k3s
    disable_wifi
    enable_vnc
    install_pia
    log_info "Rebooting system..."
    reboot
}

# Execute main function
main