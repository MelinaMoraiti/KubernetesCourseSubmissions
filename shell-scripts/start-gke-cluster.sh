#!/bin/bash

set -e

CLUSTER_NAME="dwk-cluster"
ZONE="europe-west8-c"
K8S_VERSION="1.36"
DISK_SIZE="32"
NUM_NODES="4"
MACHINE_TYPE="e2-small"

echo "Creating GKE cluster: ${CLUSTER_NAME}..."

gcloud container clusters create "$CLUSTER_NAME" \
    --zone="$ZONE" \
    --cluster-version="$K8S_VERSION" \
    --disk-size="$DISK_SIZE" \
    --num-nodes="$NUM_NODES" \
    --machine-type="$MACHINE_TYPE"

echo "Enabling Gateway API..."

gcloud container clusters update "$CLUSTER_NAME" \
    --location="$ZONE" \
    --gateway-api=standard

echo "Getting cluster credentials..."

gcloud container clusters get-credentials "$CLUSTER_NAME" \
    --zone="$ZONE"

echo "Cluster is ready."

gcloud container clusters list
kubectl get nodes
kubectl get gatewayclass
