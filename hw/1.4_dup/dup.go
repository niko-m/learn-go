package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	files := os.Args[1:]

	if len(files) == 0 {
		fmt.Println("no files passed")
		return
	}

	out := make([]string, 0, len(files))

	for _, v := range files {
		f, err := os.Open(v)

		if err != nil {
			fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
			continue
		}

		if hasDup(f) {
			out = append(out, v)
		}

		f.Close()
	}

	for _, file := range out {
		fmt.Println(file)
	}
}

func hasDup(f *os.File) bool {
	counts := make(map[string]int)
	input := bufio.NewScanner(f)

	for input.Scan() {
		text := input.Text()
		counts[text]++
		if counts[text] > 1 {
			return true
		}
	}

	return false
}
