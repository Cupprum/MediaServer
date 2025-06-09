#!/usr/bin/env bash

###############################################################################
# K3s Kubeconfig Generator Script
#
# This script:
# - Sets up local kubeconfig permissions
# - Generates a modified kubeconfig for remote access
# - Replaces localhost IP with the machine's IP
###############################################################################

# Source common utilities
# shellcheck source=../../utils.sh
source "$(dirname "${BASH_SOURCE[0]}")/../../utils.sh"

# Constants
readonly KUBECONFIG='/etc/rancher/k3s/k3s.yaml'

###############################################################################
# Functions
###############################################################################

setup_kubeconfig() {
    log_info "Setting up kubeconfig..."
    export KUBECONFIG
    sudo chmod +r "$KUBECONFIG" || { log_error "Failed to set kubeconfig permissions"; exit 1; }
}

print_separator() {
    printf -- "-%.0s" $(seq 1 50)
    echo
}

generate_remote_kubeconfig() {
    log_info "Generating kubeconfig for remote access..."
    if [ ! -f "$KUBECONFIG" ]; then
        log_error "Kubeconfig file not found at $KUBECONFIG"
        exit 1
    fi

    local new_kubeconfig
    new_kubeconfig=$(sed "s/127.0.0.1/192.168.0.100/g" "$KUBECONFIG") || { 
        log_error "Failed to generate remote kubeconfig"
        exit 1
    }
    
    print_separator
    echo "$new_kubeconfig"
}

###############################################################################
# Main Script Execution
###############################################################################

main() {
    setup_kubeconfig
    log_info "Generating kubeconfig to use with kubectl from other machine"
    generate_remote_kubeconfig
}

# Execute main function
main