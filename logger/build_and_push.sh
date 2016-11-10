#! /bin/bash

go build -ldflags "-linkmode external -extldflags -static"
docker build -t caiopo/logger .
rm logger
docker push caiopo/logger
