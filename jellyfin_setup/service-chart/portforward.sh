#!/usr/bin/env bash

# Start the port forwarding
kubectl port-forward svc/jellyfin 8096:8096 &
kubectl port-forward svc/qbittorrent 8080:8080 &
kubectl port-forward svc/prowlarr 9696:9696 &
sudo kubectl port-forward svc/heimdall 80:80 &

echo "Port forwarding started. Press Ctrl+C to stop."
# Wait for user to press Ctrl+C
trap "echo 'Stopping port forwarding...'; pkill -f 'kubectl port-forward'; exit" INT
while true; do
    sleep 1
done