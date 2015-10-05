package udp

import (
	. "fmt"
	"net"
	"os"
	"time"
)

func Foo() string {
	return "bar"
}

var (
	RecvPort *net.UDPAddr
	err      error
)

// A Simple function to verify error
func checkError(err error) {
	if err != nil {
		Println("Error: ", err)
		os.Exit(0)
	}
}

// Set the port used by Receive()
func SetRecvPort(port string) {

	// Println(RecvPort)

	RecvPort, err = net.ResolveUDPAddr("udp", ":"+port)
	checkError(err)

}

func Receive() (string, *net.UDPAddr) {

	if RecvPort == nil {
		panic("Receive port not set")
	}

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

func ReceiveCh(ch chan<- string) {

	if RecvPort == nil {
		panic("Receive port not set")
	}

	ServerConn, err := net.ListenUDP("udp", RecvPort)
	checkError(err)

	defer ServerConn.Close()

	buf := make([]byte, 1024)

	n, addr, err := ServerConn.ReadFromUDP(buf)
	Println("Received ", string(buf[0:n]), " from ", addr)

	if err != nil {
		Println("Error: ", err)
	}

	ch <- string(buf[0:n])

}

func ReceiveTimeout(timeout time.Time) string {

	if RecvPort == nil {
		panic("Receive port not set")
	}

	ServerConn, err := net.ListenUDP("udp", RecvPort)
	checkError(err)

	ServerConn.SetDeadline(timeout)

	defer ServerConn.Close()

	buf := make([]byte, 1024)

	n, _, err := ServerConn.ReadFromUDP(buf)
	// Println("Received ", string(buf[0:n]), " from ", addr)

	if err != nil {
		// Println("Error: ", err)
		return "timeout"
	}

	return string(buf[0:n])

}

func Send(msg, target string) {
	TargetAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+target)
	checkError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	checkError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, TargetAddr)
	checkError(err)

	defer Conn.Close()

	buf := []byte(msg)

	_, err = Conn.Write(buf)

	if err != nil {
		Println(msg, err)
	}

	Printf("Sent: %s to %s\n", msg, TargetAddr)

}
