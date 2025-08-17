#!/usr/bin/env bash

set -e
set -u
set -o pipefail

###############################################################################
# Bootstrap Monitoring
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
    log_info "Starting Monitoring bootstrap..."

    # Create namespaces
    log_info "Creating namespace..."
    kubectl create namespace monitoring

    log_info "Creating required directories..."
    mkdir -p /home/x42/MediaServer/config
    mkdir -p /home/x42/MediaServer/config/prometheus

    log_info "Monitoring bootstrap completed successfully."
}

# Execute main function
main "$@"