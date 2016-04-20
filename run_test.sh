#! /bin/bash

#alias kub="sudo /opt/bin/kubectl --server=192.168.15.150:8080"

for r in 50 100; do
	for c in 4 8 16 32; do
		
		echo "Starting the cluster"

		sudo /opt/bin/kubectl --server=192.168.15.150:8080 scale rc raft --replicas=0
		sudo /opt/bin/kubectl --server=192.168.15.150:8080 scale rc raft --replicas=$1

		read -n 1 -p "Check if the cluster is ready"

		echo "Starting test. Clients: $c Replicas: $r"

		./clients $c $r $1 192.168.15.103:32480 tests/rep$1 > /dev/null

		for i in $(seq 1 $1); do
			curl 192.168.15.103:32480/hash
			echo ""
		done
	done
done

