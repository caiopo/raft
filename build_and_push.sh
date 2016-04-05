#! /bin/bash

go build raft.go
docker build -t caiopo/raft .
docker push caiopo/raft