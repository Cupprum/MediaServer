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
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && rm kubectl

# Install helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Preconfigure k3s - Enable cgroup
echo ' cgroup_enable=memory' | sudo tee -a /boot/cmdline.txt

echo "- NOTE: Please create a new bash session to use docker without sudo."