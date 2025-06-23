#!/usr/bin/env bash

###############################################################################
# ArgoCD Monitoring App Management Script
#
# This script manages the ArgoCD application for the monitoring stack
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
    install     Install monitoring ArgoCD application
    delete      Remove monitoring ArgoCD application

Environment variables required:
    GRAFANA_USERNAME    Username for Grafana admin
    GRAFANA_PASSWORD    Password for Grafana admin

Example:
    $(basename "$0") install    # Install monitoring application
    $(basename "$0") delete     # Remove monitoring application
EOF
}

verify_kube_config() {
    if ! kubectl cluster-info &> /dev/null; then
        log_error "Failed to contact k8s cluster, is KUBECONFIG configured properly?"
        exit 1
    fi
}

install_monitoring() {
    log_info "Source environment variables..."
    set -a # automatically export all variables
    source .env
    set +a
    if [ -z "$GRAFANA_USERNAME" ] || [ -z "$GRAFANA_PASSWORD" ]; then
        log_error "GRAFANA_USERNAME and GRAFANA_PASSWORD must be present in the environment"
        exit 1
    fi
    log_info "Environment variables loaded successfully"

    log_info "Checking if monitoring namespace already exists..."
    if kubectl get namespace monitoring &>/dev/null; then
        log_info "Monitoring namespace already exists, skipping creation"
    else
        log_info "Creating monitoring namespace..."
        kubectl create namespace monitoring || {
            log_error "Failed to create monitoring namespace"
            exit 1
        }
        log_info "Monitoring namespace created successfully"
    fi

    log_info "Checking if grafana-credentials secret already exists..."
    if kubectl get secret grafana-credentials --namespace monitoring &>/dev/null; then
        log_info "grafana-credentials secret already exists, skipping creation"
    else
        log_info "Creating grafana-credentials secret..."
        kubectl create secret generic grafana-credentials \
            --namespace monitoring \
            --from-literal=admin-user="${GRAFANA_USERNAME}" \
            --from-literal=admin-password="${GRAFANA_PASSWORD}" \
            --dry-run=client -o yaml | kubectl apply -f - || {
                log_error "Failed to create grafana-credentials secret"
                exit 1
        }
        log_info "grafana-credentials secret created successfully"
    fi

    log_info "Installing ArgoCD application for monitoring..."
    kubectl apply --filename ./monitoring_app.yaml || {
        log_error "Failed to create ArgoCD application for monitoring"
        exit 1
    }
    log_info "ArgoCD application for monitoring was installed successfully"
}

delete_monitoring() {
    log_info "Removing ArgoCD application for monitoring..."
    kubectl delete \
        --filename ./monitoring_app.yaml \
        --namespace monitoring || {
            log_error "Failed to remove ArgoCD application for monitoring"
            exit 1
    }
    log_info "ArgoCD application for monitoring was removed successfully"

    log_info "Cleaning up secrets..."
    kubectl delete secret grafana-credentials --namespace monitoring || {
        log_error "grafana-credentials secret was not deleted"
        exit 1
    }
    log_info "grafana-credentials secret was deleted successfully"
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