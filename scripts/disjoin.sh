#!/bin/sh
source settings.source

die() { echo "$@" 1>&2 ; exit 1; }

[[ -z "$1" ]] && die "Usage: $0 CLUSTER (e.g. cluster-europe-west1-b)"

kubectl config use-context federation

kubefed --context=federation unjoin $1 --host-cluster-context=gke_${FED_PROJECT}_${FED_HOST_CLUSTER}_gce-${FED_HOST_CLUSTER} --v=8
