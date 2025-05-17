#!/usr/bin/env bash

# Delete the app
helm uninstall jellyfin
helm uninstall qbittorrent
helm uninstall prowlarr
helm uninstall heimdall