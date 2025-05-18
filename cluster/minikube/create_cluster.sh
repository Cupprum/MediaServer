#!/usr/bin/env bash

# Start the cluster
minikube start

# Enable local storage path provisioner
minikube addons enable storage-provisioner-rancher