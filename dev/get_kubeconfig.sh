#!/usr/bin/env bash

set -e
set -u
set -o pipefail


HOST="pi.local"
IP=$(dig +short pi.local)

ssh "x42@$HOST" 'sudo cat /etc/rancher/k3s/k3s.yaml' > ~/.kube/config
sed --in-place "s/127.0.0.1/$IP/g" ~/.kube/config