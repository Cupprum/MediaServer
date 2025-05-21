#!/usr/bin/env bash

# Installation script for ArgoCD
gitHubToken="${GH_TOKEN:?Error: GH_TOKEN is not set}"
helm install argocd . \
    --set "git.token=$gitHubToken"