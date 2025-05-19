#!/usr/bin/env bash

# Install k3s
echo "curl -sfL https://get.k3s.io | sh -"

export KUBECONFIG='/etc/rancher/k3s/k3s.yaml'
sudo chmod +r "$KUBECONFIG"

echo 'Read the k3s.yaml kubeconfig, update the ip address to 192.168.0.100'
echo 'Apply this kubeconfig from my local machine'

# TODO: maybe i should wait until k3s is up and running

# Install argocd
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Get argocd secret
argoPassword=$(kubectl get secret \
    argocd-initial-admin-secret \
    -o jsonpath="{ .data.password }" \
    -n argocd | base64 -d)
argocd login localhost:8080 --username 'admin' --password "$argoPassword" --insecure

# Generate a new password for argocd
newArgoPassword=$(date +%s | sha256sum | base64 | head -c 32)
argocd account update-password \
    --current-password "$argoPassword" \
    --new-password "$newArgoPassword"

echo "- New argocd password: $newArgoPassword"

# Add private repo
gitHubToken="${GH_TOKEN:?Error: GH_TOKEN is not set}"
argocd repo add https://github.com/Cupprum/MediaServer.git \
  --name 'MediaServer' \
  --project 'default' \
  --username "$(git config user.name)" \
  --password "$gitHubToken"