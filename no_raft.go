package main

import (
	"io"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func main() {
	http.HandleFunc("/", hello)
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
