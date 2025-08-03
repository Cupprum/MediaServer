#!/usr/bin/env bash

###############################################################################
# Restart deployments in server namespace
###############################################################################

main() {
    log_info "Retrieving deployment names in the 'server' namespace..."
    local deployments
    deployments=$(kubectl get deployments -n server -o jsonpath='{.items[*].metadata.name}')

    # Check if there are any deployments
    if [ -z "$deployments" ]; then
        log_error "No deployments found in the 'server' namespace."
        exit 1
    fi

    log_info "Restarting all deployments in 'server' namespace..."
    for deployment in $deployments; do
        log_info "Restarting: $deployment"
        kubectl rollout restart "deployment/$deployment" -n server
    done
    log_info "Deployments were restarted"
}

main