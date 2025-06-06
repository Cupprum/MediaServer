#!/usr/bin/env bash

# Check for required environment variables
GH_TOKEN="${GH_TOKEN:?Error: GH_TOKEN is not set}"

# Check if k3s is installed
if ! command -v k3s &> /dev/null; then
    echo "k3s not found, proceeding with installation..."

    echo "Install k3s"
    echo "curl -sfL https://get.k3s.io | sh -"

    echo "Install kubectl"
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubectl"
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && rm kubectl

    echo "Install helm"
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

    echo "Install ArgoCD cli"
    curl -sSL -o argocd-linux-arm64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-arm64
    sudo install -m 555 argocd-linux-arm64 /usr/local/bin/argocd
    rm argocd-linux-arm64
fi

export KUBECONFIG='/etc/rancher/k3s/k3s.yaml'
sudo chmod +r "$KUBECONFIG"

# TODO: maybe i should wait until k3s is up and running

# Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd \
    -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# ArgoCD login
argoPassword=$(./argo_login.sh)

# Generate a new password for ArgoCD
newArgoPassword=$(date +%s | sha256sum | base64 | head -c 32)
argocd account update-password \
    --current-password "$argoPassword" \
    --new-password "$newArgoPassword"

echo "- New argocd password: $newArgoPassword"
echo "- ArgoCD UI: http://localhost:8080"
echo "- Login with username 'admin' and new password, to update it in password manager"

# Set the new password in the argocd secret
kubectl get secret argocd-initial-admin-secret \
    --output json \
    --namespace argocd | \
    jq ".data[\"password\"]=\"$newArgoPassword\"" | \
    kubectl apply -f -

# Add private repo
argocd repo add https://github.com/Cupprum/MediaServer.git \
  --name 'MediaServer' \
  --project 'default' \
  --username "$(git config user.name)" \
  --password "$GH_TOKEN"