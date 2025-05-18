#!/usr/bin/env bash

# Linux bootstrap script for Raspberry Pi
# Update and upgrade the system
sudo apt update
sudo apt upgrade -y

# Install docker
curl -sSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Install kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
rm kubectl

# Install minikube
curl -LO https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-arm64
sudo install minikube-linux-arm64 /usr/local/bin/minikube && rm minikube-linux-arm64

echo "- NOTE: Please create a new bash session to use docker without sudo."