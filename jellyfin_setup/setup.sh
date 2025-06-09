#!/usr/bin/env bash

# Installation script for ArgoCD apps
helm install jellyfin ./argocd-config-chart --set "service=jellyfin"
helm install qbittorrent ./argocd-config-chart --set "service=qbittorrent"
helm install prowlarr ./argocd-config-chart --set "service=prowlarr"
helm install heimdall ./argocd-config-chart --set "service=heimdall"