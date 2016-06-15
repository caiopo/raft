#! /bin/bash

requests=8000
# clients=
replicas=3
url=http://192.168.1.201:$1

# echo "Preparing... requests=$requests clients=$clients replicas=$replicas"

# reset the cluster
# sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc raft --replicas=0
# sleep 10
# sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc raft --replicas=$replicas


for cli in 4 8 16 32 64; do

	echo "Ready to run! requests=$requests clients=$clients replicas=$replicas url=$url"
	echo "Press enter to continue"

	read -n1 -s

	ab -n $requests -c $cli -s 5 -e tests/tests_physical/raft_ph_${cli}_${requests}_client.csv $url/request

	echo "Done"

	read -n1 -s

	for i in $(seq 6); do
		curl $url/hash
		echo
	done

	echo "Done"

done
