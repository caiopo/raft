package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/bitly/go-simplejson"

	raft "github.com/caiopo/pontoon"
)

var (
	node *raft.Node
	myip string
)

func main() {
	if raft.RunningInKubernetes {
		// log.SetOutput(ioutil.Discard)

		var err error

		myip, err = getMyIP("18")

		if err != nil {
			os.Exit(1)
		}

		fmt.Println(myip)

		transport := &raft.HTTPTransport{Address: myip + raft.PORT}
		logger := &raft.Log{}
		applyer := &raft.StateMachine{}
		node = raft.NewNode(myip, transport, logger, applyer)
		node.Serve()

		node.Start()

		updateCluster()

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

func getIPsFromKubernetes(tag string) ([]string, error) {
	resp, err := http.Get("http://" + raft.KubernetesAPIServer + "/api/v1/namespaces/default/endpoints/" + tag)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	json, err := simplejson.NewJson(content)

	if err != nil {
		return nil, err
	}

	replicas := make([]string, 0)

	addresses := json.Get("subsets").GetIndex(0).Get("addresses")

	for i := 0; ; i++ {

		ip, err := addresses.GetIndex(i).Get("ip").String()

		if err != nil {
			break
		}

		replicas = append(replicas, ip)
	}

	return replicas, nil
}

func getMyIP(firstChars string) (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				s := ipnet.IP.String()
				if s[:len(firstChars)] == firstChars {
					return s, nil
				}
			}
		}
	}

	return "", errors.New("ip not found")
}

func updateCluster() {
	ipsAdded := make([]string, 0)

	for {
		ipsKube, err := getIPsFromKubernetes("raft")

		if err != nil {
			continue
		}

		fmt.Printf("IPs Kube: %v\nIPs Added: %v\nMy IP:%v\n\n", ipsKube, ipsAdded, myip)

		for _, ipKube := range ipsKube {
			if !find(ipKube, ipsAdded) && (ipKube != myip) {
				node.AddToCluster(ipKube + raft.PORT)
				ipsAdded = append(ipsAdded, ipKube)
			}
		}

		time.Sleep(time.Second)
	}
}
