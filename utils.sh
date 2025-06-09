#!/usr/bin/env bash

###############################################################################
# Common Utilities
#
# Shared functions and settings for bash scripts
# Meant to be sourced by other scripts
###############################################################################

set -e # Exit on fail
set -u # Treat unset variables as an error
set -o pipefail # Fail if any command in a pipeline fails

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

