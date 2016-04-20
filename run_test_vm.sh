#! /bin/bash

ip="192.168.1.201 192.168.1.201 192.168.1.202"
port=55123
req=50
rep=3


for cli in 4 8 16 32; do
	read -n 1 -p "Please reset the cluster"

	echo "Starting test. Clients: $cli Replicas: $rep"

	./clients_vm $cli $req $rep $port $(echo $ip)	
done
