#!/bin/bash

set -e

CLUSTER_NAME="dwk-cluster"
ZONE="europe-west8-c"
gcloud container clusters delete "$CLUSTER_NAME" --zone="$ZONE"

