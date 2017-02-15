# Kubernetes Federation Demo

## Requirements

### Kubernetes command line tools

Download the command line tools from [Github](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG.md#client-binaries) and 
install them in your `$PATH` (e.g. `/usr/local/sbin`).

### GCP project

Create a new project in your [Google Cloud Console](https://console.cloud.google.com/) and note down the project name.

### GKE secret key

Generate a Google Cloud Platform secret key in JSON format and store it in this location: `$HOME/.federation-key.json`

### DockerHub account

If you haven't got a DockerHub account yet, create one at [hub.docker.com](https://hub.docker.com/).

## Preparing the demo

**WARNING:** this creates several clusters in Google Cloud Platform! Watch out for your billing! /!\

### Build and push the docker image for the demo application

Adjust the variables `FED_PROJECT`, `FED_DNS_ZONE` and `IMAGE` in the `Makefile` and in `scripts/settings.source` to match your GCE project name, DNS configuration and your DockerHub account.

Then execute: `make push`

### Create the clusters and deploy the demo application

1. make sure `kubefed` and `kubectl` are in `$PATH`

2. Reserve a couple of global IP adddress

    ```
    gcloud compute addresses create kubernetes-ingress --global
    gcloud compute addresses create geoserver-admin --global
    ```

2. Create the clusters: (WARNING: this operation takes a long time)
    Change folder to `scripts`

    ```
    ./init.sh
    ```   

3. Initialise the federation:    

    ```
    ./init-federation.sh
    ```   

4. Join the clusters together:

    ```
    ./join.sh
    ```   


5. Deploy Federated Service and Ingress and the federated application

    ```
    kubectl --context=federation create -f manifests/geoserver-service.yaml
    kubectl --context=federation create -f manifests/geoserver-ingress.yaml
    kubectl --context=federation create -f manifests/geoserver-replica.yaml
    ```

6. Deploy the map

    ```
    kubectl --context=federation create -f manifests/geoserver-admin.yaml    
    ```

### Explore the demo

Point your browser to the IP address you generated in step 2 (geoserver-admin).
You can always find it again later by executing:

    `gcloud compute addresses list | grep geoserver-admin`

You should see the clusters appearing on a map, but no traffic yet.

### Generate traffic

The script `scripts/start-traffic.sh` will create a micro-instance in the desired Google 
datacenter and execute a small script that generates traffic on the federated cluster.
Launch is as often as you want with different regions as argument.

**Important:** Do not launch it in the same regions as the Kubernetes clusters, as this will not
work nicely with the Maps user interface.

The required arguments are the region where to generate traffic from, the IP address 
of the global ingress load balancer and the IP address of the admin service (from step 2 above).

    # replace the two IP addresses with your own created in step 2:
    scripts/start-traffic.sh asia-northeast1 130.211.41.245 104.155.43.73
    scripts/start-traffic.sh us-central1 130.211.41.245 104.155.43.73

To remove the VM generating traffic again, run `scripts/stop-traffic.sh REGION`. This 
may also be useful to clean up when setting up traffic generation might have failed for some reason.

## Demo time

You can now move around the map, and enable or disable either clusters or traffic sources (green markers with a house icon).
You should see the traffic numbers in the statistics window changing and lines between traffic sources and datacenters, 
depending on what the next healthy datacenter to a traffic source is.

## Known Issues

Some firewall rules need to be setup manually:

First step is to retrieve the ports to open via:

- `kubectl --context=gke_steam-ego-156812_us-east1-b_gce-us-east1-b --namespace=kube-system get services`
- `kubectl --context=gke_steam-ego-156812_europe-west1-b_gce-europe-west1-b --namespace=kube-system get services`
- `kubectl --context=gke_steam-ego-156812_asia-east1-a_gce-asia-east1-a --namespace=kube-system get services`

The ports to open are listed under: `default-http-backend` service.
It usually is higher than 30000.

Then run the following command: (please change the ports accordingly)    
    gcloud compute firewall-rules create my-federated-ingress-firewall-rule --source-ranges 130.211.0.0/22 --allow "icmp,tcp:80,tcp:443,tcp:30451,tcp:31014,tcp:30699" --target-tags "cluster-europe-west1-b,cluster-asia-east1-a,cluster-us-east1-b" --network default

See also:

- https://github.com/kubernetes/kubernetes/issues/37306

---
