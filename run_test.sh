#! /bin/bash

requests=$1
clients=$2
replicas=3

echo "Preparing... requests=$requests clients=$clients replicas=$replicas"

# reset the cluster
kub scale rc raft --replicas=0
sleep 5
kub scale rc raft --replicas=$replicas

echo "Ready to run! requests=$requests clients=$clients replicas=$replicas"
echo "Press enter to continue"

read -n1 -s

ab -n $requests -c $clients -s 5 -e tests/tests_kube/raft_ku_${clients}_${requests}_client.csv http://192.168.1.201:55123/request

echo "Done"