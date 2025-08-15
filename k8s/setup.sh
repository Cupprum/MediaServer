#!/usr/bin/env bash

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
# shellcheck source=../utils.sh
source "$(dirname "${BASH_SOURCE[0]}")/../utils.sh"

# Constants
readonly KUBECONFIG='/etc/rancher/k3s/k3s.yaml'

###############################################################################
# Functions
###############################################################################

check_prerequisites() {
    if [[ -z "${GH_TOKEN:-}" ]]; then
        log_error "GH_TOKEN environment variable is not set"
        exit 1
    fi
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
    log_info "Setting up ArgoCD..."

   # Expose ArgoCD through Ingress
    kubectl apply --namespace argocd --filename ./ingress.yaml || { log_error "Failed to apply ArgoCD ingress"; exit 1; }

    # Configure ArgoCD server to allow insecure connections
    kubectl patch configmap argocd-cmd-params-cm -n argocd \
        --type merge \
        -p '{"data":{"server.insecure":"true"}}'

    # TODO: fix this hack
    log_info "Start port forwarding argocd to localhost"
    kubectl port-forward svc/argocd-server -n argocd 3000:80 &
    sleep 5

    local argoPassword
    local newArgoPassword

    argoPassword="$(./helpful_scripts/argo_login.sh)" || { log_error "Failed to get ArgoCD password"; exit 1; }
    newArgoPassword="$(date +%s | sha256sum | base64 | head -c 32)"

    argocd account update-password \
        --current-password "$argoPassword" \
        --new-password "$newArgoPassword" || { log_error "Failed to update ArgoCD password"; exit 1; }

    # Update secret with new password
    kubectl get secret argocd-initial-admin-secret \
        --output json \
        --namespace argocd | \
        jq ".data[\"password\"]=\"$newArgoPassword\"" | \
        kubectl apply -f - || { log_error "Failed to update ArgoCD secret"; exit 1; }

    log_info "- New argocd password: $newArgoPassword"
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
    check_prerequisites
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