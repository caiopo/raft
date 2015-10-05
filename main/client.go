package main

import (
	. "fmt"
	"os"
	"strconv"
	"time"
	"udp"
)

var cluster []string = []string{"56001", "56002", "56003", "56004", "56005"}

func CheckError(err error) {
	if err != nil {
		Println("Error: ", err)
	}
}

func main() {

	index, _ := strconv.Atoi(os.Args[1])

	for {
		udp.Send("heartbeat", cluster[index])
		time.Sleep(time.Second)
	}

	// var msg string

	// Scanf("%s", &msg)

	// udp.Send(msg, cluster[index])

}
