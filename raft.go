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

var (
	node *raft.Node
	myip string
)

func main() {
	if raft.RunningInKubernetes {
		log.SetOutput(ioutil.Discard)

		myip = getMyIP("18")

		fmt.Println(myip)

		if myip == "badIPReturn" {
			myip = "localhost"
		}

		transport := &raft.HTTPTransport{Address: myip + raft.PORT}
		logger := &raft.Log{}
		applyer := &StateMachine{}
		node = raft.NewNode(myip, transport, logger, applyer)
		node.Serve()

		node.Start()
		defer node.Exit()

	} else {
		myip = os.Args[1]
		fmt.Println(myip)

		transport := &raft.HTTPTransport{Address: myip}
		logger := &raft.Log{}
		applyer := &StateMachine{}

		applyer.AddReplica("localhost:65431")
		applyer.AddReplica("localhost:65432")

		node := raft.NewNode(myip, transport, logger, applyer)
		node.Serve()

		node.Start()
		defer node.Exit()

		for _, ip := range os.Args[2:] {
			node.AddToCluster(ip)
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
		return nil
	}

	defer resp.Body.Close()

	contentByte, err := ioutil.ReadAll(resp.Body)

	if err != nil {
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

func updateCluster() {
	ipsAdded := make([]string, 0)

	for {
		ipsKube := getIPsFromKubernetes()

		fmt.Print("IPs Kube: " + fmt.Sprintf("%v", ipsKube) + "\nIPs Added: " + fmt.Sprintf("%v", ipsKube))

		for _, ipKube := range ipsKube {
			if !find(ipKube, ipsAdded) && (ipKube != myip) {
				node.AddToCluster(ipKube + raft.PORT)
				ipsAdded = append(ipsAdded, (ipKube))
			}
		}

		time.Sleep(time.Second)
	}
}
