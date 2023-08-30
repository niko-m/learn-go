package main

import (
	"fmt"
	"os"

	"github.com/niko-m/learn-go/book0/ch11/word1"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no words is passed")
		return
	}

	fmt.Println(word1.IsPalindrome(os.Args[1]))
}
