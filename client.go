package main

import (
	. "fmt"
	"net"
	"time"
)

// var cluster []string = []string{"56001", "56002", "56003", "56004", "56005"}

func CheckError(err error) {
	if err != nil {
		Println("Error: ", err)
	}
}

func main() {
	ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:56001")
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	defer Conn.Close()

	var msg string

	for {
		Scanf("%s", &msg)

		buf := []byte(msg)

		_, err = Conn.Write(buf)

		if err != nil {
			Println(msg, err)
		}

		Printf("Sent: %s to %s\n", msg, ServerAddr)

		time.Sleep(time.Second * 1)
	}
}
