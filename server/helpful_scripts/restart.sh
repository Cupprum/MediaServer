#!/bin/bash

# Get the deployment names in the "server" namespace
deployments=$(kubectl get deployments -n server -o jsonpath='{.items[*].metadata.name}')

# Iterate over the deployments and print them line by line
for deployment in $deployments; do
  echo "$deployment"
done