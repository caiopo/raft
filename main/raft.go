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

const (
	maxTimeout = time.Second * 3
	envVar     = "RAFT_PORT"
)

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

	var msg string

	for {

		if isLeader {

			listening := true

			ch := make(chan string, 10)

			for {

				for _, i := range settings.Cluster {
					if i != settings.Port {

						go udp.Send("heartbeat", i)

					}
				}

				Println("Sent heartbeats to cluster")

				for _, i := range settings.Cluster {

					countResponse := 0

					select {
					case msg = <-ch:
						if msg == "alive" {
							countResponse++
						}
					case <-time.After(time.Second):
						break
					}

				}

			}

			// time.Sleep(50 * time.Millisecond)

		} else {

			msg = udp.ReceiveTimeout()

			follower.HandleRequest(msg)

		}

	}
}
