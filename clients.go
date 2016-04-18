package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

const commandMessage = "Specify: n-clients n-requests ip:port n-replicas[this, just for logging]"

const requestBody = "BODY"

var (
	wg sync.WaitGroup

	mutex sync.Mutex

	startTime time.Time

	file *os.File

	nClients, nRequests, nReplicas int

	targetIP string
)

func main() {

	if len(os.Args) < 5 {
		fmt.Println(commandMessage)
		os.Exit(1)
	}

	nClients, err := strconv.Atoi(os.Args[1])

	if err != nil {
		fmt.Println(commandMessage, ": ", err.Error())
		os.Exit(1)
	}

	nRequests, err := strconv.Atoi(os.Args[2])

	if err != nil {
		fmt.Println(commandMessage, ": ", err.Error())
		os.Exit(1)
	}

	targetIP := "http://" + os.Args[3]

	nReplicas := os.Args[4]

	file, err = os.Create(fmt.Sprintf("raft_test_c%dreq%drep%d.txt", nClients, nRequests, nReplicas))

	wg.Add(nClients)

	for c := 0; c < nClients; c++ {
		go client(c)
	}

	wg.Wait()
}

func client(clientID int) {
	var t0, t1 time.Time

	for r := 0; r < nRequests; r++ {

		target := makeIP(targetIP, clientID, r) + requestBody

		t0 = time.Now()

		resp, err := http.Get(target)

		if err != nil {
			go writeToFile(fmt.Sprintf("Error on HTTP/GET! Client: %d, Request %d", clientID, r))
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode == 299 {
			t1 = time.Now()

			diff := t1.Sub(t0).Nanoseconds()

			// client;request;time;requestBody
			go writeToFile(fmt.Printf("%d;%d;%d;%s", clientID, r, diff, requestBody))

		} else {
			go writeToFile(fmt.Sprintf("Error on command! Status code: %d Client: %d Request %d", resp.StatusCode, clientID, r))
		}

	}

	wg.Done()
}

func writeToFile(s string) {
	mutex.Lock()

	_, err := file.WriteString(s + "\n")

	fmt.Println(s)

	mutex.Unlock()

	if err != nil {
		os.Exit(1)
	}

}

func makeIP(baseIP string, clientID int, requestID int) string {
	return fmt.Sprintf("%s/command/%d/", baseIP, 1000*clientID+requestID)
}
