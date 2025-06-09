#!/usr/bin/env bash

###############################################################################
# ArgoCD Apps Management Script
#
# This script manages Helm charts for:
# - Jellyfin
# - qBittorrent
# - Prowlarr
# - Heimdall
#
# Usage: ./script.sh [install|delete]
###############################################################################

set -e # Exit on fail
set -u # Treat unset variables as an error
set -o pipefail # Fail if any command in a pipeline fails

# Color definitions
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly NC='\033[0m' # No Color

# Array of services
readonly SERVICES=(
    "jellyfin"
    "qbittorrent"
    "prowlarr"
    "heimdall"
)

###############################################################################
# Functions
###############################################################################

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

show_help() {
    cat << EOF
Usage: $(basename "$0") [MODE]

Modes:
    install     Install all ArgoCD applications
    delete      Remove all ArgoCD applications

Applications managed:
    - Jellyfin
    - qBittorrent
    - Prowlarr
    - Heimdall

Example:
    $(basename "$0") install    # Install all applications
    $(basename "$0") delete    # Remove all applications
EOF
}

install_apps() {
    log_info "Installing ArgoCD applications..."
    
    for service in "${SERVICES[@]}"; do
        log_info "Installing $service..."
        helm install "$service" ./argocd-config-chart --set "service=$service" || {
            log_error "Failed to install $service"
            exit 1
        }
    done
    
    log_info "All applications installed successfully"
}

delete_apps() {
    log_info "Removing ArgoCD applications..."
    
    for service in "${SERVICES[@]}"; do
        log_info "Removing $service..."
        helm uninstall "$service" || {
            log_error "Failed to remove $service"
            exit 1
        }
    done
    
    log_info "All applications removed successfully"
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

    case "$1" in
        "install")
            install_apps
            ;;
        "delete")
            delete_apps
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