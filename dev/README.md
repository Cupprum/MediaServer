# Devcontainer

Development container used for development of this project. Keeps dependencies separated.

## List of commands:

* Build new devcontainer: `./devcontainer build`
* Start interactive shell in the container: `./devcontainer run`

# Docker Compose

For testing the `configure-services` we can run the services in docker compose: `docker compose up -d --force-recreate`.

# Homepage Service Widgets Configuration

Homepage dashboard of services is configured by the `services.yaml` stored in `/app/config/services.yaml`. This file can be generated using the `generate_services.yaml.sh` shell script.