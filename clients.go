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

	nClients, nRequests, nReplicas int

	targetIP, path string

	timeInit time.Time
)

func main() {

	var err error

	if len(os.Args) < 5 {
		fmt.Println(commandMessage)
		os.Exit(1)
	}

	if len(os.Args) == 6 {
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

	nReplicas, err = strconv.Atoi(os.Args[3])

	if err != nil {
		fmt.Println(commandMessage, ": ", err.Error())
		os.Exit(1)
	}

	targetIP = "http://" + os.Args[4]

	file, err = os.Create(fmt.Sprintf(path+"raft_test_c%dreq%drep%d.csv", nClients, nRequests, nReplicas))

	if err != nil {
		fmt.Println("Can't create file")
		os.Exit(1)
	}

	writer = bufio.NewWriter(file)

	wg.Add(nClients)

	timeInit = time.Now()

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

		target := fmt.Sprintf("%s/request/%d/%s", targetIP, requestID, requestBody)

		t0 = time.Now()

		resp, err := http.Get(target)

		if err != nil {
			go writeToFile(fmt.Sprintf("Error on HTTP/GET! Client: %d, Request %d", clientID, requestID))
			continue
		}

		defer resp.Body.Close()

		var leader int

		if resp.StatusCode == 290 {
			leader = 0
		} else if resp.StatusCode == 291 {
			leader = 1
		} else {
			go writeToFile(fmt.Sprintf("Error on command! Status code: %d Client: %d Request %d", resp.StatusCode, clientID, requestID))
			return
		}

		t1 = time.Now()

		diff := t1.Sub(t0).Nanoseconds()

		elapsed := time.Now().Sub(timeInit).Nanoseconds()

		// client;request;time(ns);time since start;requestBody;total clients;requests;replicas
		go writeToFile(fmt.Sprintf("%d;%d;%d;%d;%d;%s;%d;%d;%d", clientID, requestID, diff, elapsed, leader, requestBody, nClients, nRequests, nReplicas))

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
