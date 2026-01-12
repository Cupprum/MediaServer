#!/usr/bin/env bash

set -e
set -u
set -o pipefail

###############################################################################
# ArgoCD Apps Management Script
#
# This script manages Helm charts for:
# - Jellyfin
# - qBittorrent
# - Prowlarr
# - Homepage
#
# Usage: ./setup.sh [install|delete]
###############################################################################

# Source common utilities
source "$(dirname "${BASH_SOURCE[0]}")/../utils.sh"

# Array of services
readonly SERVICES=(
    "jellyfin"
    "qbittorrent"
    "prowlarr"
    "homepage"
)

###############################################################################
# Functions
###############################################################################

show_help() {
    cat << EOF
Usage: $(basename "$0") [MODE]

Modes:
    install                 Install all ArgoCD applications
    configure               Configure all services after installation
    install_and_configure   Install applications and configure services
    delete                  Remove all ArgoCD applications
    cleanup                 Remove all configuration directories
    help, --help, -h        Display this help message

Applications managed:
    - Jellyfin
    - qBittorrent
    - Prowlarr
    - Homepage

Example:
    $(basename "$0") install_and_configure  # Install and configure all applications
    $(basename "$0") delete                 # Remove all applications
EOF
}

load_env_file() {
    log_info "Loading environment variables..."
    set -a
    source "$(dirname "${BASH_SOURCE[0]}")/configuration/.env"
    set +a
}

make_sure_folders_exist() {
    log_info "Ensuring necessary folders exist..."

    if [ -z "${MEDIASERVER_CONFIG_DIR:-}" ]; then
        log_error "MEDIASERVER_CONFIG_DIR is not set. Please set it in the .env file."
        exit 1
    fi

    local folders=(
        "$MEDIASERVER_CONFIG_DIR/homepage/config"
        "$MEDIASERVER_CONFIG_DIR/jellyfin/config"
        "$MEDIASERVER_CONFIG_DIR/jellyfin/cache"
        "$MEDIASERVER_CONFIG_DIR/prowlarr/config"
        "$MEDIASERVER_CONFIG_DIR/qbittorrent/config"
    )

    for folder in "${folders[@]}"; do
        if [ ! -d "$folder" ]; then
            log_info "Creating folder: $folder"
            mkdir -p "$folder" || {
                log_error "Failed to create folder: $folder"
                exit 1
            }
        fi
    done
}

make_sure_namespace_exists() {
    local namespace="server"

    if ! kubectl get namespace "$namespace" &> /dev/null; then
        log_info "Creating namespace: $namespace"
        kubectl create namespace "$namespace" || {
            log_error "Failed to create namespace: $namespace"
            exit 1
        }
    else
        log_info "Namespace $namespace already exists"
    fi
}

wait_for_jellyfin() {
    local url="$JELLYFIN_URL/health"
    local max_retries=30
    local wait_time=10
    local attempt=1

    # While `curl http://jellyfin.pi.local/health` does not return "Healthy"
    while [ $attempt -le $max_retries ]; do
        local status
        status=$(curl -s "$url" || echo "unreachable")
        if [ "$status" == "Healthy" ]; then
            log_info "Jellyfin is healthy"
            return 0
        else
            log_info "Waiting for Jellyfin to become healthy (Attempt: $attempt/$max_retries)..."
            sleep $wait_time
            ((attempt++))
        fi
    done
    log_error "Jellyfin did not become healthy within the expected time"
    exit 1
}

install() {
    log_info "Installing ArgoCD applications..."

    load_env_file

    make_sure_folders_exist
    make_sure_namespace_exists
    
    for service in "${SERVICES[@]}"; do
        log_info "Installing $service..."
        helm install "$service" ./argocd-config-chart --set "service=$service" || {
            log_error "Failed to install $service"
            exit 1
        }
    done

    log_info "All applications installed successfully"
}

configure() {
    log_info "Waiting for services to become ready..."

    load_env_file

    # From the services, jellyfin takes longest to start
    wait_for_jellyfin

    log_info "All services are up and running"

    # Change to configuration directory
    cd "$(dirname "${BASH_SOURCE[0]}")/configuration"; 
    
    log_info "Configuring services..."
    ./configure_services run || {
        cd .. # Return to previous directory on error
        log_error "Failed to configure services"
        exit 1
    }

    sleep 10

    log_info "Testing services..."
    ./configure_services test || {
        cd .. # Return to previous directory on error
        log_error "Service tests failed"
        exit 1
    }
    log_info "All services configured and tested successfully"
    
    # Return to previous directory
    cd ..
}

install_and_configure_apps() {
    install
    configure
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

cleanup() {
    log_info "Cleaning up config folders..."

    load_env_file

    sudo rm -rf "$MEDIASERVER_CONFIG_DIR/prowlarr/*"
    sudo rm -rf "$MEDIASERVER_CONFIG_DIR/jellyfin/*"
    sudo rm -rf "$MEDIASERVER_CONFIG_DIR/qbittorrent/*"
    sudo rm -rf "$MEDIASERVER_CONFIG_DIR/heimdall/*"
    sudo rm -rf "$MEDIASERVER_CONFIG_DIR/homepage/*"

    log_info "Cleanup completed"
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
            install
            ;;
        "configure")
            configure
            ;;
        "install_and_configure")
            install_and_configure_apps
            ;;
        "delete")
            delete_apps
            ;;
        "cleanup")
            cleanup
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