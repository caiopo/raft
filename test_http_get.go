package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	resp, err := http.Get("http://127.0.0.1:5513/node")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)

	str := string(content)

	fmt.Println(str)

	if str == "\"1\"" {
		fmt.Println("yay")
	}
}
