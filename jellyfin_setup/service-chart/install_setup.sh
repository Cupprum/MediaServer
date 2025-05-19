#!/usr/bin/env bash

# Install the app

# helm install jellyfin -f ./values/jellyfin.yaml .
# helm install qbittorrent -f ./values/qbittorrent.yaml .
# helm install prowlarr -f ./values/prowlarr.yaml .
# helm install heimdall -f ./values/heimdall.yaml .

argocd app create jellyfin \
    --repo 'https://github.com/Cupprum/MediaServer.git' \
    --path 'jellyfin_setup/service-chart' \
    --dest-server 'https://kubernetes.default.svc' \
    --values 'values/jellyfin.yaml' \
    --dest-namespace 'default'