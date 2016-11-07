package main

import (
	"fmt"
	"os"

	"github.com/bitly/go-simplejson"
)

const jsonstr = `{
  "kind": "Endpoints",
  "apiVersion": "v1",
  "metadata": {
    "name": "raft",
    "namespace": "default",
    "selfLink": "/api/v1/namespaces/default/endpoints/raft",
    "uid": "91d166b3-a501-11e6-857c-94de802df35a",
    "resourceVersion": "2612943",
    "creationTimestamp": "2016-11-07T15:48:04Z",
    "labels": {
      "run": "raft"
    }
  },
  "subsets": [
    {
      "addresses": [
        {
          "ip": "18.16.82.2",
          "targetRef": {
            "kind": "Pod",
            "namespace": "default",
            "name": "raft-pet5e",
            "uid": "86b10116-a509-11e6-857c-94de802df35a",
            "resourceVersion": "2612936"
          }
        },
        {
          "ip": "18.16.99.2",
          "targetRef": {
            "kind": "Pod",
            "namespace": "default",
            "name": "raft-0aqsp",
            "uid": "86b1127f-a509-11e6-857c-94de802df35a",
            "resourceVersion": "2612933"
          }
        }
      ],
      "notReadyAddresses": [
        {
          "ip": "18.16.39.2",
          "targetRef": {
            "kind": "Pod",
            "namespace": "default",
            "name": "raft-2foa5",
            "uid": "86b11954-a509-11e6-857c-94de802df35a",
            "resourceVersion": "2612942"
          }
        }
      ],
      "ports": [
        {
          "port": 55123,
          "protocol": "TCP"
        }
      ]
    }
  ]
}`

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
	}
}

func main() {
	json, err := simplejson.NewJson([]byte(jsonstr))

	check(err)

	// addresses, err := json.Get("subsets").Array()

	addresses := json.Get("subsets").GetIndex(0).Get("addresses")

	for i := 0; ; i++ {

		ip, err := addresses.GetIndex(i).Get("ip").String()

		if err != nil {
			break
		}

		fmt.Println(ip)

	}

	// fmt.Println(addresses[0].(map[string]interface{})["addresses"].([]interface{})[0].)

	// for _, addr := range addresses {
	// 	fmt.Println(addr["ip"])
	// 	replicas = append(replicas, addr["ip"].(string))
	// }
}
