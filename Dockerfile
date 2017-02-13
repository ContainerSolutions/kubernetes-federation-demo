FROM golang:1.7.4-alpine3.5

RUN apk -q update && apk -q add ca-certificates && apk -q add git && mkdir -p /opt/geoserver

ADD . /opt/geoserver/

WORKDIR /opt/geoserver/

RUN go build

CMD ["/opt/geoserver/geoserver"]
