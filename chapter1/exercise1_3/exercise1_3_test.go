package main

import (
	"os"
	"testing"
)

func BenchmarkJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		using_strings_join(os.Args[1:])
	}
}

func BenchmarkConcat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		using_range_concat(os.Args[1:])
	}
}
