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

create_secrets() {
    log_info "Checking grafana-admin-credentials secret..."
    if kubectl get secret grafana-admin-credentials --namespace monitoring &>/dev/null; then
        log_info "grafana-admin-credentials secret already exists, skipping."
    else
        log_info "Creating grafana-admin-credentials secret..."
        kubectl create secret generic grafana-admin-credentials \
            --namespace monitoring \
            --from-literal=admin-user="${MEDIASERVER_GRAFANA_USERNAME}" \
            --from-literal=admin-password="${MEDIASERVER_GRAFANA_PASSWORD}" \
            --dry-run=client -o yaml | kubectl apply -f -
        log_info "grafana-admin-credentials secret created."
    fi
}

install_monitoring() {
    log_info "Starting Monitoring installation..."

    load_env_vars "$(dirname "${BASH_SOURCE[0]}")/../.env"
    require_env_vars "MEDIASERVER_GRAFANA_USERNAME" "MEDIASERVER_GRAFANA_PASSWORD" "MEDIASERVER_CONFIG_DIR"

    ensure_namespace "monitoring"
    ensure_directory "$MEDIASERVER_CONFIG_DIR/prometheus"
    create_secrets

    log_info "Installing ArgoCD application for monitoring..."
    kubectl apply --filename ./monitoring_app.yaml
    log_info "ArgoCD application for monitoring installed successfully."
}

delete_monitoring() {
    log_info "Removing ArgoCD application for monitoring..."
    kubectl delete --filename ./monitoring_app.yaml || log_warn "App not found or already deleted."

    log_info "Cleaning up grafana-admin-credentials secret..."
    kubectl delete secret grafana-admin-credentials --namespace monitoring || log_warn "Secret not found or already deleted."
    
    log_info "Monitoring deletion completed."
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

main "$@"