package main

import (
	. "fmt"
	"time"
	// "util"
)

var i int

func main() {

	alive := make([]bool, 10)

	for j := 0; j < 10; j++ {

		alive[j] = true

		go foo(&alive[j], j)

		i++

		Println(i)

	}

	time.Sleep(time.Second)

	for j := 0; j < 10; j++ {

		alive[j] = false

	}

	Println(i)

}

func foo(alive *bool, index int) (x int) {

	for *alive {

		x++

		Println(index, *alive)

	}

	Println(index, *alive)

	i--

	return

}
