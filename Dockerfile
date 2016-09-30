FROM golang:1.7-alpine

MAINTAINER Grigory Aksentyev <grigory.aksentiev@gmail.com>

RUN mkdir -p /go/src/github.com/aksentyev/postgres_exporter
COPY . /go/src/github.com/aksentyev/postgres_exporter
RUN apk add --no-cache git build-base

ENV GOROOT /usr/local/go
RUN cd /go/src/github.com/aksentyev/postgres_exporter \
    && go get -d \
    && go build -o /bin/postgres_exporter \
    && rm -rf /go/src/github.com/aksentyev/postgres_exporter

RUN apk del --purge git build-base

RUN mkdir /config
COPY exporter/queries.yaml /config/queries.yaml

WORKDIR /config

ENTRYPOINT ["/bin/postgres_exporter", "-log.level", "debug"]
CMD ["-update-interval", "3", "-scrape-interval", "3"]
