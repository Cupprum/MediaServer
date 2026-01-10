#!/usr/bin/env bash

set -e
set -u
set -o pipefail

###############################################################################
# Bootstrap MediaServer
#
# This script:
# - creates namespace
# - creates required folders
#
# Usage: ./bootstrap.sh
###############################################################################

# Source common utilities
source "$(dirname "${BASH_SOURCE[0]}")/../utils.sh"

main() {
    log_info "Starting MediaServer bootstrap..."

    # Create namespaces
    log_info "Creating namespace..."
    kubectl create namespace server

    log_info "Creating required directories..."
    mkdir -p /home/x42/MediaServer/config
    mkdir -p /home/x42/MediaServer/config/homepage/config
    mkdir -p /home/x42/MediaServer/config/jellyfin/config
    mkdir -p /home/x42/MediaServer/config/jellyfin/cache
    mkdir -p /home/x42/MediaServer/config/prowlarr/config
    mkdir -p /home/x42/MediaServer/config/qbittorrent/config

    log_info "MediaServer bootstrap completed successfully."
}

# Execute main function
main "$@"