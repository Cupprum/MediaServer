#!/usr/bin/env bash

# Configure the system to use k3s
export KUBECONFIG='/etc/rancher/k3s/k3s.yaml'
sudo chmod +r "$KUBECONFIG"

echo "- Generating kubeconfig to use with kubectl from other machine"
sed 's/127.0.0.1/192.168.0.100/g' "$KUBECONFIG"
cat "$KUBECONFIG"