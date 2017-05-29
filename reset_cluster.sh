#! /bin/bash

sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc fw1 --replicas=0
sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc raft --replicas=0
sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc logger --replicas=0

sleep 10

sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc logger --replicas=$1
sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc raft --replicas=3
sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc fw1 --replicas=1
