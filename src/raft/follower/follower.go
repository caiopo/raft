package follower

import (
	. "fmt"
	"raft"
	"raft/settings"
	"udp"
)

var (
	history     []string
	leaderPort  string
	state       int
	promisedMsg raft.Message
)

const (
	IDLE    int = 0
	PROMISE int = 1
)

func init() {

	history = make([]string, 0)
	state = IDLE

}

func HandleRequest(initMsg string, ch chan string) {

	Println(history)

	if initMsg == udp.TIMEOUT {

		startElection()

	} else {

		decompMsg := raft.DecomposeMessage(initMsg)

		if decompMsg.State == "heartbeat" {

			aliveMsg := raft.Message{settings.Port, " ", "alive"}
			go aliveMsg.SendTo(decompMsg.Sender)

		} else {

			switch state {
			case IDLE:

				if decompMsg.State == "proposal" {

					promiseMsg := raft.Message{settings.Port, decompMsg.Content, "promise"}
					go promiseMsg.SendTo(decompMsg.Sender)

					state = PROMISE

				}

			case PROMISE:

				if decompMsg.State == "accepted" {

					finishedMsg := raft.Message{settings.Port, decompMsg.Content, "finished"}
					go finishedMsg.SendTo(decompMsg.Sender)

					history = append(history, decompMsg.Content)

					state = IDLE

				}

			}

		}

	}

}

func startElection() {
	Println(">> Starting election")

}

// func HandleRequest(initMsg string, ch chan string) {

// 	Println(history)

// 	if initMsg == "timeout" {

// 		Println(">> Starting election")

// 		startElection()

// 	} else {

// 		decompMsg := raft.DecomposeMessage(initMsg)

// 		if decompMsg.Content == "heartbeat" {

// 			aliveMsg := raft.Message{settings.Port, " ", "alive"}

// 			aliveMsg.SendTo(decompMsg.Sender)

// 			// udp.Send("alive", decompMsg.Sender)

// 		} else if decompMsg.State == "proposal" {

// 			promiseMsg := raft.Message{settings.Port, decompMsg.Content, "promise"}

// 			util.EmptyChan(ch)

// 			go promiseMsg.SendTo(decompMsg.Sender)

// 			tempMsg, _ := udp.Receive()

// 			finalMsg := raft.DecomposeMessage(tempMsg)

// 			if finalMsg.State == "accepted" {

// 				finishedMsg := raft.Message{settings.Port, decompMsg.Content, "finished"}

// 				finishedMsg.SendTo(decompMsg.Sender)

// 				// udp.Send("finished:"+decompMsg.Content, decompMsg.Sender)

// 				history = append(history, decompMsg.Content)

// 			}

// 		}

// 	}

// }
