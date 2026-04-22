#!/usr/bin/env bash

set -e
set -u
set -o pipefail

###############################################################################
# Raspberry Pi Media Server Installation Script
#
# This script configures MediaServer on Raspberry Pi, specifically:
# - Configure the Raspberry Pi
# - Setup and configure Kubernetes
# - Deploy MediaServer to Kubernetes
#
# Usage: ./setup.sh
###############################################################################

# Source common utilities
source "$(dirname "${BASH_SOURCE[0]}")/utils.sh"

readonly HOSTNAME='x42@pi.local'
readonly REPO_URL='git@github.com:Cupprum/MediaServer.git'

main() {
    log_info "Checking if MediaServer folder exists on ${HOSTNAME}..."
    
    if ssh "$HOSTNAME" '! test -d MediaServer'; then
        log_info "MediaServer folder not found. Cloning repository..."

        ssh "$HOSTNAME" "git clone ${REPO_URL}"
        if [ $? -eq 128 ]; then
            log_warn "Could not clone repository. Adding GitHub to known hosts..."
            ssh "$HOSTNAME" 'ssh-keyscan github.com >> ~/.ssh/known_hosts 2>/dev/null'
            
            log_info "Retrying repository clone..."
            ssh "$HOSTNAME" "git clone ${REPO_URL}"
        fi
        log_info "Repository cloned successfully."

        log_info "Copying secrets to MediaServer..."
        scp .env "$HOSTNAME:~/MediaServer/.env"
        scp .piactl_login "$HOSTNAME:~/MediaServer/.piactl_login"
        log_info "Secrets copied."
    else
        log_info "MediaServer folder already exists."
    fi

    log_info "Configuring Raspberry Pi..."
    ssh "$HOSTNAME" 'cd MediaServer/rpi && sudo ./setup.sh'

    log_info "Restarting the Raspberry Pi..."
    ssh "$HOSTNAME" 'sudo reboot'

    # Note: the next ssh command would be waiting there either way
    log_info "Sleeping for 120 seconds to allow Raspberry Pi to reboot..."
    sleep 120

    log_info "Setting up and configuring Kubernetes..."
    ssh "$HOSTNAME" 'cd MediaServer/k8s && ./setup.sh'

    log_info "Deploying MediaServer to Kubernetes..."
    ssh "$HOSTNAME" 'cd MediaServer/server && ./setup.sh install'

    log_info "Deploying Monitoring to Kubernetes..."
    ssh "$HOSTNAME" 'cd MediaServer/monitoring && ./setup.sh install'
    
    log_info "MediaServer installation completed successfully!"
}

# Execute main function
main