package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

const (
	KubernetesAPIServer = "192.168.1.200:8080"
	TAG                 = "raft"
	PORT                = ":55123"
)

var (
	targetIP string
)

func UpdateIP() error {
	hst, err := getIPsFromKubernetes(TAG)

	if err != nil {
		return err
	}

	hosts := make([]string, len(hst))

	for i, ip := range hst {
		resp, err := http.Get("http://" + ip + PORT + "/leader")

		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		hosts[i] = strings.Trim(string(body), "\" ")
	}

	if len(hosts) != 0 {
		targetIP = hosts[0]
		return nil
	}

	return errors.New(":(")
}

func Update(rw http.ResponseWriter, req *http.Request) {
	err := UpdateIP()

	if err != nil {
		fmt.Fprintln(rw, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(rw, "Target: ", targetIP)
}

func Target(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "Target: ", targetIP)
}

func Forward(rw http.ResponseWriter, req *http.Request) {
	url := targetIP + req.URL.Path

	log.Printf("fowarding: %v", url)

	resp, err := http.Get("http://" + url)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, "Error: get: ", err.Error())
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, "Error: read body: ", err.Error())
		return
	}

	rw.WriteHeader(resp.StatusCode)
	rw.Write(body)
}

func Version(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, 1)
}

func main() {

	http.HandleFunc("/update", Update)
	http.HandleFunc("/target", Target)
	http.HandleFunc("/version", Version)
	http.HandleFunc("/", Forward)

	time.Sleep(time.Second)

	if UpdateIP() != nil {
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
