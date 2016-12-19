#! /bin/bash

requests_per_client=8000
replicas=3
urlraft=http://192.168.1.201:32100
urlapp=http://192.168.1.201:31000

# reset the cluster
# sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc raft --replicas=0
# sleep 10
# sudo /opt/bin/kubectl --server=192.168.1.200:8080 scale rc raft --replicas=$replicas

for clients in 4 8 16; do

	echo "Ready to run! requests=$requests clients=$clients replicas=$replicas url=$urlraft"
	echo "Press enter to continue"

	read -n1 -s

	ab -n $requests -c $clients -s 5 -e tests/tests_physical/raft_${clients}_${requests}_client.csv $urlraft/request

	echo "Done"

	read -n1 -s

	for i in $(seq 6); do
		curl $urlraft/hash
		echo
	done


	for i in $(seq 6); do
		curl $urlapp/hash
		echo
	done

	echo "Done"

done
