package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	var s string
	s = "Hello, OTUS!"
	s = stringutil.Reverse(s)
	fmt.Println(s)
}
