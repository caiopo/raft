package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const commandMessage = "Specify: n-clients n-requests n-replicas[just for logging] ip:port (optional: path)"

const requestBody = "BODY"

var (
	wg sync.WaitGroup

	mutex  sync.Mutex
	file   *os.File
	writer *bufio.Writer

	nClients, nRequests int

	targetIP, path string
)

func main() {

	var err error

	if len(os.Args) < 5 {
		fmt.Println(commandMessage)
		os.Exit(1)
	}

	if len(os.Args == 6) {
		path = os.Args[5]
	}

	nClients, err = strconv.Atoi(os.Args[1])

	if err != nil {
		fmt.Println(commandMessage, ": ", err.Error())
		os.Exit(1)
	}

	nRequests, err = strconv.Atoi(os.Args[2])

	if err != nil {
		fmt.Println(commandMessage, ": ", err.Error())
		os.Exit(1)
	}

	nReplicas := os.Args[3]

	targetIP = "http://" + os.Args[4]

	file, err = os.Create(fmt.Sprintf(path+"raft_test_c%dreq%drep%s.txt", nClients, nRequests, nReplicas))

	if err != nil {
		fmt.Println("Can't create file")
		os.Exit(1)
	}

	writer = bufio.NewWriter(file)

	wg.Add(nClients)

	for c := 1; c <= nClients; c++ {
		go client(c)
	}

	wg.Wait()
	file.Sync()
}

func client(clientID int) {
	var t0, t1 time.Time

	for r := 0; r < nRequests; r++ {

		requestID := 1000*clientID + r

		target := fmt.Sprintf("%s/command/%d/%s", targetIP, requestID, requestBody)

		t0 = time.Now()

		resp, err := http.Get(target)

		if err != nil {
			go writeToFile(fmt.Sprintf("Error on HTTP/GET! Client: %d, Request %d", clientID, requestID))
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode == 299 || resp.StatusCode == 200 { // request accepted
			t1 = time.Now()

			diff := t1.Sub(t0).Nanoseconds()

			// client;request;time(ns);time(ms);requestBody
			go writeToFile(fmt.Sprintf("%d;%d;%d;%d;%s", clientID, requestID, diff, int64(diff/1000000), requestBody))

		} else {
			go writeToFile(fmt.Sprintf("Error on command! Status code: %d Client: %d Request %d", resp.StatusCode, clientID, requestID))
		}

	}

	wg.Done()
}

func writeToFile(s string) {
	mutex.Lock()

	_, err := writer.WriteString(s + "\n")

	writer.Flush()

	fmt.Println(s)

	mutex.Unlock()

	if err != nil {
		os.Exit(1)
	}

}
