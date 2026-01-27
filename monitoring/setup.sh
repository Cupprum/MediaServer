#!/usr/bin/env bash

set -e
set -u
set -o pipefail

###############################################################################
# Monitoring Management Script
###############################################################################

# Source common utilities
source "$(dirname "${BASH_SOURCE[0]}")/../utils.sh"

###############################################################################
# Functions
###############################################################################

show_help() {
    cat << EOF
Usage: $(basename "$0") [MODE]

Modes:
    install     Install monitoring ArgoCD application
    delete      Remove monitoring ArgoCD application
    help        Display this help message

Example:
    $(basename "$0") install
    $(basename "$0") delete
EOF
}

load_env_file() {
    log_info "Loading environment variables..."
    set -a
    source "$(dirname "${BASH_SOURCE[0]}")/../.env"
    set +a

    if [ -z "${MEDIASERVER_GRAFANA_USERNAME:-}" ] || [ -z "${MEDIASERVER_GRAFANA_PASSWORD:-}" ]; then
        log_error "MEDIASERVER_GRAFANA_USERNAME and MEDIASERVER_GRAFANA_PASSWORD must be present in the .env file"
        exit 1
    fi

    if [ -z "${MEDIASERVER_CONFIG_DIR:-}" ]; then
        log_error "MEDIASERVER_CONFIG_DIR is not set in the .env file"
        exit 1
    fi
}

verify_namespace() {
    local namespace="monitoring"

    if ! kubectl get namespace "$namespace" &> /dev/null; then
        log_info "Creating namespace: $namespace..."
        kubectl create namespace "$namespace" || {
            log_error "Failed to create namespace: $namespace"
            exit 1
        }
    else
        log_info "Namespace $namespace already exists"
    fi
}

verify_folders() {
    log_info "Ensuring necessary folders exist..."
    
    local folder="$MEDIASERVER_CONFIG_DIR/config/prometheus"

    if [ ! -d "$folder" ]; then
        log_info "Creating folder: $folder"
        mkdir -p "$folder"
    fi
}

create_secrets() {
    log_info "Checking if grafana-admin-credentials secret already exists..."
    if kubectl get secret grafana-admin-credentials --namespace monitoring &>/dev/null; then
        log_info "grafana-admin-credentials secret already exists, skipping creation"
    else
        log_info "Creating grafana-admin-credentials secret..."
        kubectl create secret generic grafana-admin-credentials \
            --namespace monitoring \
            --from-literal=admin-user="${MEDIASERVER_GRAFANA_USERNAME}" \
            --from-literal=admin-password="${MEDIASERVER_GRAFANA_PASSWORD}" \
            --dry-run=client -o yaml | kubectl apply -f - || {
                log_error "Failed to create grafana-admin-credentials secret"
                exit 1
        }
        log_info "grafana-admin-credentials secret was created successfully"
    fi
}

install_monitoring() {
    log_info "Starting Monitoring installation..."

    load_env_file
    verify_namespace
    verify_folders
    create_secrets

    log_info "Installing ArgoCD application for monitoring..."
    kubectl apply --filename ./monitoring_app.yaml || {
        log_error "Failed to create ArgoCD application for monitoring"
        exit 1
    }
    log_info "ArgoCD application for monitoring was installed successfully"
}

delete_monitoring() {
    log_info "Removing ArgoCD application for monitoring..."
    kubectl delete --filename ./monitoring_app.yaml || {
        log_error "Failed to remove ArgoCD application for monitoring"
        exit 1
    }
        log_info "ArgoCD application for monitoring was removed successfully"

    log_info "Cleaning up grafana-admin-credentials secret..."
    kubectl delete secret grafana-admin-credentials --namespace monitoring || {
        log_error "grafana-admin-credentials secret was not deleted"
        exit 1
    }
    log_info "grafana-admin-credentials secret was deleted successfully"
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
            install_monitoring
            ;;
        "delete")
            delete_monitoring
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