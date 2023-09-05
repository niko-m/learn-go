package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	echo1(os.Args[1:])
	echo2(os.Args[1:])
	echo3(os.Args[1:])
}

func echo1(args []string) {
	s, sep := "", ""

	for _, arg := range args {
		s += sep + arg
		sep = " "
	}

	fmt.Println(s)
}

func echo2(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func echo3(args []string) {
	fmt.Println(args)
}
