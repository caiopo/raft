#! /bin/bash

requests=$2
clients=$3
replicas=3
url=http://192.168.1.201:$1

# echo "Preparing... requests=$requests clients=$clients replicas=$replicas"

# reset the cluster
# sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc raft --replicas=0
# sleep 10
# sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc raft --replicas=$replicas

echo "Ready to run! requests=$requests clients=$clients replicas=$replicas"
echo "Press enter to continue"

read -n1 -s

ab -n $requests -c $clients -s 5 -e tests/tests_physical/raft_ph_${clients}_${requests}_client.csv $url/

echo "Done"

read -n1 -s

# for i in $(seq 6); do
# 	curl $url/hash
# 	echo
# done

echo "Done"