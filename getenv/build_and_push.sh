#! /bin/bash

go build -ldflags "-linkmode external -extldflags -static" getenv.go
docker build -t caiopo/getenv .
docker push caiopo/getenv
rm getenv
