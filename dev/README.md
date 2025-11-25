# Devcontainer

Development container used for development of this project. Keeps dependencies separated.

## List of commands:

* Build new devcontainer: `./devcontainer build`
* Start interactive shell in the container: `./devcontainer run`

# Docker Compose

For testing the `configure-services` go code, we can simply run `docker compose up -d --force-recreate`.

This will run all services in docker compose.

# Homepage Service Widgets Configuration

Homepage dashboard of services is configured by the `services.yaml` stored in `/app/config/services.yaml`. This file can be generated using the `generate_services.yaml.sh` shell script.