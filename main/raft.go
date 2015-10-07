package main

import (
	. "fmt"
	"os"
	"raft/follower"
	// "raft/leader"
	"raft/settings"
	"strconv"
	"time"
	"udp"
)

const maxTimeout = time.Second * 3

var isLeader bool = false

func init() {

	index, _ := strconv.Atoi(os.Args[1])

	settings.SetMyPort(settings.Cluster[index])

	udp.SetRecvPort(settings.Port)

	udp.SetTimeout(time.Now().Add(maxTimeout))

	if index == 0 {
		isLeader = true
	}

	Printf("My port is %s\n", settings.Port)

}

func main() {

	for {
		if isLeader {

			for _, i := range settings.Cluster {
				if i != settings.Port {
					go udp.Send("heartbeat", i)

				}
			}

			Println("Sent heartbeats to cluster")

			time.Sleep(50 * time.Millisecond)

		} else {

			// timeout :=
			msg := udp.ReceiveTimeout()

			follower.HandleRequest(msg)

		}

	}
}
