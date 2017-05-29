package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/bitly/go-simplejson"
)

const (
	KubernetesAPIServer = "192.168.1.200:8080"
	Port                = ":55123"
	Tag                 = "raft"
)

var (
	hosts    []string
	next     int = 0
	nextLock sync.Mutex
)

func UpdateHosts() error {
	hst, err := getIPsFromKubernetes(Tag)

	hosts = hst

	return err
}

func Update(rw http.ResponseWriter, req *http.Request) {
	err := UpdateHosts()

	if err != nil {
		fmt.Fprintln(rw, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(rw, "Hosts: ", hosts)
}

func Hosts(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "Hosts: ", hosts)
}

func Forward(rw http.ResponseWriter, req *http.Request) {
	nextLock.Lock()
	host := hosts[next]
	next = (next + 1) % len(hosts)
	nextLock.Unlock()

	url := host + Port + req.URL.Path

	log.Printf("fowarding: %v", url)

	resp, err := http.Get("http://" + url)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(resp.StatusCode)
	rw.Write(body)
}

func main() {

	http.HandleFunc("/update", Update)
	http.HandleFunc("/hosts", Hosts)
	http.HandleFunc("/", Forward)

	if UpdateHosts() != nil {
		fmt.Println("Error on initial update!")
		os.Exit(1)
	}

	http.ListenAndServe(":12345", nil)

}

func getIPsFromKubernetes(tag string) ([]string, error) {
	log.Printf("starting getIPsFromKubernetes(%s)\n", tag)

	resp, err := http.Get("http://" + KubernetesAPIServer + "/api/v1/namespaces/default/endpoints/" + tag)

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
