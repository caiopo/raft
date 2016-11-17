package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	sm   *StateMachine
)

func main() {
	if raft.RunningInKubernetes {
		logfile, err := os.OpenFile("raft.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

		if err != nil {
			fmt.Printf("error opening file: %v", err)
		}

		// don't forget to close it
		defer logfile.Close()

		log.SetOutput(logfile)

		myip, err = getMyIP("18")

		if err != nil {
			os.Exit(1)
		}

		fmt.Println(myip)

		transport := &raft.HTTPTransport{Address: myip + raft.PORT}
		logger := &raft.Log{}

		sm = &StateMachine{} // with application
		// applyer := &raft.StateMachine{} // without application

		node = raft.NewNode(myip, transport, logger, sm)
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
	log.Printf("starting getIPsFromKubernetes(%s)\n", tag)

	resp, err := http.Get("http://" + raft.KubernetesAPIServer + "/api/v1/namespaces/default/endpoints/" + tag)

	if err != nil {
		return nil, err
	}

	log.Println("GET successful")

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	log.Println("Read body")

	json, err := simplejson.NewJson(content)

	if err != nil {
		return nil, err
	}

	log.Println("Parsing json")

	replicas := make([]string, 0)

	addresses := json.Get("subsets").GetIndex(0).Get("addresses")

	for i := 0; ; i++ {

		ip, err := addresses.GetIndex(i).Get("ip").String()

		if err != nil {
			break
		}

		replicas = append(replicas, ip)
	}

	log.Printf("Done: %v\n", replicas)

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
	ipsAddedRaft := make([]string, 0)
	ipsAddedApp := make([]string, 0)

	appName := os.Getenv("RAFT_APP")

	if appName == "" {
		// TODO: write log to file
		log.Println("RAFT_APP environment variable not found")
		os.Exit(1)
	}

	for {
		ipsRaft, err := getIPsFromKubernetes("raft")

		if err != nil {
			continue
		}

		for _, ipRaft := range ipsRaft {
			if !find(ipRaft, ipsAddedRaft) && (ipRaft != myip) {
				node.AddToCluster(ipRaft + raft.PORT)
				ipsAddedRaft = append(ipsAddedRaft, ipRaft)
			}
		}

		log.Printf("IPs Raft: %v IPs Added: %v My IP:%v\n", ipsRaft, ipsAddedRaft, myip)

		time.Sleep(500 * time.Millisecond)

		ipsApp, err := getIPsFromKubernetes(appName)

		if err != nil {
			continue
		}

		for _, ipApp := range ipsApp {
			if !find(ipApp, ipsAddedApp) {
				sm.AddReplica(ipApp)
				ipsAddedApp = append(ipsAddedApp, ipApp)
			}
		}

		log.Printf("IPs App: %v IPs Added: %v\n", ipsApp, ipsAddedApp)

		time.Sleep(500 * time.Millisecond)
	}
}
