# AGENTS.md

## Rules
- Do not read, print, or copy `.env` or `.piactl_login`. They hold secrets, are
  gitignored, and all config flows through `MEDIASERVER_*` env vars.

## What this repo is
Raspberry Pi home media server. A Go program configures Jellyfin/Prowlarr/
qBittorrent over their HTTP APIs, deployed as Helm charts on single-node k3s
with ArgoCD, Prometheus, and Grafana.

## Layout
- `configuration/` — Go module `MediaServer/configuration`. The only code CI
  builds (published to `ghcr.io/cupprum/configure_services`). Entrypoint is the
  root `main.go`; per-service packages are `jellyfin/`, `prowlarr/`, `qbittorrent/`.
- `kube/` — ArgoCD bootstrap (`bootstrap/`), apps-of-apps (`apps/`), and Helm
  charts `service-chart/` (media apps) and `monitoring-chart/` (Prometheus/Grafana).
- `init/` — one-time Pi provisioning (`setup-rpi.sh`, `setup-k8s.sh`), run over SSH.
- `dev/` — docker-compose stack to exercise the Go configurator against live services.
- `utils.sh` — shared bash helpers (`log_*`, `load_env_vars`, `require_env_vars`, …).

## Commands
- Full deploy from a dev machine: `./setup.sh` (SSHes to `x42@pi.local`).
- Run / test configurator locally: `cd configuration && ./configure_services run|test` (loads `../.env`).
- Cluster apps: `cd kube && ./setup.sh install|delete|secrets|cleanup-images`.
- Local service stack: `cd dev && ./compose_recreate [jellyfin|qbittorrent|prowlarr|flaresolverr]`.
- On the Pi: `sudo init/setup-rpi.sh` (root only), then `init/setup-k8s.sh`.

## Re-running config after a code change
The configurator is idempotent (Jellyfin/Prowlarr skip when already configured),
so a plain ArgoCD sync won't re-apply changes to `configuration/`. To force it:
1. `ssh x42@pi.local`
2. delete the relevant config dir, e.g. `/home/x42/MediaServer/config/<service>/config`
   (paths are in `kube/service-chart/values.yaml`)
3. delete the running pod so it restarts against the empty config:
   `kubectl delete pod -n server -l app=<service>`
4. resync the services app: `argocd app sync mediaserver-app` — the PostSync hook
   re-runs `configure_services` + its tests

Service pods live in the `server` namespace (names: jellyfin, prowlarr, qbittorrent, flaresolverr).

## Gotchas
- `configuration/**` "tests" are integration tests hitting live services; they
  need the services up and `MEDIASERVER_*` set. `configure_services test` = `go test ./...`.
- The Docker image compiles each service's tests into `/app/bin/{qbittorrent,prowlarr,jellyfin}.test`;
  an ArgoCD PostSync hook (`kube/service-chart/templates/hook.yaml`) runs them after each sync.
- CI only fires on pushes to `main` touching `configuration/**` (or the workflow) and tags `latest`.
  ArgoCD syncs `targetRevision: HEAD` with self-heal.
- Go module path is `MediaServer/configuration` (not a URL); imports like
  `MediaServer/configuration/utils`. Go 1.25 expected.
- Deployments hardcode `nodeSelector: kubernetes.io/hostname: pi` and `hostPath`
  mounts (`/home/x42/…`, `/media/x42/T5/…`). Pi-specific — don't "fix" casually.
- `configuration/*/cmd/` binaries are unused leftovers (not referenced or built).
- `setup-rpi.sh` has side effects: disables WiFi, enables VNC, installs PIA with
  killswitch. Re-running is not idempotent.

## Style
- Shell scripts follow Google Shell Style (`#!/bin/bash`, `set -euo pipefail`,
  `log_*` helpers). Full guides: `docs/prompts/{go,k8s,shell}.md`.
