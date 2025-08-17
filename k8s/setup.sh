#!/usr/bin/env bash

set -e
set -u
set -o pipefail

###############################################################################
# K3s and ArgoCD Setup Script
#
# This script automates the installation and configuration of:
# - K3s
# - kubectl
# - Helm
# - ArgoCD
# - Repository setup
###############################################################################

# Source common utilities
source "$(dirname "${BASH_SOURCE[0]}")/../utils.sh"

# Constants
readonly KUBECONFIG='/etc/rancher/k3s/k3s.yaml'

###############################################################################
# Functions
###############################################################################

source_env_vars() {
    log_info "Source environment variables..."
    set -a # automatically export all variables
    source ../.env
    set +a
    if [ -z "$GH_TOKEN" ] || [ -z "$ARGO_PASSWORD" ]; then
        log_error "GH_TOKEN and ARGO_PASSWORD must be present in the environment"
        exit 1
    fi
    log_info "Environment variables loaded successfully"
}

install_k3s() {
    if ! command -v k3s &> /dev/null; then
        log_info "Installing k3s..."
        curl -sfL https://get.k3s.io | sh - || { log_error "Failed to install k3s"; exit 1; }
        log_info "k3s installation completed"
    else
        log_info "k3s is already installed"
    fi
}

install_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        log_info "Installing kubectl..."
        local kubectl_version
        kubectl_version="$(curl -L -s https://dl.k8s.io/release/stable.txt)"
        curl -LO "https://dl.k8s.io/release/${kubectl_version}/bin/linux/arm64/kubectl" || { log_error "Failed to download kubectl"; exit 1; }
        sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && rm kubectl
        log_info "kubectl installation completed"
    else
        log_info "kubectl is already installed"
    fi
}

install_helm() {
    if ! command -v helm &> /dev/null; then
        log_info "Installing Helm..."
        curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash || { log_error "Failed to install Helm"; exit 1; }
        log_info "Helm installation completed"
    else
        log_info "Helm is already installed"
    fi
}

install_argocd_cli() {
    if ! command -v argocd &> /dev/null; then
        log_info "Installing ArgoCD CLI..."
        curl -sSL -o argocd-linux-arm64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-arm64 || { log_error "Failed to download ArgoCD CLI"; exit 1; }
        sudo install -m 555 argocd-linux-arm64 /usr/local/bin/argocd
        rm argocd-linux-arm64
        log_info "ArgoCD CLI installation completed"
    else
        log_info "ArgoCD CLI is already installed"
    fi
}

setup_kubeconfig() {
    log_info "Setting up kubeconfig..."
    export KUBECONFIG
    sudo chmod +r "$KUBECONFIG" || { log_error "Failed to set kubeconfig permissions"; exit 1; }
}

wait_for_k8s() {
    log_info "Waiting for Kubernetes API server to be ready..."
    until kubectl get nodes &>/dev/null; do
        log_info "Waiting for Kubernetes API server..."
        sleep 5
    done
    log_info "Kubernetes API server is ready"
}

install_argocd() {
    log_info "Installing ArgoCD..."
    if ! kubectl get namespace argocd &>/dev/null; then
        log_info "Creating argocd namespace..."
        kubectl create namespace argocd || { log_error "Failed to create argocd namespace"; exit 1; }
        kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml || { log_error "Failed to install ArgoCD"; exit 1; }
    else
        log_info "argocd namespace already exists, skipping creation and configuration"
    fi
}

setup_argocd() {
    # TODO: this should only be executed if Argo NS was created and not configured
    log_info "Setting up ArgoCD..."

   # Expose ArgoCD through Ingress
    kubectl apply --namespace argocd --filename ./ingress.yaml || { log_error "Failed to apply ArgoCD ingress"; exit 1; }

    # Configure ArgoCD server to allow insecure connections
    kubectl patch configmap argocd-cmd-params-cm -n argocd \
        --type merge \
        -p '{"data":{"server.insecure":"true"}}'


    log_info "Waiting for ArgoCD server to be ready..."
    until curl --fail argocd.pi.local; do
        log_info "Waiting for ArgoCD server..."
        sleep 30
    done
    log_info "ArgoCD server is ready"
    
    local currentArgoPassword
    local argoServerPod

    log_info "Getting initial ArgoCD Server password..."
    currentArgoPassword=$(kubectl get secret argocd-initial-admin-secret \
        --output jsonpath="{ .data.password }" \
        --namespace argocd | \
            base64 --decode
    )
    log_info "Current ArgoCD password retrieved successfully"

    log_info "Getting ArgoCD server pod..."
    argoServerPod=$(kubectl get pods \
        --namespace argocd \
        --selector app.kubernetes.io/name=argocd-server \
        --output custom-columns=NAME:.metadata.name \
        --no-headers
    )
    log_info "ArgoCD Server is running on pod: $argoServerPod"

    log_info "Logging into ArgoCD server..."
    kubectl exec "$argoServerPod" \
        --namespace argocd -- \
            argocd login 'localhost:8080' \
                --username 'admin' \
                --password "$currentArgoPassword" \
                --plaintext
    log_info "ArgoCD login successful"

    log_info "Configuring ArgoCD password..."
    kubectl exec "$argoServerPod" \
        --namespace argocd -- \
            argocd account update-password \
                --current-password "$currentArgoPassword" \
                --new-password "$ARGO_PASSWORD"
    log_info "ArgoCD password configured"

    log_info "- ArgoCD UI: http://argocd.pi.local/"
    log_info "- Login with username 'admin' and new password, to update it in password manager"
}

setup_repository() {
    log_info "Adding repository to ArgoCD..."
    argocd repo add https://github.com/Cupprum/MediaServer.git \
        --name 'MediaServer' \
        --project 'default' \
        --username "$(git config user.name)" \
        --password "$GH_TOKEN" || { log_error "Failed to add repository to ArgoCD"; exit 1; }
}

###############################################################################
# Main Script Execution
###############################################################################

main() {
    source_env_vars
    install_k3s
    install_kubectl
    install_helm
    install_argocd_cli
    setup_kubeconfig
    wait_for_k8s
    install_argocd
    setup_argocd
    setup_repository
    log_info "Setup completed successfully"
}

# Execute main function
main