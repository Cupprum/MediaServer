# Devcontainer

Development container used for development of this project. Keeps dependencies separated.

## List of commands:

* Build new devcontainer: `./devcontainer build`
* Start interactive shell in the container: `./devcontainer run`

# Docker Compose

For testing the `server/configuration` go code, we can run the services in docker compose using the following utility: `./compose_recreate`, it also supports just recreating a single service: `./compose_recreate jellyfin`.