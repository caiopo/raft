package main

import (
	"net/http"
	"os"
	"strings"
)

func handler(rw http.ResponseWriter, req *http.Request) {
	path := strings.Split(req.URL.Path, "/")[1]

	if path == "" {

		for _, e := range os.Environ() {
			rw.Write([]byte(e + "\n"))
		}

	} else {
		rw.Write([]byte(os.Getenv(path)))
	}
}

func main() {

	http.HandleFunc("/", handler)

	http.ListenAndServe(":12345", nil)

}
