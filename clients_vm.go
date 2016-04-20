package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const commandMessage = "Specify: n-clients n-requests n-replicas[just for logging] port ips"

const requestBody = "BODY"

var (
	wg sync.WaitGroup

	mutex  sync.Mutex
	file   *os.File
	writer *bufio.Writer

	nClients, nRequests, nReplicas, nIPs int

	targetIPs []string

	timeInit time.Time
)

func main() {

	var err error

	if len(os.Args) < 5 {
		fmt.Println(commandMessage)
		os.Exit(1)
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

	port := ":" + os.Args[4]

	targetIPs = os.Args[5:]
	nIPs = len(targetIPs)

	for i, ip := range targetIPs {
		targetIPs[i] = "http://" + ip + port
	}

	file, err = os.Create(fmt.Sprintf("raft_test_c%dreq%drep%d.csv", nClients, nRequests, nReplicas))

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

		ip := targetIPs[(clientID+r)%nIPs]

		target := fmt.Sprintf("%s/request/%d/%s", ip, requestID, requestBody)

		t0 = time.Now()

		var resp *http.Response

		err := errors.New("")

		resp, err = http.Get(target)

		if err != nil {
			go writeToFile(fmt.Sprintf("%d;%d;%d;%d;%d;%s;%d;%d;%d", 0, requestID, 0, 0, 0, requestBody, nClients, nRequests, nReplicas))
			fmt.Printf("c%d: %s\n", clientID, ip)
			continue
		}

		resp.Body.Close()

		var leader int

		if resp.StatusCode == 290 {
			leader = 0
		} else if resp.StatusCode == 291 {
			leader = 1
		} else {
			go writeToFile(fmt.Sprintf("%d;%d;%d;%d;%d;%s;%d;%d;%d", 0, requestID, 0, 0, leader, requestBody, nClients, nRequests, nReplicas))
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

	// fmt.Println(s)

	mutex.Unlock()

	if err != nil {
		os.Exit(1)
	}
}
