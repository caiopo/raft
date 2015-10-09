package main

import (
	. "fmt"
	"time"
	"udp"
)

func CheckError(err error) {
	if err != nil {
		Println("Error: ", err)
	}
}

func main() {

	// content := []string{"heartbeat", "heartbeat", "teste1:proposal", "accept", "teste2:proposal", "accept"}

	// udp.SetRecvPort("1234")

	// index, _ := strconv.Atoi(os.Args[1])

	// for _, msg := range content {
	// 	udp.Send(msg, settings.Cluster[index])
	// 	time.Sleep(time.Millisecond * 500)
	// }

	var msg string

	for {

		Scanf("%s", &msg)

		udp.Send(msg, "56000")
		time.Sleep(time.Second)
	}

}
