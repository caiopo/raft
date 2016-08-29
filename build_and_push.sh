#! /bin/bash

go build -ldflags "-linkmode external -extldflags -static" raft.go
docker build -t caiopo/raft .
docker push caiopo/raft
rm raft
