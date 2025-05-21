#!/usr/bin/env bash

# Linux bootstrap script for Raspberry Pi
# Update and upgrade the system
sudo apt update
sudo apt upgrade -y

# TODO: Is this needed?
# Install docker
curl -sSL https://get.docker.com | sh
sudo usermod -aG docker "$USER"

# Preconfigure k3s - Enable cgroup
echo ' cgroup_enable=memory' | sudo tee -a /boot/cmdline.txt

# Disable wifi
echo "dtoverlay=disable-wifi" | sudo tee -a /boot/firmware/config.txt