#! /bin/bash

go build -ldflags "-linkmode external -extldflags -static"
docker build -t caiopo/raft:latest .
rm raft
docker push caiopo/raft:latest
