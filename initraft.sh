#! /bin/bash

kub run raft --image=caiopo/raft --replicas=0
kub expose rc raft --port=55123 --type=LoadBalancer
kub scale rc raft --replicas=$1