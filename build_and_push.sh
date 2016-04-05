#! /bin/bash

go build test_raft.go
docker build -t caiopo/raft .
docker push caiopo/raft