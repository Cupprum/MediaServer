#!/usr/bin/env bash

set -e
set -u
set -o pipefail

###############################################################################
# ArgoCD Apps Management Script
###############################################################################

# Source common utilities
source "$(dirname "${BASH_SOURCE[0]}")/../utils.sh"

readonly SERVICES=(
    "jellyfin"
    "qbittorrent"
    "prowlarr"
    "flaresolverr"
)

###############################################################################
# Functions
###############################################################################

show_help() {
    cat << EOF
Usage: $(basename "$0") [MODE]

Modes:
    install_services        Install all ArgoCD applications
    configure_services      Configure all services after installation
    install                 Install applications and configure services
    delete_services         Remove all ArgoCD applications
    cleanup_services        Remove all configuration directories
    delete                  Remove all ArgoCD application and all configs
    help, --help, -h        Display this help message

Applications managed:
    - Jellyfin
    - qBittorrent
    - Prowlarr
    - Flaresolverr

Example:
    $(basename "$0") install  # Install and configure all applications
    $(basename "$0") delete   # Remove all applications and their configs
EOF
}

get_service_dirs() {
    require_env_vars "MEDIASERVER_CONFIG_DIR"

    echo "$MEDIASERVER_CONFIG_DIR/qbittorrent/config"
    echo "$MEDIASERVER_CONFIG_DIR/flaresolverr/config"
    echo "$MEDIASERVER_CONFIG_DIR/prowlarr/config"
    echo "$MEDIASERVER_CONFIG_DIR/jellyfin/config"
    echo "$MEDIASERVER_CONFIG_DIR/jellyfin/cache"
}

setup_directories() {
    log_info "Ensuring necessary folders exist..."
    local dirs
    mapfile -t dirs < <(get_service_dirs)
    
    for dir in "${dirs[@]}"; do
        ensure_directory "$dir"
    done
}

wait_for_jellyfin() {
    require_env_vars "MEDIASERVER_JELLYFIN_URL"
    
    local url="${MEDIASERVER_JELLYFIN_URL}/health"
    local max_retries=30
    local wait_time=10
    local attempt=1

    while [ $attempt -le $max_retries ]; do
        local status
        status=$(curl -s "$url" || echo "unreachable")
        if [ "$status" == "Healthy" ]; then
            log_info "Jellyfin is healthy."
            return 0
        else
            log_info "Waiting for Jellyfin to become healthy (Attempt: $attempt/$max_retries)..."
            sleep $wait_time
            ((attempt++))
        fi
    done
    
    log_error "Jellyfin did not become healthy within the expected time."
    exit 1
}

install_services() {
    log_info "Installing ArgoCD applications..."

    load_env_vars "$(dirname "${BASH_SOURCE[0]}")/../.env"
    setup_directories
    ensure_namespace "server"
    
    for service in "${SERVICES[@]}"; do
        log_info "Installing $service..."
        helm install "$service" ./argocd-config-chart --set "service=$service"
    done

    log_info "All applications installed successfully."
}

configure_services() {
    log_info "Waiting for services to become ready..."

    load_env_vars "$(dirname "${BASH_SOURCE[0]}")/../.env"
    wait_for_jellyfin

    log_info "All services are up and running."

    cd "$(dirname "${BASH_SOURCE[0]}")/configuration"
    
    log_info "Configuring services..."
    ./configure_services run || {
        log_error "Failed to configure services."
        exit 1
    }

    sleep 10

    log_info "Testing services..."
    ./configure_services test || {
        log_error "Service tests failed."
        exit 1
    }
    
    log_info "All services configured and tested successfully."
    cd ..
}

delete_services() {
    log_info "Removing ArgoCD applications..."
    for service in "${SERVICES[@]}"; do
        log_info "Removing $service..."
        helm uninstall "$service" || log_warn "Failed to remove $service or it doesn't exist."
    done
    log_info "All applications removed successfully."
}

cleanup_services() {
    log_info "Cleaning up config folders..."
    
    load_env_vars "$(dirname "${BASH_SOURCE[0]}")/../.env"
    
    local dirs
    mapfile -t dirs < <(get_service_dirs)

    for dir in "${dirs[@]}"; do
        if [ -d "$dir" ]; then
            log_info "Removing: $dir"
            sudo rm -rf "$dir" || ( sleep 5 && sudo rm -rf "$dir" )
        fi
    done

    log_info "Cleanup completed."
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
        "install_services")
            install_services
            ;;
        "configure_services")
            configure_services
            ;;
        "install")
            install_services
            configure_services
            ;;
        "delete_services")
            delete_services
            ;;
        "cleanup_services")
            cleanup_services
            ;;
        "delete")
            delete_services
            cleanup_services
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

main "$@"