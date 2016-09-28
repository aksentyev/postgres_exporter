FROM golang:1.7-onbuild

MAINTAINER Grigory Aksentyev <grigory.aksentiev@gmail.com>

COPY . /go/src/app

RUN mkdir /config

COPY exporter/queries.yaml /config/queries.yaml

WORKDIR /config

ENTRYPOINT ["/go/bin/app", "-log.level", "debug"]
CMD ["-update-interval", "3", "-scrape-interval", "3"]
