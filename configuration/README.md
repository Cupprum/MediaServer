# Server configuration

This project configures the services, by using their APIs.

It sets up QBittorrent, Prowlarr and Jellyfin.

It also supports running the tests, to verify, that the services are actually correctly configured.

## Development

Verify env config file: `../.env`.

**Run locally:**
```
./configure_services run
```

**Run tests:**
```
./configure_services test
```

## Deployment

[![Build and push configure_services](https://github.com/Cupprum/MediaServer/actions/workflows/build_configure_services.yml/badge.svg)](https://github.com/Cupprum/MediaServer/actions/workflows/build_configure_services.yml)

The `build_configure_services.yml` pipeline triggers automatically on change to this folder, and uploads a new latest container with a sha as its tag. ArgoCD afterwards can automatically update the container on next run, because the latest tag is going to point to a different sha.

The outcome of the execution of the pipeline is a tiny scratch docker image is published into [ghcr.io](https://github.com/Cupprum/MediaServer/pkgs/container/configure_services)

## Local Docker

**Build:**
```
docker build -t configure_services .
```

**Run:**
```
docker run configure_services /app/bin/configure_services
```