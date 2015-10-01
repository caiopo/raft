package main

import (
	. "fmt"
	"os"
	"strconv"
	"udp"
)

var cluster []string = []string{"56001", "56002", "56003", "56004", "56005"}

var leader bool = false

func init() {

	index, _ := strconv.Atoi(os.Args[1])

	udp.SetRecvPort(cluster[index])

	Printf("My port is %s\n", cluster[index])

}

func main() {

	Println(udp.Receive())

	if leader {

	} else {

	}

}
