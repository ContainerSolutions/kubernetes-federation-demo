#!/usr/bin/env bash

die() { echo "$@" 1>&2 ; exit 1; }
REGION=$1

[[ -z "$REGION" ]] && die "Usage: $0 REGION"

ZONE="$REGION-b"
INSTANCE="traffic-$REGION"

echo "Removing $INSTANCE in $ZONE..."
gcloud compute instances delete $INSTANCE --zone $ZONE -q
