package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	n = flag.Bool("n", false, "опустить символы новой строки")
	s = flag.String("s", " ", "разделитель")
)

var out io.Writer = os.Stdout

func main() {
	flag.Parse()

	err := echo(!*n, *s, flag.Args())

	if err != nil {
		fmt.Fprintf(os.Stderr, "echo: %v/n", err)
		os.Exit(1)
	}
}

func echo(newline bool, sep string, args []string) error {
	fmt.Fprint(out, strings.Join(args, sep))
	if newline {
		fmt.Fprintln(out)
	}
	return nil
}
