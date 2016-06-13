package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	raft "github.com/caiopo/pontoon"
)

func main() {
	if raft.RunningInKubernetes {
		log.SetOutput(ioutil.Discard)

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
		myip := os.Args[1]
		fmt.Println(myip)

		transport := &raft.HTTPTransport{Address: myip + raft.PORT}
		logger := &raft.Log{}
		applyer := &raft.StateMachine{}
		node := raft.NewNode(myip, transport, logger, applyer)
		node.Serve()

		node.Start()
		defer node.Exit()

		for _, ip := range os.Args[2:] {
			node.AddToCluster(ip + raft.PORT)
		}

		for {
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
	resp, err := http.Get("http://" + raft.KubernetesAPIServer + "/api/v1/endpoints")

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

// 	for {

// 		for _, ip := range cluster {

// 			go func(ip string) {

// 				for _, port := range raft.ValidPorts {

// 					if port == myport && ip == myip {
// 						continue
// 					}

// 					if find(ip+port, ipsAdded) {
// 						continue
// 					}

// 					go func(ip string, port string) {
// 						resp, err := http.Get("http://" + ip + port + "/ping")

// 						if err != nil {
// 							return
// 						}

// 						defer resp.Body.Close()

// 						body, err := ioutil.ReadAll(resp.Body)

// 						if err != nil {
// 							return
// 						}

// 						ss := string(body[:])

// 						fmt.Println(ss)

// 						if ss != "" {
// 							mutex.Lock()
// 							ipsAdded = append(ipsAdded, ip+port)
// 							node.AddToCluster(ip + port)
// 							mutex.Unlock()
// 						}

// 					}(ip, port)

// 				}

// 			}(ip)

// 		}

// 		fmt.Println(ipsAdded)
// 		time.Sleep(time.Second)

// 	}
// }
