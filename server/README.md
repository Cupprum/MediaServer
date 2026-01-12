# Configure Media Server

The Media Server is based on Jellyfin, QBittorrent, Prowlarr and Homepage.

## Installation and configuration

Install setup using the following command: `./setup.sh install` - This will create the ArgoCD apps, wait for the applications to be up and running and afterwards configure the services and test the they are configured properly.

## Deletion process

Delete the argocd apps with their configuration: `./setup.sh delete`