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
