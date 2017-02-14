#!/usr/bin/env bash

die() { echo -e "$@" 1>&2 ; exit 1; }

REGION=$1
INGRESS_IP="$2"
ADMIN_IP="$3"
SKIP_CREATE="$4"

[[ -z "$INGRESS_IP" || -z "$ADMIN_IP" || -z "$REGION" ]] && die "Usage: $0 REGION INGRESS_IP ADMIN_IP [skip-create]\n\n    REGION e.g. us-central1 - should be different from regions with clusters deployed\n"

ZONE="$REGION-b" # b seems to exist in all regions...
INSTANCE="traffic-$REGION"
SCRIPT_NAME="gen-traffic.sh"
SCRIPT_SRC="$(dirname $0)/$SCRIPT_NAME"

if [ "$SKIP_CREATE" == "skip-create" ]; then
    echo "Skipping creation of instance $INSTANCE..."
else
    echo "Creating instance $INSTANCE..."
    gcloud compute instances create $INSTANCE --image-family ubuntu-1604-lts --image-project ubuntu-os-cloud --zone $ZONE --machine-type=f1-micro -q
    test $? -ne 0 && die "Unable to create instance"
    sleep 10
fi

echo
echo "Copying $SCRIPT_SRC to $INSTANCE..."
gcloud compute copy-files $SCRIPT_SRC $INSTANCE:~/ --zone $ZONE -q
test $? -ne 0 && die "Unable to copy traffic script to instance"

echo
echo "Installing Apache Bench..."
gcloud compute ssh $INSTANCE --zone $ZONE -q -- "sudo apt-get install apache2-utils -y"
test $? -ne 0 && die "Unable to install ab"

echo
echo "Launching traffic script (admin: $ADMIN_IP, ingress: $INGRESS_IP)..."
gcloud compute ssh $INSTANCE --zone $ZONE -q -- "nohup ~/$SCRIPT_NAME $INGRESS_IP $ADMIN_IP >/tmp/traffic.log 2>&1 &"
test $? -ne 0 && die "Unable to launch traffic script"

echo
echo "****************************************************************************"
echo "Traffic script running. You can also log into the VM for debugging purposes:"
echo "gcloud compute ssh $INSTANCE --zone $ZONE"
echo "****************************************************************************"