package main

import (
	. "fmt"
	"os"
	// "raft"
	"raft/follower"
	"raft/leader"
	"raft/settings"
	"strconv"
	// "strings"
	"time"
	"udp"
	// "util"
)

var (
	isLeader bool = false
)

const (
	maxTimeout = time.Second * 3
)

func init() {

	index, _ := strconv.Atoi(os.Args[1])

	settings.Port = settings.Cluster[index]

	udp.SetRecvPort(settings.Port)

	if index == 0 {
		isLeader = true
	}

	Printf("My port is %s\n", settings.Port)

}

func main() {

	ch := make(chan string, 10)

	// chAlive := true

	var msg string

	for {

		if isLeader {

			leader.Heartbeat(ch)

		} else {

			udp.SetTimeout(time.Now().Add(maxTimeout))

			msg = udp.ReceiveTimeout()

			follower.HandleRequest(msg, ch)

		}

	}
}
