package main

import (
	. "fmt"
	// "net"
	"os"
	"strconv"
	"time"
	"udp"
)

type TimeCounter bool

const maxTimeout = time.Second * 3

var cluster []string = []string{"56001", "56002", "56003", "56004", "56005"}

var leader bool = false

func init() {

	index, _ := strconv.Atoi(os.Args[1])

	udp.SetRecvPort(cluster[index])

	if index == 0 {
		leader = true
	}

	Printf("My port is %s\n", cluster[index])

}

func main() {

	for {
		if leader {

			for _, i := range cluster {
				go udp.Send("heartbeat", i)
				time.Sleep(50 * time.Millisecond)
			}

		} else {

			timeout := time.Now().Add(maxTimeout)
			msg := udp.ReceiveTimeout(timeout)

			if msg == "timeout" {

				Println(">> Starting election")
			} else {
				Println(msg)
			}
		}

	}
}
