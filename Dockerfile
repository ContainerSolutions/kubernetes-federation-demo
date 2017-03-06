#FROM golang:1.7.4-alpine3.5
FROM golang:wheezy

#RUN apk -q update && apk -q add ca-certificates && apk -q add git && mkdir -p /opt/geoserver
RUN apt-get update && apt-get install -y && apt-get install git && mkdir -p /opt/geoserver

ADD ./vendor/k8s.io /usr/local/go/src/k8s.io

ADD . /opt/geoserver/

#ADD ./kubernetes/client/bin/kubefed /opt/geoserver/kubefed

WORKDIR /opt/geoserver/

RUN go build

CMD ["/opt/geoserver/geoserver"]
