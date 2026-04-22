#!/bin/bash

set -e
set -u
set -o pipefail

###############################################################################
# Common Utilities
# Purpose: Shared functions and settings for bash scripts.
###############################################################################

readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1" >&2; }

verify_kube_config() {
  if ! kubectl cluster-info &> /dev/null; then
    log_error "Failed to contact k8s cluster, is KUBECONFIG configured properly?"
    exit 1
  fi
}

load_env_vars() {
  local env_file="$1"
  if [[ -f "${env_file}" ]]; then
    log_info "Loading environment variables from ${env_file}..."
    set -a
    # shellcheck disable=SC1090
    source "${env_file}"
    set +a
  else
    log_error "Environment file ${env_file} not found."
    exit 1
  fi
}

require_env_vars() {
  local missing=()
  local var
  for var in "$@"; do
    if [[ -z "${!var:-}" ]]; then
      missing+=("${var}")
    fi
  done
  if [[ ${#missing[@]} -ne 0 ]]; then
    log_error "Missing required environment variables: ${missing[*]}"
    exit 1
  fi
}

ensure_namespace() {
  local ns="$1"
  if ! kubectl get namespace "${ns}" &> /dev/null; then
    log_info "Creating namespace: ${ns}..."
    kubectl create namespace "${ns}"
  else
    log_info "Namespace ${ns} already exists."
  fi
}

ensure_directory() {
  local dir="$1"
  if [[ ! -d "${dir}" ]]; then
    log_info "Creating directory: ${dir}..."
    mkdir -p "${dir}"
  else
    log_info "Directory ${dir} already exists."
  fi
}