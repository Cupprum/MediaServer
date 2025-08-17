# Cluster

Manage k3s kubernetes cluster.

## Cluster creation and configuration

Install tooling around Kubernetes and configure ArgoCD.

Create and configure new k3s cluster: `cd ./k8s; ./setup.sh`

## Helpful scripts

Collection of scripts, used to simplify managing the server.

### Get Kubeconfig

Execute on Raspberry Pi; to print the Kubeconfig. The script updates the IP Address of the cluster. Copy the content of Kubeconfig to different machine on network. Don't forget to specify the location of kubeconfig before comunicating with kubernetes.

Cat and update Kubeconfig: `cd ./k8s/helpful_scripts; ./get_kubeconfig.sh`