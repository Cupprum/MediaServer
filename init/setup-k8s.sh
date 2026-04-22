#!/bin/bash

set -e
set -u
set -o pipefail

###############################################################################
# K3s and ArgoCD Setup Script
# Purpose: Automates the installation and configuration of K3s and ArgoCD.
###############################################################################

source "$(dirname "${BASH_SOURCE[0]}")/../utils.sh"

readonly KUBE_CONFIG='/etc/rancher/k3s/k3s.yaml'

install_k3s() {
  if ! command -v k3s &> /dev/null; then
    log_info "Installing k3s..."
    curl -sfL https://get.k3s.io | sh -
    log_info "k3s installation completed."
  else
    log_info "k3s is already installed."
  fi
}

install_kubectl() {
  if ! command -v kubectl &> /dev/null; then
    log_info "Installing kubectl..."
    local kubectl_version
    kubectl_version="$(curl -L -s https://dl.k8s.io/release/stable.txt)"
    curl -LO "https://dl.k8s.io/release/${kubectl_version}/bin/linux/arm64/kubectl"
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    rm kubectl
    log_info "kubectl installation completed."
  else
    log_info "kubectl is already installed."
  fi
}

install_helm() {
  if ! command -v helm &> /dev/null; then
    log_info "Installing Helm..."
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
    log_info "Helm installation completed."
  else
    log_info "Helm is already installed."
  fi
}

install_argocd_cli() {
  if ! command -v argocd &> /dev/null; then
    log_info "Installing ArgoCD CLI..."
    curl -sSL -o argocd-linux-arm64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-arm64
    sudo install -m 555 argocd-linux-arm64 /usr/local/bin/argocd
    rm argocd-linux-arm64
    log_info "ArgoCD CLI installation completed."
  else
    log_info "ArgoCD CLI is already installed."
  fi
}

setup_kubeconfig() {
  log_info "Setting up kubeconfig permissions..."
  export KUBECONFIG="${KUBE_CONFIG}"
  sudo chmod +r "${KUBECONFIG}"
}

wait_for_k8s() {
  log_info "Waiting for Kubernetes API server to be ready..."
  until kubectl get nodes &>/dev/null; do
        log_info "Waiting for Kubernetes API server..."
    sleep 5
  done
  log_info "Kubernetes API server is ready."
}

create_and_configure_argocd() {
  ensure_namespace "argocd"

  log_info "Applying ArgoCD manifests..."
  kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

  log_info "Exposing ArgoCD through Ingress..."
  kubectl apply --namespace argocd --filename ./ingress.yaml

  log_info "Configuring ArgoCD server to allow insecure connections..."
  kubectl patch configmap argocd-cmd-params-cm -n argocd \
    --type merge \
    -p '{"data":{"server.insecure":"true"}}'

  log_info "Waiting for ArgoCD server to be ready..."
    until curl --fail argocd.pi.local; do
        log_info "Waiting for ArgoCD server..."
        sleep 30
  done
  log_info "ArgoCD server is reachable."

  local current_argo_password
  local argo_server_pod

  log_info "Retrieving initial ArgoCD Server password..."
  current_argo_password=$(kubectl get secret argocd-initial-admin-secret \
    --output jsonpath="{ .data.password }" \
    --namespace argocd | base64 --decode)

  argo_server_pod=$(kubectl get pods \
    --namespace argocd \
    --selector app.kubernetes.io/name=argocd-server \
    --output custom-columns=NAME:.metadata.name \
    --no-headers)
  
  log_info "Logging into ArgoCD server..."
    kubectl exec "$argo_server_pod" \
        --namespace argocd -- \
            argocd login 'localhost:8080' \
                --username 'admin' \
                --password "$current_argo_password" \
                --plaintext

  log_info "Updating ArgoCD password..."
    kubectl exec "$argo_server_pod" \
        --namespace argocd -- \
            argocd account update-password \
                --current-password "$current_argo_password" \
                --new-password "$MEDIASERVER_ARGO_PASSWORD"

  log_info "Adding MediaServer repository to ArgoCD..."
  kubectl exec "${argo_server_pod}" --namespace argocd -- \
    argocd repo add https://github.com/Cupprum/MediaServer.git \
      --name 'MediaServer' \
      --project 'default' \
      --username "$(git config user.name)" \
      --password "${MEDIASERVER_GITHUB_TOKEN}"

  log_info "- ArgoCD UI: http://argocd.pi.local/"
    log_info "- Login with username 'admin' and new password, to update it in password manager"
}

setup_argocd() {
  if kubectl get namespace argocd &>/dev/null; then
    log_info "argocd namespace already exists, skipping creation and configuration."
  else
    create_and_configure_argocd
  fi
}

main() {
  load_env_vars "$(dirname "${BASH_SOURCE[0]}")/../.env"
  require_env_vars "MEDIASERVER_GITHUB_TOKEN" "MEDIASERVER_ARGO_PASSWORD"

  install_k3s
  install_kubectl
  install_helm
  install_argocd_cli
  setup_kubeconfig
  wait_for_k8s
  setup_argocd
  
  log_info "K8s and ArgoCD setup completed successfully."
}

main