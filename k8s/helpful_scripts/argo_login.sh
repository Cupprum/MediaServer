#!/usr/bin/env bash

###############################################################################
# ArgoCD Login Script
#
# This script:
# - Retrieves the ArgoCD initial admin password
# - Logs in to ArgoCD
# - Outputs the password for use in other scripts
###############################################################################

# Source common utilities
# shellcheck source=../../utils.sh
source "$(dirname "${BASH_SOURCE[0]}")/../../utils.sh"

###############################################################################
# Functions
###############################################################################

get_argocd_password() {
    local password
    password=$(kubectl get secret argocd-initial-admin-secret \
        --output jsonpath="{ .data.password }" \
        --namespace argocd) || { log_error "Failed to get ArgoCD secret"; exit 1; }
    
    if [[ -z "$password" ]]; then
        log_error "ArgoCD password is empty"
        exit 1
    fi
    
    echo "$password" | base64 -d
}

login_argocd() {
    local password="$1"
    log_info "Logging into ArgoCD..."
    argocd login 127.0.0.1:3000 \
        --username 'admin' \
        --password "$password" \
        --insecure > /dev/null || { log_error "Failed to login to ArgoCD"; exit 1; }
    log_info "Successfully logged into ArgoCD"
}

###############################################################################
# Main Script Execution
###############################################################################

main() {
    local argoPassword
    
    argoPassword=$(get_argocd_password)
    login_argocd "$argoPassword"
    
    # Output password for use in other scripts
    echo -n "$argoPassword"
}

# Execute main function
main