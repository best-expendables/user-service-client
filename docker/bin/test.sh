#!/usr/bin/env bash
set -e

cd /go/src/user-service-go-client

go tool vet *.go
go test -redisAddr redis:6379 --cover
FILES=$(gofmt -l *.go)
echo $FILES
exit $([[ $FILES = "" ]] && echo 0 || echo 1)
