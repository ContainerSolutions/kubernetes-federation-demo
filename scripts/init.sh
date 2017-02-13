#!/bin/sh
source settings.source

gcloud container clusters create gce-us-east1-b --project=${FED_PROJECT} --zone=us-east1-b --cluster-version=${CLUSTER_VERSION} --scopes cloud-platform,storage-ro,logging-write,monitoring-write,service-control,service-management,https://www.googleapis.com/auth/ndev.clouddns.readwrite

gcloud container clusters create gce-europe-west1-b --project=${FED_PROJECT} --zone=europe-west1-b --cluster-version=${CLUSTER_VERSION} --scopes cloud-platform,storage-ro,logging-write,monitoring-write,service-control,service-management,https://www.googleapis.com/auth/ndev.clouddns.readwrite

gcloud container clusters create gce-asia-east1-a --project=${FED_PROJECT} --zone=asia-east1-a --cluster-version=${CLUSTER_VERSION} --scopes cloud-platform,storage-ro,logging-write,monitoring-write,service-control,service-management,https://www.googleapis.com/auth/ndev.clouddns.readwrite

gcloud container clusters create gce-us-west1-a --project=${FED_PROJECT} --zone=us-west1-a --cluster-version=${CLUSTER_VERSION} --scopes cloud-platform,storage-ro,logging-write,monitoring-write,service-control,service-management,https://www.googleapis.com/auth/ndev.clouddns.readwrite
