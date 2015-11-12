package main

import (
	. "fmt"
	"raft"
)

func main() {

	m := raft.Message{"port", " ", "state"}

	Printf("%v\n", m)

	Println(m.ToString())

	n := raft.DecomposeMessage(m.ToString())

	Printf("%+v\n", n)
}
