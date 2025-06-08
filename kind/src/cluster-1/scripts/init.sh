#!/bin/sh

# initialized cluster
kind create cluster --config /src/cluster/config.yaml

# apply nginx
kubectl apply -f /src/pods/nginx.yaml

# get nodes
kubectl get nodes -o json

# get pods
kubectl get pods -o json

# describe deployment
kubectl describe deployment
