FROM golang:1.8

COPY ./ /go/src/user-service-go-client
WORKDIR /go/src/user-service-go-client

RUN cp -R /go/src/user-service-go-client/docker/bin /usr/local/bin/app \
    && go get -v \
    && chmod +x /usr/local/bin/app/*
