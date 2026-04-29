# Configure Media Server

The Media Server is based on Jellyfin, QBittorrent, Prowlarr and Flaresolverr.

## Setup

Install and configure the setup: `./setup.sh install` - This will create the ArgoCD apps, wait for the applications to be up and running and afterwards configure the services and test the they are configured properly using ArgoCD Hook.

Delete the argocd apps: `./setup.sh delete`

In case we have to be sure argocd is going to fetch new image, cleanup old docker images from k3s using: `./setup.sh cleanup-images`

## Monitoring

Configure Monitoring setup for Raspberry Pi: 

* Setup `Promethes Server` with local storage for storing metrics.
* Export metrics about Raspberri Pi using `Prometheus Node Exporter`.
* Export metrics about Kubernetes using `Kube State Metrics`.

### Helm charts:

Monitoring setup consists of following helm charts:

* [Prometheus](https://artifacthub.io/packages/helm/prometheus-community/prometheus)
* [Grafana](https://artifacthub.io/packages/helm/grafana/grafana)

### Enabled dependencies of Prometheus Helm chart:

* Server
* Node exporter
* Kube state metrics

### Grafana configuration

Grafana is configured in a following way:

**Datasources:**
    
* Prometheus: `http://monitoring-prometheus-server.monitoring.svc.cluster.local:9090`

**Dashboards:**

* [Kubernetes cluster monitoring (via Prometheus)](https://grafana.com/grafana/dashboards/315-kubernetes-cluster-monitoring-via-prometheus/)
* [Node Exporter Full](https://grafana.com/grafana/dashboards/1860-node-exporter-full/)
