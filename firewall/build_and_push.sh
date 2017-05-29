#! /bin/bash

go build -ldflags "-linkmode external -extldflags -static"
docker build -t caiopo/firewall .
rm firewall
docker push caiopo/firewall
