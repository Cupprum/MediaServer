#!/usr/bin/env bash

###############################################################################
# Restart deployments accessing media storage
###############################################################################

# Source common utilities
# shellcheck source=../../utils.sh
source "$(dirname "${BASH_SOURCE[0]}")/../../utils.sh"

main() {
    log_info "Restarting Jellyfin and qBittorrent deployments in the 'server' namespace..."
    kubectl rollout restart "deployment/jellyfin" -n server
    kubectl rollout restart "deployment/qbittorrent" -n server
    log_info "Deployments were restarted"
}

main