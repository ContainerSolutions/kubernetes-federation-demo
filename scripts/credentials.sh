#!/bin/sh
source settings.source

gcloud -q container clusters get-credentials gce-us-east1-b --zone=us-east1-b --project=${FED_PROJECT}

gcloud -q container clusters get-credentials gce-europe-west1-b --zone=europe-west1-b --project=${FED_PROJECT}

#gcloud -q container clusters get-credentials gce-asia-east1-a --zone=asia-east1-a --project=${FED_PROJECT}
