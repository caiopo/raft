#! /bin/bash

ip=$1
path=$2
rep=$3

req=50


for cli in 4 8 16 32; do

	echo "Starting the cluster"

	sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc raft --replicas=0
	sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc raft --replicas=$rep

	read -n 1 -p "Check if the cluster is ready"

	echo "Starting test. Clients: $cli Replicas: $rep"

	./clients $cli $req $rep $ip $path > /dev/null

	for i in $(seq 1 $rep); do
		curl $ip/hash
		echo ""
	done
done
