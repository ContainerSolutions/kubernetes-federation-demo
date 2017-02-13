#!/bin/sh
source settings.source

gcloud -q container clusters delete gce-us-east1-b --project=${FED_PROJECT} --zone=us-east1-b
gcloud -q container clusters delete gce-europe-west1-b --project=${FED_PROJECT} --zone=europe-west1-b
# gcloud -q container clusters delete gce-asia-east1-a --project=${FED_PROJECT} --zone=asia-east1-a
