#!/usr/bin/env bash

# Create kind cluster
cat cluster.yaml | kind create cluster --config  -