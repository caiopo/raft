package main

import (
	. "fmt"
	"time"
)

var i int

func main() {

	i = 0
	// for j := 0; j < 10; j++ {
	for {
		chString := make(chan string)

		alive := true

		go foo(chString, &alive)

		time.Sleep(time.Second)

		alive = false

		time.Sleep(50 * time.Millisecond)

		close(chString)

		i++

		Println(i)
	}

}

func foo(chString chan<- string, alive *bool) {

	for *alive {

		chString <- "teste"

	}

	i--
}
