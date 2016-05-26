package main

import (
	"io"
	"net"
	"net/http"
)

func make_response() func(w http.ResponseWriter, r *http.Request) {
	response := "My IP is "+getMyIP("18")

	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, response)
	}
}

func main() {
	http.HandleFunc("/", make_response())
	http.ListenAndServe(":8000", nil)
}

func getMyIP(firstChars string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "couldNotConfigurateIP:" + err.Error()
	}

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

	return "badIPReturn"
}
