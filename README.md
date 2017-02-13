# Kubernetes Federation Demo

## Requirements

### Kubernetes command line tools

Download the Kubernetes command line tools from `https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG.md#client-binaries` and 
install the them in your `$PATH`.

### Generate GKE secret key

Generate a Google platform secret key in JSON format and store it in `$HOME/.federation-key.json`.

## Preparing the demo

### Setup GKE clusters in different regions

**Warning:** this creates several clusters in Google Cloud Platform! Watch out for your billing! /!\

#### Build and push the docker image for the demo application

1. make push

#### Create the clusters and deploy the demo application

1. make sure `kubefed` and `kubectl` are in `$PATH`

2. Reserve a global IP adddress

    ```
    gcloud compute addresses create kubernetes-ingress --global
    ```

2. Create the clusters and deploy the manifests

    ```
    scripts/create.sh
    ```   

#### Explore the demo

Point your browser to the IP address you generated in step 2.
You can always find it again later by executing:

    gcloud compute addresses list | grep kubernetes-ingress

You should see the clusters appearing on a map, but no traffic yet.

#### Generate traffic

The script `scripts/start-traffic.sh` will create a micro-instance in the desired Google 
datacenter and execute a small script that generates traffic on the federated cluster.
Launch is as often as you want with different regions as argument.

**Important:** Do not launch it in the same regions as the Kubernetes clusters, as this will not
work nicely with the Maps user interface.

The required arguments are the region where to generate traffic from, the IP address 
of the global ingress load balancer (from step 2) and the IP address of the admin service (FIXME where do you find it easily?).

    # replace the two IP addresses with your own:
    scripts/start-traffic.sh asia-northeast1 130.211.41.245 104.155.43.73
    scripts/start-traffic.sh us-central1 130.211.41.245 104.155.43.73

To remove the VM generating traffic again, run `scripts/stop-traffic.sh REGION`. This 
may also be useful to clean up when setting up traffic generation might have failed for some reason.

## Demo time

You can now move around the map, and enable or disable either clusters or traffic sources (green markers with a house icon).
You should see the traffic numbers in the statistics window changing and lines between traffic sources and datacenters, 
depending on what the next healthy datacenter to a traffic source is.

## Known Issues

Some firewall rules need to be setup:

    # FIXME - this is wrong, ports need to be determined with kubectl after the federation has been setup
    gcloud compute firewall-rules create my-federated-ingress-firewall-rule --source-ranges 130.211.0.0/22 --allow "icmp,tcp:80,tcp:443" --target-tags "cluster-europe-west1-b,cluster-asia-east1-a,cluster-us-east1-b" --network default

See also:

- https://github.com/kubernetes/kubernetes/issues/37306
