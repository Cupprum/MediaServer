#!/bin/bash

set -e
set -u
set -o pipefail

###############################################################################
# Raspberry Pi Media Server Installation Script
# Purpose: Configures MediaServer on Raspberry Pi, K8s, and deployments.
# Usage: ./setup.sh
###############################################################################

source "$(dirname "${BASH_SOURCE[0]}")/utils.sh"

readonly TARGET_HOST='x42@pi.local'
readonly REPO_URL='git@github.com:Cupprum/MediaServer.git'

main() {
  log_info "Checking if MediaServer folder exists on ${TARGET_HOST}..."
  
  if ssh "${TARGET_HOST}" '[[ ! -d MediaServer ]]'; then
    log_info "MediaServer folder not found. Cloning repository..."

    ssh "${TARGET_HOST}" "git clone ${REPO_URL}"
    if [ $? -eq 128 ]; then
      log_warn "Could not clone repository. Adding GitHub to known hosts..."
      ssh "${TARGET_HOST}" 'ssh-keyscan github.com >> ~/.ssh/known_hosts 2>/dev/null'
      
      log_info "Retrying repository clone..."
      ssh "${TARGET_HOST}" "git clone ${REPO_URL}"
    fi
    log_info "Repository cloned successfully."

    log_info "Copying secrets to MediaServer..."
    scp .env "${TARGET_HOST}:~/MediaServer/.env"
    scp .piactl_login "${TARGET_HOST}:~/MediaServer/.piactl_login"
    log_info "Secrets copied."
  else
    log_info "MediaServer folder already exists."
  fi

  log_info "Configuring Raspberry Pi..."
  ssh "${TARGET_HOST}" 'cd MediaServer/rpi && sudo ./setup.sh'

  log_info "Restarting the Raspberry Pi..."
  ssh "${TARGET_HOST}" 'sudo reboot' || true

  # Note: the next ssh command would be waiting there either way
  log_info "Sleeping for 120 seconds to allow Raspberry Pi to reboot..."
  sleep 120

  log_info "Setting up and configuring Kubernetes..."
  ssh "${TARGET_HOST}" 'cd MediaServer/k8s && ./setup.sh'

  log_info "Deploying MediaServer to Kubernetes..."
  ssh "${TARGET_HOST}" 'cd MediaServer/server && ./setup.sh install'

  log_info "Deploying Monitoring to Kubernetes..."
  ssh "${TARGET_HOST}" 'cd MediaServer/monitoring && ./setup.sh install'
  
  log_info "MediaServer installation completed successfully!"
}

main