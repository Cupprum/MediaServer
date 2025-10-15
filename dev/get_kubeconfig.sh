#!/usr/bin/env bash

set -e
set -u
set -o pipefail

scp 'x42@pi.local:/etc/rancher/k3s/k3s.yaml' ~/.kube/config
sed --in-place 's/127.0.0.1/192.168.0.100/g' ~/.kube/config