# Kubernetes Federation Demo

* [Requirements](#requirements)
    * [Install the latest Kubernetes command line tools](#install-the-latest-kubernetes-command-line-tools)
    * [Create a fresh GCP project](#create-a-fresh-gcp-project)
    * [Generate a key for the GCP service account](#generate-a-key-for-the-gcp-service-account)
* [Preparing the demo](#preparing-the-demo)
    * [Create the clusters and deploy the demo application](#create-the-clusters-and-deploy-the-demo-application)
    * [Explore the demo](#explore-the-demo)
    * [Generate traffic](#generate-traffic)
* [Demo time](#demo-time)
* [Simulating the demo](#simulating-the-demo)
* [Modifying the demo](#modifying-the-demo)


## Requirements

### Install the latest Kubernetes command line tools

Download the command line tools from [Github](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG.md#client-binaries) and 
install them in your `$PATH` (e.g. `/usr/local/sbin`).

### Create a fresh GCP project

Create a new project in your [Google Cloud Console](https://console.cloud.google.com/) and note down the project name.

### Generate a key for the GCP service account

Generate a Google Cloud Platform secret key in JSON format and store it in this location: `$HOME/.federation-key.json`.
You can generate it under *IAM & Admin* / *Service accounts*. Select *Create key* from the menu next to the default service account.

## Preparing the demo

**WARNING:** this creates several clusters in Google Cloud Platform! Watch out for your billing!

### Create the clusters and deploy the demo application

1. make sure `kubefed` and `kubectl` are in `$PATH`

2. Adjust some variables:

    Change `FED_PROJECT` and `FED_DNS_ZONE` in the `Makefile` and in `scripts/settings.source` to match your GCE project name 
    and some domain name that you can control (although actually changing DNS configuration it is not required for this demo).

2. Reserve two global IP adddress:

    ```
    gcloud compute addresses create kubernetes-ingress --global
    gcloud compute addresses create geoserver-admin --global
    ```
    Note down those IPs, they are needed later.

2. Set up a firewall rule:

    Assuming you have clusters in `us-east1-b`, `europe-west1-b` and `asia-east1-a`, run:

    ```
    gcloud compute firewall-rules create my-federated-ingress-firewall-rule --source-ranges 130.211.0.0/22 --allow "icmp,tcp:80,tcp:443 tcp:30000-33000" --target-tags "cluster-europe-west1-b,cluster-asia-east1-a,cluster-us-east1-b" --network default
    ```

    See https://github.com/kubernetes/kubernetes/issues/37306 for the reason behind this step.

3. Create the clusters:

    Change folder to `scripts`

    ```
    ./init.sh
    ```
    *This operation make take a long time.*

4. Update the local kubeconfig file    

    ```
    ./credentials.sh
    ```   

5. Initialise the federation:    

    ```
    ./init-federation.sh
    ```   

6. Join the clusters together:

    ```
    ./join.sh
    ```   

7. Setup the variable to the `geoserver-admin` manifest for handling the federation via API calls.

    ```
    ./clusters
    ```

8. Deploy Federated Service and Ingress and the federated application

    ```
    kubectl --context=federation create -f manifests/geoserver-service.yaml
    kubectl --context=federation create -f manifests/geoserver-ingress.yaml
    kubectl --context=federation create -f manifests/geoserver-replica.yaml
    ```

9. Deploy the demo application:

    ```
    kubectl --context=federation create -f manifests/geoserver-admin.yaml    
    ```

### Explore the demo

Point your browser to the IP address you generated in step 2 (geoserver-admin).
You can always find it again later by executing:

    gcloud compute addresses list | grep geoserver-admin

You should see the clusters appearing on a map, but no traffic yet.

### Generate traffic

The script `scripts/start-traffic.sh` will create a micro-instance in the desired Google 
datacenter and execute a small script that generates traffic on the federated cluster.
Launch it as often as you want with different regions as argument.

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

## Simulating the demo

If you don't want to or can't create a cluster, or if you want to hack on the map frontend, you can simulate the demo.
Just build and run the binary in `ADMIN` mode, then open your browser and add `?simulate` to the URL:

    make build
    ADMIN=1 ./kubernetes-federation-demo
    open http://localhost:8080/?simulate

You can modify the demo by changing the canned JSON responses in `static/simulate.js`.

## Modifying the demo

1. Make sure you have a DockerHub or GCR account

    If you haven't got a DockerHub account yet, create one at [hub.docker.com](https://hub.docker.com/).
    Alternatively, you should be equally successful using the [Google Container Registry](https://cloud.google.com/container-registry/).

2. Adjust the IMAGE variable

    Change the `IMAGE` variable in the `Makefile` to match your own image registry account.

3. Make your changes to the code, then build and push

    `make push`

---
