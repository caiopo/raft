package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	kubernetesAPIServer = "192.168.15.150:8080"
)

func main() {

	fmt.Println(getIPsFromKubernetes())
}

func getIPsFromKubernetes() []string {
	resp, err := http.Get("http://" + kubernetesAPIServer + "/api/v1/endpoints")

	if err != nil {
		fmt.Println("ERROR getting endpoints in kubernetes API: ", err.Error())
		return nil
	}
	defer resp.Body.Close()
	contentByte, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println("ERROR reading data from endpoints: ", err2.Error())
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
