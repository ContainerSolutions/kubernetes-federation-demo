# *** CHANGE THIS *** (Google Cloud project name)
FED_PROJECT=steam-ego-156812

# Image name in version
IMAGE=containersoluploader/geoserver
VERSION=0.9.0

build:
	go build

docker: build
	docker build -t ${IMAGE}:${VERSION} .

push: docker
	docker push ${IMAGE}:${VERSION}

europe:	
	gcloud container clusters get-credentials gce-europe-west1-b --zone europe-west1-b --project ${FED_PROJECT}

deploy-admin: europe        
	kubectl --namespace=federation-system create -f manifests/geoserver-admin.yaml	

destroy-admin: europe
	kubectl --namespace=federation-system delete -f manifests/geoserver-admin.yaml

deploy-replica:
	kubectl --context=federation create -f manifests/geoserver-replica.yaml

destroy-replica:
	kubectl --context=federation delete -f manifests/geoserver-replica.yaml

deploy-service:
	kubectl --context=federation create -f manifests/geoserver-service.yaml

deploy-ingress:
	kubectl --context=federation create -f manifests/geoserver-ingress.yaml

clusters:
	kubectl get clusters	

destroy:	
	kubectl --context=federation delete -f manifests/geoserver-ingress.yaml
	kubectl --context=federation delete -f manifests/geoserver-service.yaml
	kubectl --context=federation delete -f manifests/geoserver-replica.yaml
	kubectl --context=federation delete ns federation-system
	gcloud container clusters get-credentials gce-europe-west1-b --zone europe-west1-b --project ${FED_PROJECT}
	kubectl delete -f manifests/geoserver-admin.yaml
	

all: build
