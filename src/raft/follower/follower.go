package follower

import (
	. "fmt"
	// "raft/settings"
	"strings"
	"udp"
	"util"
)

var history []string

var leaderPort string

func init() {

	history = make([]string, 0)

}

func HandleRequest(initMsg string) {

	Println(history)

	if initMsg == "timeout" {

		Println(">> Starting election")

		startElection()

	} else {

		leaderPort, content, state := separate(initMsg)

		if content == "heartbeat" {

			udp.Send("alive", leaderPort)

		} else if state != "" {

			if state == "proposal" {

				udp.Send(content+":promise", leaderPort)

				acceptMsg, _ := udp.Receive()

				_, confirmation, _ := separate(acceptMsg)

				if confirmation == "accept" {

					udp.Send("finished:"+content, leaderPort)
					history = append(history, content)

				}

			}

		}

	}

}

func startElection() {

}

func separate(initMsg string) (target, content, state string) {

	sep := util.CreateSeparator(':')

	msg := strings.FieldsFunc(initMsg, sep)

	target = msg[0]

	content = msg[1]

	if len(msg) == 3 {

		state = msg[2]
	} else {
		state = ""
	}

	return
}
