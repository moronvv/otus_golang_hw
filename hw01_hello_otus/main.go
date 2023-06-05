package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	str := "Hello, OTUS!"
	reversedStr := stringutil.Reverse(str)

	fmt.Print(reversedStr)
}
