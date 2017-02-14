IMAGE=containersoluploader/geoserver
VERSION=0.4.8
FED_PROJECT=steam-ego-156812

build:
	go build

docker: build
	docker build -t ${IMAGE}:${VERSION} .

push: docker
	docker push ${IMAGE}:${VERSION}

europe:	
	gcloud container clusters get-credentials gce-europe-west1-b --zone europe-west1-b --project ${FED_PROJECT}

deploy-admin: europe        
	kubectl create -f manifests/geoserver-admin.yaml

destroy-admin: europe
	kubectl delete -f manifests/geoserver-admin.yaml

deploy-replica:
	kubectl --context=federation create -f manifests/geoserver-replica.yaml

deploy-service:
	kubectl --context=federation create -f manifests/geoserver-service.yaml

deploy-ingress:
	kubectl --context=federation create -f manifests/geoserver-ingress.yaml

status:
	kubectl get clusters
	# kubectl --context=federation describe ingress geoserver-ingress
	# kubectl --context=federation describe service geoserver-service
	# kubectl --context=federation describe replicaset geoserver-replica

destroy:	
	kubectl --context=federation delete -f manifests/geoserver-ingress.yaml
	kubectl --context=federation delete -f manifests/geoserver-service.yaml
	kubectl --context=federation delete -f manifests/geoserver-replica.yaml	
	kubectl delete ns federation-system
	

all: build
