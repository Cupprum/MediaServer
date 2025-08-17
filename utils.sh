#!/usr/bin/env bash

set -e
set -u
set -o pipefail

###############################################################################
# Common Utilities
#
# Shared functions and settings for bash scripts
# Meant to be sourced by other scripts
###############################################################################

# Color definitions
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

verify_kube_config() {
    if ! kubectl cluster-info &> /dev/null; then
        log_error "Failed to contact k8s cluster, is KUBECONFIG configured properly?"
        exit 1
    fi
}