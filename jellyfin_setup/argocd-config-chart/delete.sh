#!/usr/bin/env bash

# Delete script for ArgoCD apps
helm uninstall jellyfin
helm uninstall qbittorrent
helm uninstall prowlarr
helm uninstall heimdall