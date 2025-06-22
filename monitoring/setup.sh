# Install Prometheus

KUBECONFIG='/etc/rancher/k3s/k3s.yaml'
sudo chmod +r "$KUBECONFIG"

# Install Prometheus
cd './monitoring-chart' || exit 1
helm dependency build

helm install monitoring . \
    --namespace monitoring \
    --create-namespace

helm uninstall monitoring \
    --namespace monitoring

# Get Grafana admin password
kubectl get secrets monitoring-grafana \
    --namespace monitoring \
    --output jsonpath="{.data.admin-password}" | base64 -d