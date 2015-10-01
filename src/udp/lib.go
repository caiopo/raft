package udp

import (
	. "fmt"
	"net"
	"os"
	// "reflect"
)

func Foo() string {
	return "bar"
}

var (
	RecvPort *net.UDPAddr
	err      error
)

/* A Simple function to verify error */
func checkError(err error) {
	if err != nil {
		Println("Error: ", err)
		os.Exit(0)
	}
}

func SetRecvPort(port string) {

	RecvPort, err = net.ResolveUDPAddr("udp", ":"+port)
	checkError(err)

}

func Receive() (string, *net.UDPAddr) {

	ServerConn, err := net.ListenUDP("udp", RecvPort)

	checkError(err)

	defer ServerConn.Close()

	buf := make([]byte, 1024)

	n, addr, err := ServerConn.ReadFromUDP(buf)
	Println("Received ", string(buf[0:n]), " from ", addr)

	if err != nil {
		Println("Error: ", err)
	}

	return string(buf[0:n]), addr

}

func Send(target, msg string) {

}
