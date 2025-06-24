#!/usr/bin/env bash

###############################################################################
# Ingress Management Script
#
# This script configures Ingress to services:
#
# Usage: ./script.sh [install|delete]
###############################################################################

# Source common utilities
# shellcheck source=../utils.sh
source "$(dirname "${BASH_SOURCE[0]}")/../utils.sh"

###############################################################################
# Functions
###############################################################################

show_help() {
    cat << EOF
Usage: $(basename "$0") [MODE]

Modes:
    install     Create Ingress
    delete      Delete Ingress

Example:
    $(basename "$0") install    # Create Ingress
    $(basename "$0") delete    # Delete Ingress
EOF
}

verify_kube_config() {
    if ! kubectl cluster-info &> /dev/null; then
        log_error "Failed to contact k8s cluster, is KUBECONFIG configured properly?"
        exit 1
    fi
}

install_ingress() {
    log_info "Installing Ingress..."
    
    kubeclt apply --filename "./ingress.yaml"
    
    log_info "Ingress installed successfully"
}

delete_ingress() {
    log_info "Removing Ingress..."

    kubeclt delete --filename "./ingress.yaml"
    
    log_info "Ingress removed successfully"
}

###############################################################################
# Main Script Execution
###############################################################################

main() {
    if [ $# -ne 1 ]; then
        log_error "Invalid number of arguments"
        show_help
        exit 1
    fi

    verify_kube_config

    case "$1" in
        "install")
            install_ingress
            ;;
        "delete")
            delete_ingress
            ;;
        "help"|"--help"|"-h")
            show_help
            exit 0
            ;;
        *)
            log_error "Invalid mode: $1"
            show_help
            exit 1
            ;;
    esac
}

# Execute main function
main "$@"