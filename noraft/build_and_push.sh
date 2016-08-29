#! /bin/bash

go build -ldflags "-linkmode external -extldflags -static" noraft.go
docker build -t caiopo/noraft .
docker push caiopo/noraft
rm noraft
