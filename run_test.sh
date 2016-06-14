#! /bin/bash

requests=$2
clients=$3
replicas=3
url=http://192.168.1.201:$1

echo "Preparing... requests=$requests clients=$clients replicas=$replicas"

# reset the cluster
sudo /opt/bin/kubectl --server=192.168.15.150:8080 scale rc raft --replicas=0
sleep 5
sudo /opt/bin/kubectl --server=192.168.15.150:8080 scale rc raft --replicas=$replicas

echo "Ready to run! requests=$requests clients=$clients replicas=$replicas"
echo "Press enter to continue"

read -n1 -s

ab -n $requests -c $clients -s 5 -e tests/tests_kube/raft_ku_$clients_$requests_client.csv $url/request

sleep 2

for i in seq 10; do
	curl $url/hash
done

echo "Done"