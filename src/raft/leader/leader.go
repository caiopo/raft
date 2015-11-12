package leader

import (
	. "fmt"
	"raft"
	"raft/settings"
	"time"
)

var (
	heartbeatMsg = raft.Message{settings.Port, " ", "heartbeat"}
)

func HandleRequest(msg string) {

}

func Heartbeat(ch chan string) {

	for _, target := range settings.Cluster {

		if target != settings.Port {
			heartbeatMsg.SendTo(target)
		}

	}

	Print("Sent heartbeats to cluster...")

	countResponse := 0

	timeout := time.After(time.Second)

Loop:
	for countResponse < len(settings.Cluster) {

		select {
		case heartbeatResponse := <-ch:

			aliveMsg := raft.DecomposeMessage(heartbeatResponse)

			if aliveMsg.State == "alive" {

				countResponse++

			}

		case <-timeout:
			break Loop
		}

	}

	Printf("%d heartbeats received")
}
