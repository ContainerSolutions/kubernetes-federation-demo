#!/bin/sh
source settings.source

kubectl config use-context federation

kubefed --context=federation unjoin cluster-europe-west1-b --host-cluster-context=gke_${FED_PROJECT}_${FED_HOST_CLUSTER}_gce-${FED_HOST_CLUSTER} --v=8

# kubefed --context=federation join cluster-us-east1-b --cluster-context=gke_${FED_PROJECT}_us-east1-b_gce-us-east1-b --host-cluster-context=gke_${FED_PROJECT}_${FED_HOST_CLUSTER}_gce-${FED_HOST_CLUSTER}

# kubefed --context=federation join cluster-asia-east1-a --cluster-context=gke_${FED_PROJECT}_asia-east1-a_gce-asia-east1-a --host-cluster-context=gke_${FED_PROJECT}_${FED_HOST_CLUSTER}_gce-${FED_HOST_CLUSTER}
