# Configure Media Server

The Media Server is based on Jellyfin, QBittorrent, Prowlarr and Flaresolverr.

## Installation and configuration

Install and configure the setup using the following command: `./setup.sh install_and_configure` - This will create the ArgoCD apps, wait for the applications to be up and running and afterwards configure the services and test the they are configured properly.

This tool also supports other commands, like only installation only mode: `./setup.sh install`, or configuration only mode: `./setup.sh configure`

## Deletion process

Delete the argocd apps with their configuration: `./setup.sh delete`

Cleanup config folders of the services: `./setup.sh cleanup`