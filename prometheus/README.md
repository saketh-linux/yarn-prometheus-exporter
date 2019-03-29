## Prometheus Monitoring On Kubernetes

## Create A Namespace
    kubectl create namespace monitoring
    
## Cluster Role
    kubectl create -f clusterRole.yaml
    
## Create A Config Map
    We should create a config map with all the prometheus scrape config and alerting rules, 
    which will be mounted to the Prometheus container in /etc/prometheus as prometheus.yaml 
    and prometheus.rules files. The prometheus.yaml contains all the configuration to dynamically
    discover pods and services running in the kubernetes cluster. prometheus.rules will contain all 
    the alert rules for sending alerts to alert manager.
    
    kubectl create -f config-map.yaml -n monitoring
    
## Create A Prometheus Deployment   
    kubectl create  -f prometheus-deployment.yaml --namespace=monitoring
    
## Exposing Prometheus As A Service
    To access the Prometheus dashboard over a IP or a DNS name, you need to expose it as 
    kubernetes service.
    
    kubectl create -f prometheus-service.yaml --namespace=monitoring
    Once created, you can access the Prometheus dashboard using any Kubernetes node IP on port 30090.