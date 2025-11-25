#!/usr/bin/env bash

set -e
set -u
set -o pipefail

set -a
source ../configure-services/.env
set +a

envsubst < ./services.yaml.tpl > ./services.yaml