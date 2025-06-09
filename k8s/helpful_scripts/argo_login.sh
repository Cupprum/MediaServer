#!/usr/bin/env bash

###############################################################################
# ArgoCD Login Script
#
# This script:
# - Retrieves the ArgoCD initial admin password
# - Logs in to ArgoCD
# - Outputs the password for use in other scripts
###############################################################################

set -e # Exit on fail
set -u # Treat unset variables as an error
set -o pipefail # Fail if any command in a pipeline fails

# Color definitions
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly NC='\033[0m' # No Color

###############################################################################
# Functions
###############################################################################

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

get_argocd_password() {
    local password
    password=$(kubectl get secret argocd-initial-admin-secret \
        --output jsonpath="{ .data.password }" \
        --namespace argocd) || { log_error "Failed to get ArgoCD secret"; exit 1; }
    
    if [[ -z "$password" ]]; then
        log_error "ArgoCD password is empty"
        exit 1
    fi
    
    echo "$password"
}

login_argocd() {
    local password="$1"
    log_info "Logging into ArgoCD..."
    argocd login localhost:8080 \
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