package util

// import (
// 	. "fmt"
// )

func CreateSeparator(separate int32) func(rune) bool {

	separate = rune(separate)

	return func(char rune) bool {
		return (char == separate)
	}

}

func EmptyChan(ch chan string) {

	for len(ch) > 0 {
		<-ch
	}

}
