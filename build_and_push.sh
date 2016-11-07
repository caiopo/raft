#! /bin/bash

go build -ldflags "-linkmode external -extldflags -static"
docker build -t caiopo/raft .
rm raft
docker push caiopo/raft
