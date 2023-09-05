package main

import (
	"os"
	"testing"
)

var input = []string{"hello", "world", "dude", "fuck", "you"}

func BenchmarkEcho1(t *testing.B) {
	os.Stdout, _ = os.Open(os.DevNull)

	for i := 0; i < t.N; i++ {
		echo1(input)
	}
}

func BenchmarkEcho2(t *testing.B) {
	os.Stdout, _ = os.Open(os.DevNull)

	for i := 0; i < t.N; i++ {
		echo2(input)
	}
}

func BenchmarkEcho3(t *testing.B) {
	os.Stdout, _ = os.Open(os.DevNull)

	for i := 0; i < t.N; i++ {
		echo3(input)
	}
}
