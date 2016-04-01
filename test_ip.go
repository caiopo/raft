package main

import (
	"fmt"
	"http"
	"net"
)

func main() {
	// addrs, _ := net.InterfaceAddrs()

	// for _, addr := range addrs {
	// 	fmt.Println(addr)
	// }

	// fmt.Println(getMyIP(""))

	http.Get

}

func getMyIP(firstChars string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "couldNotConfigurateIP:" + err.Error()
	} else {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					s := ipnet.IP.String()
					if s[:len(firstChars)] == firstChars {
						return s
					}
				}
			}
		}
	}
	return "badIPReturn"
}
