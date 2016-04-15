package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	raft "github.com/caiopo/pontoon"
)

const (
	kubernetesAPIServer = "192.168.1.200:8080"
	runningInKubernetes = false
)

var validPorts = []string{":55125", ":55126", ":55127", ":55128", ":55129", ":55130"}

func main() {

	if runningInKubernetes {

		myip := getMyIP("18")

		fmt.Println(myip)

		if myip == "badIPReturn" {
			myip = "localhost"
		}

		transport := &raft.HTTPTransport{Address: myip + raft.PORT}
		logger := &raft.Log{}
		applyer := &raft.StateMachine{}
		node := raft.NewNode(myip, transport, logger, applyer)
		node.Serve()

		node.Start()
		defer node.Exit()

		ipsAdded := make([]string, 0)

		for {
			ipsKube := getIPsFromKubernetes()

			fmt.Print("IPs Kube: ")
			fmt.Println(ipsKube)

			fmt.Print("IPs Added: ")
			fmt.Println(ipsAdded)

			for _, ipKube := range ipsKube {
				if !find(ipKube, ipsAdded) && (ipKube != myip) {
					node.AddToCluster(ipKube + raft.PORT)
					ipsAdded = append(ipsAdded, (ipKube))
				}
			}

			time.Sleep(time.Second)
		}

	} else {
		myip := getMyIP("")
		fmt.Println(myip)

		myport := ":" + os.Args[1]

		if !find(myport, validPorts) {
			panic("port must be between 55125 and 55130")
		}

		transport := &raft.HTTPTransport{Address: myip + myport}
		logger := &raft.Log{}
		applyer := &raft.StateMachine{}
		node := raft.NewNode(myip, transport, logger, applyer)
		node.Serve()

		node.Start()
		defer node.Exit()

		// cluster := make([]string, 0)

		cluster := os.Args[2:]

		// for _, ip := range =os.Args[2:] {
		// 	if ip != myip {
		// 		cluster = append(cluster, ip)
		// 	}
		// }

		fmt.Println(cluster)

		mutex := &sync.Mutex{}
		ipsAdded := make([]string, 0)

		for {

			for _, remoteip := range cluster {

				go func(nodeip string, mut *sync.Mutex) {
					// fmt.Println("node: " + nodeip)

					for _, remoteport := range validPorts {
						// fmt.Println("port: " + remoteport)

						if remoteport == myport && nodeip == myip {
							continue
						}

						if find(nodeip+remoteport, ipsAdded) {
							continue
						}

						go func(ip string, port string, m *sync.Mutex) {

							// fmt.Println("get: http://" + ip + port + "/ping")

							resp, err := http.Get("http://" + ip + port + "/ping")

							if err != nil {
								return
							}

							defer resp.Body.Close()

							body, err := ioutil.ReadAll(resp.Body)

							if err != nil {
								return
							}

							ss := string(body[:])

							fmt.Println(ss)

							if ss != "" {
								m.Lock()
								ipsAdded = append(ipsAdded, ip+port)
								node.AddToCluster(ip + port)
								m.Unlock()

							}

						}(nodeip, remoteport, mut)

					}

				}(remoteip, mutex)

			}

			fmt.Println(ipsAdded)
			time.Sleep(time.Second)

		}
	}

}

func find(needle string, haystack []string) bool {
	for _, h := range haystack {
		if needle == h {
			return true
		}
	}
	return false
}

func getIPsFromKubernetes() []string {
	resp, err := http.Get("http://" + kubernetesAPIServer + "/api/v1/endpoints")

	if err != nil {
		// raft.Debug += fmt.Sprintln("ERROR getting endpoints in kubernetes API: ", err.Error())
		return nil
	}
	defer resp.Body.Close()
	contentByte, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		// raft.Debug += fmt.Sprintln("ERROR reading data from endpoints: ", err2.Error())
		return nil
	}

	content := string(contentByte)

	replicas := make([]string, 0)

	words := strings.Split(content, "\"ip\":")
	for _, v := range words {
		if v[1:7] == "18.16." { //18.x.x.x, IP of PODS
			parts := strings.Split(v, ",")
			theIP := parts[0]
			theIP = theIP[1 : len(theIP)-1] //remove " chars from IP
			replicas = append(replicas, theIP)
		}
	}

	return replicas
}

func getMyIP(firstChars string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "couldNotConfigurateIP:" + err.Error()
	} else {
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
	}
	return "badIPReturn"
}

// func findLeader(node *raft.Node) (ip string) {

// 	ipchan := make(chan string)

// 	for _, n := range node.Cluster {
// 		go func(ip string) {
// 			resp, err := http.Get(ip + raft.PORT + "/node")
// 			if err == nil {
// 				defer resp.Body.Close()

// 				body, err := ioutil.ReadAll(resp.Body)
// 				if err == nil {
// 					ss := string(body[:])
// 					c := "1"
// 					b := string(ss[1])

// 					fmt.Println(ss)

// 					if b == c {
// 						ipchan <- ip
// 					}
// 				}
// 			}
// 		}(n.ID)
// 	}

// 	return <-ipchan
// }
