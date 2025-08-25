#!/usr/bin/env bash

set -e
set -u
set -o pipefail

###############################################################################
# Raspberry Pi Media Server Installation Script
#
# This script configures MediaServer on Raspberry Pi, specifically:
# - Configure the Raspberry Pi
# - Setup and configure Kubernetes
# - Deploy MediaServer to Kubernetes
#
# Usage: ./setup.sh
###############################################################################

# Source common utilities
source "$(dirname "${BASH_SOURCE[0]}")/utils.sh"

Hostname='x42@pi.local'

log_info 'Checking if MediaServer folder exists...'
if ssh "$Hostname" '! test -d MediaServer'; then
    log_info 'MediaServer folder not found, cloning repository...'

    ssh "$Hostname" 'git clone git@github.com:Cupprum/MediaServer.git'
    if [ $? -eq 128 ]; then
        log_info 'Could not clone the repository, adding GitHub public ssh keys to known hosts...'
        ssh "$Hostname" 'ssh-keyscan github.com >> ~/.ssh/known_hosts 2>/dev/null'
        log_info 'Cloning the repository again'
        ssh "$Hostname" 'git clone git@github.com:Cupprum/MediaServer.git'
    fi
    log_success 'Repository cloned successfully'
else
    log_info 'MediaServer folder already exists'
fi

log_info 'Configuring Raspberry Pi...'
ssh "$Hostname" 'cd MediaServer/rpi; sudo ./setup.sh'

log_info 'Restarting the Raspberry Pi...'
ssh "$Hostname" 'sudo reboot'

# Note: the next ssh command would be waiting there either way
log_info 'Going to sleep for 120 seconds, to wait for Raspberry to restart'
sleep 120

log_info 'Seting up and configuring Kubernetes...'
ssh "$Hostname" 'cd MediaServer/k8s; ./setup.sh'

log_info 'Deploying MediaServer to Kubernetes...'
ssh "$Hostname" 'cd MediaServer/server; ./bootstrap.sh'
ssh "$Hostname" 'cd MediaServer/server; ./setup.sh install'

log_info 'Deploying Monitoring to Kubernetes...'
ssh "$Hostname" 'cd MediaServer/monitoring; ./setup.sh install'