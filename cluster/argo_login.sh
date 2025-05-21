#!/usr/bin/env bash

# Log in to ArgoCD

# Get ArgoCD secret, decode it
argoPassword=$(kubectl get secret argocd-initial-admin-secret \
    --output jsonpath="{ .data.password }" \
    --namespace argocd | base64 -d)

# Log in to ArgoCD
argocd login localhost:8080 \
    --username 'admin' \
    --password "$argoPassword" \
    --insecure

# Output the password, for use in different scripts
echo "$argoPassword"