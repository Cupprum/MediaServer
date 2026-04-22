# Docker Compose

For testing the `server/configuration` go code, we can run the services in docker compose using the following utility: `./compose_recreate`, it also supports just recreating a single service: `./compose_recreate jellyfin`.

# Docker jail for gemini cli

Executing `./docker.sh build` and `./docker.sh run` spawn a fedora based docker container, which has gemini cli preinstalled. 