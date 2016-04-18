#! /bin/bash

IP=$1

for r in 3 4 7 10 15; do
	for c in 4 8 16 32 64; do
		echo "Starting test. Clients: $c Replicas: $r"

		./initraft.sh

		PORT=$(kub get svc raft -o json | ./filter.py)

		kub get pods

		read -p "Waiting for pods to be ready" -n 1 -s

		./clients $c 100 $r $IP:$PORT

		read -p "Waiting for requests to complete" -n 1 -s

		read -n 1 -s

		./stopraft.sh

		read -p "Go to next test?" -n 1 -s

		read -n 1 -s
	done
done

