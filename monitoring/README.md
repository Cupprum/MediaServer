# Configure monitoring

Configure Monitoring setup for Raspberry PI: 

* Setup `Promethes Server` with local storage for storing metrics.
* Export metrics about Raspberri PI using `Prometheus Node Exporter`.
* Export metrics about Kubernetes using `Kube State Metrics`.

## Helm charts:

Monitoring setup consists of following helm charts:

* [Prometheus](https://artifacthub.io/packages/helm/prometheus-community/prometheus)
* [Grafana](https://artifacthub.io/packages/helm/grafana/grafana)

## Enabled dependencies of Prometheus Helm chart:

* Server
* Node exporter
* Kube state metrics

## Grafana configuration

Grafana is configured in a following way:

### Datasources:
    
* Prometheus: `http://monitoring-prometheus-server.monitoring.svc.cluster.local:9090`

### Dashboards:

* [Kubernetes cluster monitoring (via Prometheus)](https://grafana.com/grafana/dashboards/315-kubernetes-cluster-monitoring-via-prometheus/)
* [Node Exporter Full](https://grafana.com/grafana/dashboards/1860-node-exporter-full/)

## Installation
TODO