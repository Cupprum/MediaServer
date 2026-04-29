# MediaServer

[![Build and push configure_services](https://github.com/Cupprum/MediaServer/actions/workflows/build_configure_services.yml/badge.svg)](https://github.com/Cupprum/MediaServer/actions/workflows/build_configure_services.yml)

Setup and configure Raspberry Pi to host media center

## Initial setup

This project configures a clean Raspberry PI with kubernetes based mediaserver consisting of Jellyfin, Prowlarr, QBittorrent and Flaresolverr. It also deploys a Grafana and Prometheus for monitoring purposes and sets up DNS and PIA VPN.

## Architecture:

![mediaserver-chart](docs/mediaserver.excalidraw.svg)

## Initial setup:

Install the latest Raspbian OS using `Raspberry Pi Imager.app`.

Pin the desired IP Address of Pi in admin console of router, in the `LAN` setting tab (I like to pin Pi to `192.168.0.100`).

Copy over public SSH key of my main developer device to Pi: `ssh-copy-id -i ~/.ssh/id_rsa.pub x42@pi.local`.

Create a fine grained PAT in GitHub and store it as a `MEDIASERVER_GITHUB_TOKEN` Env var on Raspberry Pi. Make sure its configured to allow read only to `MediaServer` repository.

SSH to Pi: `ssh x42@pi.local` and setup SSH key: `ssh-keygen` and press enter couple of times. Copy the public key and add it to GitHub: `cat ~/.ssh/id_rsa.pub`.

Execute the `setup.sh` shell script located in the root folder of the project on developer computer.
This script will configure the Raspberry Pi, deploy kubernetes cluster into it and spin up all services to run and monitor MediaServer.
