# Cluster

Manage k3s kubernetes cluster.

## Configure kubectl

Execute on Raspberry Pi; to print the Kubeconfig. The script updates the IP Address of the cluster. Copy the content of Kubeconfig to different machine on network. Don't forget to specify the location of kubeconfig before comunicating with kubernetes.

Cat and update Kubeconfig: `./configure.sh`

## Cluster creation

Install things running on Kubernetes and configure ArgoCD.

NOTE: This script is meant to be executed only once when initially creating the k3s cluster.

Create and configure new cluster: `./create_cluster.sh`
