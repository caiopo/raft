package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var Log []string

func AddToLog(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	url := r.URL.Path[1:]

	fmt.Println(url)

	Log = append(Log, string(url))

}

func Print(w http.ResponseWriter, r *http.Request) {
	for i, e := range Log {
		io.WriteString(w, fmt.Sprintf("%d - %s\n", i, e))
	}

}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <port>", os.Args[0])
		os.Exit(1)
	}

	fmt.Printf("Port: %s\n", os.Args[1])

	http.HandleFunc("/print", Print)
	http.HandleFunc("/", AddToLog)
	http.ListenAndServe(":"+os.Args[1], nil)
}
