#! /bin/bash

while [ 1 ]; do
	top -b -d 0.1|head -n12|tail -n10
	free -h|grep -i Mem
	ifstat 0.1 1
done