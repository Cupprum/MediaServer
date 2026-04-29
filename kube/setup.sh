#!/bin/bash

set -e
set -u
set -o pipefail

###############################################################################
# ArgoCD Apps Management Script
# Purpose: Manage application installation, configuration, and cleanup.
###############################################################################

source "$(dirname "${BASH_SOURCE[0]}")/../utils.sh"

show_help() {
  cat << EOF
Usage: $(basename "$0") [MODE]

Modes:
    install                 Install applications and configure services
    delete                  Remove all ArgoCD application and all configs
    cleanup-images          Remove all images in k3s
    help, --help, -h        Display this help message
EOF
}

configs() {
  kubectl delete secret grafana-admin-credentials -n monitoring --ignore-not-found
  kubectl create secret generic grafana-admin-credentials \
    --namespace monitoring \
    --from-literal=admin-user="${MEDIASERVER_GRAFANA_USERNAME}" \
    --from-literal=admin-password="${MEDIASERVER_GRAFANA_PASSWORD}" \
    --dry-run=client -o yaml | kubectl apply -f -
  log_info "Grafana Secret has been created."

  kubectl delete configmap mediaserver-config -n server --ignore-not-found
  kubectl create configmap mediaserver-config -n server --from-env-file=../.env
  log_info "ConfigMap has been created."
}

install_services() {
  log_info "Installing ArgoCD applications..."

  load_env_vars "$(dirname "${BASH_SOURCE[0]}")/../.env"
  require_env_vars "MEDIASERVER_GRAFANA_USERNAME" "MEDIASERVER_GRAFANA_PASSWORD" "MEDIASERVER_CONFIG_DIR"
  ensure_directory "${MEDIASERVER_CONFIG_DIR}/prometheus"

  configs

  kubectl apply -f ./bootstrap/app.yaml
  
  log_info "Applications installed successfully."
}

delete_services() {
  log_info "Removing ArgoCD applications..."
  kubectl delete -f ./bootstrap/app.yaml
  kubectl delete secret grafana-admin-credentials -n monitoring
  kubectl delete configmap mediaserver-config -n server --ignore-not-found
  log_info "Applications removed successfully."
}

cleanup_images() {
  log_info "Removing images..."
  sudo k3s crictl rmi --prune
}

main() {
  if [[ $# -ne 1 ]]; then
    log_error "Invalid number of arguments"
    show_help
    exit 1
  fi

  verify_kube_config

  case "$1" in
    "install")
      install_services
      ;;
    "delete")
      delete_services
      ;;
    "cleanup-images")
      cleanup_images
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