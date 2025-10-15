#!/usr/bin/env bash

set -e
set -u
set -o pipefail

# Configure Kubeconfig
source /workspace/dev/get_kubeconfig.sh

# Open interactive shell
bash