package main_test

import (
	popcount "github.com/chibby0ne/go_book_exercises/chapter11/exercise11_6"
	"testing"
)

func runPopCountImplementation(t *testing.T, f func(uint64) int) {
	tests := []struct {
		input uint64
		want  int
	}{
		{81, 3},
		{123, 6},
		{437, 6},
		{412384, 8},
		{1000123, 12},
		{78956781234, 17},
		{123890491234, 21},
		{12389049123432141234, 34},
		{15986149086411123128, 37},
	}

	for _, test := range tests {
		if got := f(test.input); got != test.want {
			t.Errorf("PopCount(%d) = %d, want %d", test.input, got, test.want)
		}
	}
}

func TestPopCount(t *testing.T) {
	runPopCountImplementation(t, popcount.PopCount)
}

func TestPopCountShiftLoop(t *testing.T) {
	runPopCountImplementation(t, popcount.PopCountShiftLoop)
}

func TestPopCountClearingLoop(t *testing.T) {
	runPopCountImplementation(t, popcount.PopCountClearingLoop)
}

func benchmark(b *testing.B, f func(uint64) int, times int) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < times; j++ {
			f(uint64(j))
		}
	}
}

//
// Benchmarks:
// goos: linux
// goarch: amd64
// pkg: github.com/chibby0ne/go_book_exercises/chapter11/exercise11_6
// BenchmarkPopCount1-8                            315182368                3.81 ns/op
// BenchmarkPopCount10-8                           35638393                33.6 ns/op
// BenchmarkPopCount100-8                           3526356               340 ns/op
// BenchmarkPopCount1000-8                           356150              3342 ns/op
// BenchmarkPopCount10000-8                           36013             33333 ns/op
// BenchmarkPopCount100000-8                           3484            333665 ns/op
// BenchmarkPopCountShiftLoop1-8                   537559364                2.21 ns/op
// BenchmarkPopCountShiftLoop10-8                  35536447                33.3 ns/op
// BenchmarkPopCountShiftLoop100-8                  2245789               534 ns/op
// BenchmarkPopCountShiftLoop1000-8                  184384              6488 ns/op
// BenchmarkPopCountShiftLoop10000-8                  14928             80337 ns/op
// BenchmarkPopCountShiftLoop100000-8                  1242            961462 ns/op
// BenchmarkPopCountClearingLoop1-8                554151693                2.21 ns/op
// BenchmarkPopCountClearingLoop10-8               47114060                25.5 ns/op
// BenchmarkPopCountClearingLoop100-8               3348643               357 ns/op
// BenchmarkPopCountClearingLoop1000-8               308137              3878 ns/op
// BenchmarkPopCountClearingLoop10000-8               26122             45995 ns/op
// BenchmarkPopCountClearingLoop100000-8               2066            567241 ns/op
// PASS
// ok      github.com/chibby0ne/go_book_exercises/chapter11/exercise11_6   25.629s
//

// Seems like the table based approach (PopCount) starts to break even between
// the values of 100 and 1000

//
// PopCount
//
func BenchmarkPopCount1(b *testing.B) {
	benchmark(b, popcount.PopCount, 1)
}

func BenchmarkPopCount10(b *testing.B) {
	benchmark(b, popcount.PopCount, 10)
}

func BenchmarkPopCount100(b *testing.B) {
	benchmark(b, popcount.PopCount, 100)
}

func BenchmarkPopCount1000(b *testing.B) {
	benchmark(b, popcount.PopCount, 1000)
}

func BenchmarkPopCount10000(b *testing.B) {
	benchmark(b, popcount.PopCount, 10000)
}

func BenchmarkPopCount100000(b *testing.B) {
	benchmark(b, popcount.PopCount, 100000)
}

//
// PopCountShiftLoop
//
func BenchmarkPopCountShiftLoop1(b *testing.B) {
	benchmark(b, popcount.PopCountShiftLoop, 1)
}

func BenchmarkPopCountShiftLoop10(b *testing.B) {
	benchmark(b, popcount.PopCountShiftLoop, 10)
}

func BenchmarkPopCountShiftLoop100(b *testing.B) {
	benchmark(b, popcount.PopCountShiftLoop, 100)
}

func BenchmarkPopCountShiftLoop1000(b *testing.B) {
	benchmark(b, popcount.PopCountShiftLoop, 1000)
}

func BenchmarkPopCountShiftLoop10000(b *testing.B) {
	benchmark(b, popcount.PopCountShiftLoop, 10000)
}

func BenchmarkPopCountShiftLoop100000(b *testing.B) {
	benchmark(b, popcount.PopCountShiftLoop, 100000)
}

//
// PopCountClearingLoop
//
func BenchmarkPopCountClearingLoop1(b *testing.B) {
	benchmark(b, popcount.PopCountClearingLoop, 1)
}

func BenchmarkPopCountClearingLoop10(b *testing.B) {
	benchmark(b, popcount.PopCountClearingLoop, 10)
}

func BenchmarkPopCountClearingLoop100(b *testing.B) {
	benchmark(b, popcount.PopCountClearingLoop, 100)
}

func BenchmarkPopCountClearingLoop1000(b *testing.B) {
	benchmark(b, popcount.PopCountClearingLoop, 1000)
}

func BenchmarkPopCountClearingLoop10000(b *testing.B) {
	benchmark(b, popcount.PopCountClearingLoop, 10000)
}

func BenchmarkPopCountClearingLoop100000(b *testing.B) {
	benchmark(b, popcount.PopCountClearingLoop, 100000)
}
