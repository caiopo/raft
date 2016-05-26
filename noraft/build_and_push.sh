#! /bin/bash

go build noraft.go
docker build -t caiopo/noraft .
docker push caiopo/noraft
rm noraft