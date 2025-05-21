#!/usr/bin/env bash

# Installation script for ArgoCD apps
helm install jellyfin . --set "service=jellyfin"
helm install qbittorrent . --set "service=qbittorrent"
helm install prowlarr . --set "service=prowlarr"
helm install heimdall . --set "service=heimdall"