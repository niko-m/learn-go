package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Print some text below:")

	r := bufio.NewReader(os.Stdin)

	for {
		data, prefix, err := r.ReadLine()

		if len(data) == 0 {
			break
		}

		fmt.Printf("string: (%v) %v\nprefix: %v\nerr: %v\n\n", len(data), string(data), prefix, err)
	}
}
