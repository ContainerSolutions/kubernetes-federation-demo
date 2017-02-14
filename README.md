# Kubernetes Federation Demo

## Requirements

### Kubernetes command line tools
Download the command line binaries located at: `https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG.md#client-binaries`
Install the binaries in your $PATH (Ex. /usr/local/sbin)

### GKE secret key
Generate a Google platform secret key, in JSON format and store it in this location: `$HOME/.federation-key.json`

### Preparing the demo

## Setup GKE clusters in different regions

*WARNING:* this creates several clusters in Google Cloud Platform! Watch out for your billing! /!\

### Build and push the docker image for the demo application
Adjust the version of the image inside the Dockerfile.

1. make push

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
But don't launch it in the same regions as the Kubernetes clusters, as this will not
work nicely with the Maps user interface.

The required arguments are the region where to generate traffic from, the IP address 
of the global ingress load balancer (from step 2) and the IP address of the admin service (FIXME where do you find it easily?).

    scripts/start-traffic.sh asia-east1 130.211.41.245 104.155.43.73

To remove the VM generating traffic again, run `scripts/stop-traffic.sh REGION`. This 
may also be useful to clean up when setting up traffic generation might have failed for some reason.

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

