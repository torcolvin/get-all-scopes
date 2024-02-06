#!/bin/bash
#
set -eux -o pipefail

DOCKER_IMAGE=build-docker.couchbase.com:8020/cb-vanilla/server:7.6.1-3080
#DOCKER_IMAGE=build-docker.couchbase.com:8020/cb-vanilla/server:7.6.0 # works

env DOCKER_IMAGE=${DOCKER_IMAGE} ./start_server.sh

#docker exec couchbase couchbase-cli bucket-create --cluster localhost --username Administrator --password password --bucket bucket1 --bucket-type couchbase --bucket-ramsize 200 --enable-flush 1

cd gocbonly
go run main.go
