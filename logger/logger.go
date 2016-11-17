package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var Log []string

func AddToLog(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	url := r.URL.Path[1:]

	if url == "" {
		fmt.Fprintln(w, "Empty request!")
		return
	}

	fmt.Println(url)

	Log = append(Log, url)

	fmt.Fprintf(w, "Saved: %v\n", url)
}

func Print(w http.ResponseWriter, r *http.Request) {
	for i, e := range Log {
		fmt.Fprintf(w, "%d - %s\n", i, e)
	}
}

func Hash(w http.ResponseWriter, r *http.Request) {
	hash := sha256.Sum256([]byte(strings.Join(Log, "\n")))

	fmt.Fprintf(w, "Hash: %v", hex.EncodeToString(hash[:]))
}

func main() {
	port := "65432"

	if len(os.Args) >= 2 {
		port = os.Args[1]
	}

	fmt.Printf("Port: %s\n", port)

	http.HandleFunc("/print", Print)
	http.HandleFunc("/hash", Hash)
	http.HandleFunc("/", AddToLog)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}
