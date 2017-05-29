package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

var Log []string
var LogLock sync.RWMutex

func AddToLog(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	url := r.URL.Path[1:]

	if url == "" {
		fmt.Fprintln(w, "Empty request!")
		return
	}

	fmt.Println(url)

	LogLock.Lock()
	Log = append(Log, url)
	LogLock.Unlock()

	fmt.Fprintf(w, "Saved: %v\n", url)
}

func Print(w http.ResponseWriter, r *http.Request) {
	LogLock.RLock()
	for i, e := range Log {
		fmt.Fprintf(w, "%d - %s\n", i, e)
	}
	LogLock.RUnlock()
}

func Hash(w http.ResponseWriter, r *http.Request) {

	LogLock.RLock()
	joinLog := strings.Join(Log, "\n")
	LogLock.RUnlock()

	hash := sha256.Sum256([]byte(joinLog))

	fmt.Fprintf(w, "Hash: %v", hex.EncodeToString(hash[:]))
}

func Length(w http.ResponseWriter, r *http.Request) {

	LogLock.RLock()
	logLen := len(Log)
	LogLock.RUnlock()

	fmt.Fprintf(w, "%d", logLen)
}

func main() {
	port := "65432"

	if len(os.Args) >= 2 {
		port = os.Args[1]
	}

	fmt.Printf("Port: %s\n", port)

	http.HandleFunc("/print", Print)
	http.HandleFunc("/hash", Hash)
	http.HandleFunc("/len", Length)
	http.HandleFunc("/", AddToLog)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}
