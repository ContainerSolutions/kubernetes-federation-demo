#!/usr/bin/env bash
# This script must run from within a VM inside Google Compute Engine.

die() { echo "$@" 1>&2 ; exit 1; }

INGRESS_IP="$1"
ADMIN_IP="$2"

[[ -z "$INGRESS_IP" || -z "$ADMIN_IP" ]] && die "Usage: $0 INGRESS_IP ADMIN_IP"

which ab || die "Apache Bench not installed."

echo "Generating traffic to $INGRESS_IP (admin: $ADMIN_IP) ..."
PROTO_ZONE=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/zone" -H "Metadata-Flavor: Google") || die "unable to get zone"
REGION=$(echo $PROTO_ZONE | rev | cut -d/ -f1 | rev | cut -d- -f1-2)
HEADER="X-Origin-Region: $REGION"

while true; do
    active=`curl http://$ADMIN_IP/trafficSourceActive?dc=$REGION 2>/dev/null`
    if [ $? -ne 0 ]; then
        echo "Error asking admin if I should generate traffic." 2>&1
        sleep 2
    else
        if echo $active | grep "true" >/dev/null; then 
            ab -t 2 -H "$HEADER" http://$INGRESS_IP/ >/dev/null 2>&1
        else
            echo "I'm currently not supposed to generate traffic - waiting 2 seconds..."
            sleep 2
        fi
    fi
done
