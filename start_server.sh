#!/bin/bash

set -eux -o pipefail

# kill couchbase if it exists
docker kill couchbase || true
docker rm couchbase || true

docker run --rm -d --name couchbase -p 8091-8096:8091-8096 -p 11207:11207 -p 11210:11210 -p 11211:11211 -p 18091-18094:18091-18094 $DOCKER_IMAGE

docker exec couchbase curl -L --retry-all-errors --connect-timeout 5 --max-time 10 --retry 20 --retry-delay 0 --retry-max-time 1200 'http://127.0.0.1:8091/'

# Set up CBS
docker exec couchbase couchbase-cli node-init -c 127.0.0.1 --username Administrator --password password --ipv4
docker exec couchbase couchbase-cli cluster-init --cluster-username Administrator --cluster-password password --cluster-ramsize 3072 --cluster-index-ramsize 3072 --cluster-fts-ramsize 256 --services data,index,query
docker exec couchbase couchbase-cli setting-index --cluster couchbase://localhost --username Administrator --password password --index-threads 4 --index-log-level verbose --index-max-rollback-points 10 --index-storage-setting default --index-memory-snapshot-interval 150 --index-stable-snapshot-interval 40000

curl -u Administrator:password -v -X POST http://127.0.0.1:8091/node/controller/rename -d 'hostname=127.0.0.1'

echo ""
