#!/usr/bin/env bash

# Install k3s
echo "curl -sfL https://get.k3s.io | sh -"

export KUBECONFIG='/etc/rancher/k3s/k3s.yaml'
sudo chmod +x "$KUBECONFIG"