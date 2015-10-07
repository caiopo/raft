package main

import (
	. "fmt"
	"raft/leader"
	"strings"
	"util"
)

func main() {

	separator := util.CreateSeparator(':')
	Println(strings.FieldsFunc("teste:a:b:gasd", separator))

	Println(leader.Abc())

}
