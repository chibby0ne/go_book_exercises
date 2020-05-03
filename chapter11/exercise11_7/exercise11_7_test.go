package main

import (
	"math/rand"
	"testing"
)

func initMapSet(vals []int) *MapSet {
	var x MapSet
	for _, v := range vals {
		x.Add(v)
	}
	return &x
}

func initIntSet(vals []int) *IntSet {
	var x IntSet
	for _, v := range vals {
		x.Add(v)
	}
	return &x
}

func TestIntSetHas(t *testing.T) {
	var tests = []struct {
		input  []int
		checks []int
	}{
		{
			[]int{4, 8, 15, 16},
			[]int{0, 4, 30, 15},
		},
		{
			[]int{4, 8, 15, 16},
			[]int{4, 8, 15, 16},
		},
	}
	for _, test := range tests {
		mapSet := initMapSet(test.input)
		intSet := initIntSet(test.input)
		for _, v := range test.checks {
			got := intSet.Has(v)
			want := mapSet.Has(v)
			if got != want {
				t.Errorf("intSet: %v intSet.Has(%v) = %v", intSet, v, got)
			}
		}
	}
}

func TestIntSetAdd(t *testing.T) {
	var tests = []struct {
		input []int
	}{
		{
			[]int{4, 8, 15, 16},
		},
		{
			[]int{40, 8, 48, 160},
		},
	}
	for _, test := range tests {
		var intSet IntSet
		var mapSet MapSet
		for _, in := range test.input {
			intSet.Add(in)
			mapSet.Add(in)
			if intSet.String() != mapSet.String() {
				t.Errorf("intSet: %v mapSet: %v", intSet, mapSet)
			}
		}
	}
}

func TestIntSetUnionWith(t *testing.T) {
	var tests = []struct {
		input1 []int
		input2 []int
	}{
		{
			[]int{4, 8, 15, 16},
			[]int{0, 4, 30, 40},
		},
		{
			[]int{4, 8, 100, 16},
			[]int{0, 4, 30, 38},
		},
	}
	for _, test := range tests {
		intSet1 := initIntSet(test.input1)
		intSet2 := initIntSet(test.input2)
		mapSet1 := initMapSet(test.input1)
		mapSet2 := initMapSet(test.input2)
		intSet1.UnionWith(intSet2)
		mapSet1.UnionWith(mapSet2)

		if intSet1.String() != mapSet1.String() {
			t.Errorf("intSet1: %v mapSet1: %v", intSet1, mapSet1)
		}
	}
}

//
// Benchmarks
//

var rng = rand.New(rand.NewSource(4))

func benchmarkAdd(b *testing.B, s Set, upTo int64) {
	for i := 0; i < b.N; i++ {
		s.Add(int(rng.Int63n(upTo)))
	}
}

func benchmarkHas(b *testing.B, s Set, upTo int64) {
	for i := 0; i < b.N; i++ {
		s.Has(int(rng.Int63n(upTo)))
	}
}

func benchmarkUnionWith(b *testing.B, s Set, times int) {
	var set Set
	switch s.(type) {
	case *IntSet:
		set = &IntSet{}
	case *MapSet:
		set = &MapSet{}
	default:
		b.Fatal("this should not happen")
	}
	for i := 0; i < b.N; i++ {
		val := rng.Int63n(int64(times))
		set.Add(int(val))
		s.UnionWith(set)
	}
}

//
// Benchmarks for IntSet
//
func BenchmarkIntSetAddIntsUpTo100(b *testing.B) {
	var intSet IntSet
	benchmarkAdd(b, &intSet, 100)
}

func BenchmarkIntSetAddIntsUpTo1000(b *testing.B) {
	var intSet IntSet
	benchmarkAdd(b, &intSet, 1000)
}

func BenchmarkIntSetAddIntsUpTo10000(b *testing.B) {
	var intSet IntSet
	benchmarkAdd(b, &intSet, 10000)
}

func BenchmarkIntSetAddIntsUpTo100000(b *testing.B) {
	var intSet IntSet
	benchmarkAdd(b, &intSet, 100000)
}

func BenchmarkIntSetAddIntsUpTo1000000(b *testing.B) {
	var intSet IntSet
	benchmarkAdd(b, &intSet, 1000000)
}

func BenchmarkIntSetAddIntsUpTo10000000(b *testing.B) {
	var intSet IntSet
	benchmarkAdd(b, &intSet, 10000000)
}

func BenchmarkIntSetAddIntsUpTo100000000(b *testing.B) {
	var intSet IntSet
	benchmarkAdd(b, &intSet, 100000000)
}

func BenchmarkIntSetAddIntsUpTo1000000000(b *testing.B) {
	var intSet IntSet
	benchmarkAdd(b, &intSet, 1000000000)
}

func BenchmarkIntSetHasIntsUpTo100(b *testing.B) {
	var intSet IntSet
	benchmarkHas(b, &intSet, 100)
}

func BenchmarkIntSetHasIntsUpTo1000(b *testing.B) {
	var intSet IntSet
	benchmarkHas(b, &intSet, 1000)
}

func BenchmarkIntSetHasIntsUpTo10000(b *testing.B) {
	var intSet IntSet
	benchmarkHas(b, &intSet, 10000)
}

func BenchmarkIntSetHasIntsUpTo100000(b *testing.B) {
	var intSet IntSet
	benchmarkHas(b, &intSet, 100000)
}

func BenchmarkIntSetHasIntsUpTo1000000(b *testing.B) {
	var intSet IntSet
	benchmarkHas(b, &intSet, 1000000)
}

func BenchmarkIntSetHasIntsUpTo10000000(b *testing.B) {
	var intSet IntSet
	benchmarkHas(b, &intSet, 10000000)
}

func BenchmarkIntSetHasIntsUpTo100000000(b *testing.B) {
	var intSet IntSet
	benchmarkHas(b, &intSet, 100000000)
}

func BenchmarkIntSetHasIntsUpTo1000000000(b *testing.B) {
	var intSet IntSet
	benchmarkHas(b, &intSet, 1000000000)
}

func BenchmarkIntSetUnionWithIntsUpTo100(b *testing.B) {
	var intSet IntSet
	benchmarkUnionWith(b, &intSet, 100)
}

func BenchmarkIntSetUnionWithIntsUpTo1000(b *testing.B) {
	var intSet IntSet
	benchmarkUnionWith(b, &intSet, 1000)
}

func BenchmarkIntSetUnionWithIntsUpTo10000(b *testing.B) {
	var intSet IntSet
	benchmarkUnionWith(b, &intSet, 10000)
}

func BenchmarkIntSetUnionWithIntsUpTo100000(b *testing.B) {
	var intSet IntSet
	benchmarkUnionWith(b, &intSet, 100000)
}

func BenchmarkIntSetUnionWithIntsUpTo1000000(b *testing.B) {
	var intSet IntSet
	benchmarkUnionWith(b, &intSet, 1000000)
}

func BenchmarkIntSetUnionWithIntsUpTo10000000(b *testing.B) {
	var intSet IntSet
	benchmarkUnionWith(b, &intSet, 10000000)
}

func BenchmarkIntSetUnionWithIntsUpTo100000000(b *testing.B) {
	var intSet IntSet
	benchmarkUnionWith(b, &intSet, 100000000)
}

func BenchmarkIntSetUnionWithIntsUpTo1000000000(b *testing.B) {
	var intSet IntSet
	benchmarkUnionWith(b, &intSet, 1000000000)
}

//
// Benchmarks for MapSet
//
func BenchmarkMapSetAddIntsUpTo100(b *testing.B) {
	var mapSet MapSet
	benchmarkAdd(b, &mapSet, 100)
}

func BenchmarkMapSetAddIntsUpTo1000(b *testing.B) {
	var mapSet MapSet
	benchmarkAdd(b, &mapSet, 1000)
}

func BenchmarkMapSetAddIntsUpTo10000(b *testing.B) {
	var mapSet MapSet
	benchmarkAdd(b, &mapSet, 10000)
}

func BenchmarkMapSetAddIntsUpTo100000(b *testing.B) {
	var mapSet MapSet
	benchmarkAdd(b, &mapSet, 100000)
}

func BenchmarkMapSetAddIntsUpTo1000000(b *testing.B) {
	var mapSet MapSet
	benchmarkAdd(b, &mapSet, 1000000)
}

func BenchmarkMapSetAddIntsUpTo10000000(b *testing.B) {
	var mapSet MapSet
	benchmarkAdd(b, &mapSet, 10000000)
}

func BenchmarkMapSetAddIntsUpTo100000000(b *testing.B) {
	var mapSet MapSet
	benchmarkAdd(b, &mapSet, 100000000)
}

func BenchmarkMapSetAddIntsUpTo1000000000(b *testing.B) {
	var mapSet MapSet
	benchmarkAdd(b, &mapSet, 1000000000)
}

func BenchmarkMapSetHasIntsUpTo100(b *testing.B) {
	var mapSet MapSet
	benchmarkHas(b, &mapSet, 100)
}

func BenchmarkMapSetHasIntsUpTo1000(b *testing.B) {
	var mapSet MapSet
	benchmarkHas(b, &mapSet, 1000)
}

func BenchmarkMapSetHasIntsUpTo10000(b *testing.B) {
	var mapSet MapSet
	benchmarkHas(b, &mapSet, 10000)
}

func BenchmarkMapSetHasIntsUpTo100000(b *testing.B) {
	var mapSet MapSet
	benchmarkHas(b, &mapSet, 100000)
}

func BenchmarkMapSetHasIntsUpTo1000000(b *testing.B) {
	var mapSet MapSet
	benchmarkHas(b, &mapSet, 1000000)
}

func BenchmarkMapSetHasIntsUpTo10000000(b *testing.B) {
	var mapSet MapSet
	benchmarkHas(b, &mapSet, 10000000)
}

func BenchmarkMapSetHasIntsUpTo100000000(b *testing.B) {
	var mapSet MapSet
	benchmarkHas(b, &mapSet, 100000000)
}

func BenchmarkMapSetHasIntsUpTo1000000000(b *testing.B) {
	var mapSet MapSet
	benchmarkHas(b, &mapSet, 1000000000)
}

func BenchmarkMapSetUnionWithIntsUpTo100(b *testing.B) {
	var mapSet MapSet
	benchmarkUnionWith(b, &mapSet, 100)
}

func BenchmarkMapSetUnionWithIntsUpTo1000(b *testing.B) {
	var mapSet MapSet
	benchmarkUnionWith(b, &mapSet, 1000)
}

func BenchmarkMapSetUnionWithIntsUpTo10000(b *testing.B) {
	var mapSet MapSet
	benchmarkUnionWith(b, &mapSet, 10000)
}

func BenchmarkMapSetUnionWithIntsUpTo100000(b *testing.B) {
	var mapSet MapSet
	benchmarkUnionWith(b, &mapSet, 100000)
}

func BenchmarkMapSetUnionWithIntsUpTo1000000(b *testing.B) {
	var mapSet MapSet
	benchmarkUnionWith(b, &mapSet, 1000000)
}

func BenchmarkMapSetUnionWithIntsUpTo10000000(b *testing.B) {
	var mapSet MapSet
	benchmarkUnionWith(b, &mapSet, 10000000)
}

func BenchmarkMapSetUnionWithIntsUpTo100000000(b *testing.B) {
	var mapSet MapSet
	benchmarkUnionWith(b, &mapSet, 100000000)
}

func BenchmarkMapSetUnionWithIntsUpTo1000000000(b *testing.B) {
	var mapSet MapSet
	benchmarkUnionWith(b, &mapSet, 1000000000)
}

//
// Having the UnionWith set with the add operations and the UnionWith operation done in the s set
//
// goos: linux
// goarch: amd64
// pkg: github.com/chibby0ne/go_book_exercises/chapter11/exercise11_7
// BenchmarkIntSetAddIntsUpTo100-8                 52394488                22.9 ns/op
// BenchmarkIntSetAddIntsUpTo1000-8                51494299                22.8 ns/op
// BenchmarkIntSetAddIntsUpTo10000-8               52289509                22.9 ns/op
// BenchmarkIntSetAddIntsUpTo100000-8              52094341                23.1 ns/op
// BenchmarkIntSetAddIntsUpTo1000000-8             52124618                22.9 ns/op
// BenchmarkIntSetAddIntsUpTo10000000-8            48956962                23.3 ns/op
// BenchmarkIntSetAddIntsUpTo100000000-8           27520870                43.8 ns/op
// BenchmarkIntSetAddIntsUpTo1000000000-8          15881080                72.7 ns/op
// BenchmarkIntSetHasIntsUpTo100-8                 59937487                19.5 ns/op
// BenchmarkIntSetHasIntsUpTo1000-8                59699850                19.6 ns/op
// BenchmarkIntSetHasIntsUpTo10000-8               61231165                19.6 ns/op
// BenchmarkIntSetHasIntsUpTo100000-8              61233337                19.6 ns/op
// BenchmarkIntSetHasIntsUpTo1000000-8             61188028                19.6 ns/op
// BenchmarkIntSetHasIntsUpTo10000000-8            59396986                19.7 ns/op
// BenchmarkIntSetHasIntsUpTo100000000-8           61338043                19.6 ns/op
// BenchmarkIntSetHasIntsUpTo1000000000-8          60916760                19.6 ns/op
// BenchmarkIntSetUnionWithIntsUpTo100-8           43060333                27.9 ns/op
// BenchmarkIntSetUnionWithIntsUpTo1000-8          30769065                38.7 ns/op
// BenchmarkIntSetUnionWithIntsUpTo10000-8          8347999               144 ns/op
// BenchmarkIntSetUnionWithIntsUpTo100000-8         1000000              1158 ns/op
// BenchmarkIntSetUnionWithIntsUpTo1000000-8         101601             11756 ns/op
// BenchmarkIntSetUnionWithIntsUpTo10000000-8          8708            116685 ns/op
// BenchmarkIntSetUnionWithIntsUpTo100000000-8          718           1518887 ns/op
// BenchmarkIntSetUnionWithIntsUpTo1000000000-8          55          18212293 ns/op
// BenchmarkMapSetAddIntsUpTo100-8                 26266586                46.3 ns/op
// BenchmarkMapSetAddIntsUpTo1000-8                26596009                45.4 ns/op
// BenchmarkMapSetAddIntsUpTo10000-8               24899136                47.7 ns/op
// BenchmarkMapSetAddIntsUpTo100000-8              20776390                54.9 ns/op
// BenchmarkMapSetAddIntsUpTo1000000-8             12736574                94.3 ns/op
// BenchmarkMapSetAddIntsUpTo10000000-8             9405726               129 ns/op
// BenchmarkMapSetAddIntsUpTo100000000-8            9284557               153 ns/op
// BenchmarkMapSetAddIntsUpTo1000000000-8           9103368               154 ns/op
// BenchmarkMapSetHasIntsUpTo100-8                 53536906                22.8 ns/op
// BenchmarkMapSetHasIntsUpTo1000-8                51573970                22.7 ns/op
// BenchmarkMapSetHasIntsUpTo10000-8               52191912                22.3 ns/op
// BenchmarkMapSetHasIntsUpTo100000-8              53721255                22.6 ns/op
// BenchmarkMapSetHasIntsUpTo1000000-8             52741389                22.2 ns/op
// BenchmarkMapSetHasIntsUpTo10000000-8            52826398                22.6 ns/op
// BenchmarkMapSetHasIntsUpTo100000000-8           52811232                22.3 ns/op
// BenchmarkMapSetHasIntsUpTo1000000000-8          53119748                22.7 ns/op
// BenchmarkMapSetUnionWithIntsUpTo100-8             346527              3482 ns/op
// BenchmarkMapSetUnionWithIntsUpTo1000-8             35773             36947 ns/op
// BenchmarkMapSetUnionWithIntsUpTo10000-8            10000            156514 ns/op
// BenchmarkMapSetUnionWithIntsUpTo100000-8           10000            219344 ns/op
// BenchmarkMapSetUnionWithIntsUpTo1000000-8          10000            223572 ns/op
// BenchmarkMapSetUnionWithIntsUpTo10000000-8         10000            223031 ns/op
// BenchmarkMapSetUnionWithIntsUpTo100000000-8        10000            225320 ns/op
// BenchmarkMapSetUnionWithIntsUpTo1000000000-8       10000            232264 ns/op
// PASS
// ok      github.com/chibby0ne/go_book_exercises/chapter11/exercise11_7   67.382s

//
//
// Having the UnionWith set with the add operations and the UnionWith operation done in the set Set (i.e: adding an empty Set)
//
// goos: linux
// goarch: amd64
// pkg: github.com/chibby0ne/go_book_exercises/chapter11/exercise11_7
// BenchmarkIntSetAddIntsUpTo100-8                 50385336                23.0 ns/op
// BenchmarkIntSetAddIntsUpTo1000-8                51667009                23.0 ns/op
// BenchmarkIntSetAddIntsUpTo10000-8               52210525                22.9 ns/op
// BenchmarkIntSetAddIntsUpTo100000-8              52156922                23.0 ns/op
// BenchmarkIntSetAddIntsUpTo1000000-8             51606097                23.0 ns/op
// BenchmarkIntSetAddIntsUpTo10000000-8            48750493                23.2 ns/op
// BenchmarkIntSetAddIntsUpTo100000000-8           27357582                44.0 ns/op
// BenchmarkIntSetAddIntsUpTo1000000000-8          14674228                72.9 ns/op
// BenchmarkIntSetHasIntsUpTo100-8                 61549233                19.5 ns/op
// BenchmarkIntSetHasIntsUpTo1000-8                60419586                19.5 ns/op
// BenchmarkIntSetHasIntsUpTo10000-8               60950173                19.4 ns/op
// BenchmarkIntSetHasIntsUpTo100000-8              61963610                19.6 ns/op
// BenchmarkIntSetHasIntsUpTo1000000-8             61378341                19.4 ns/op
// BenchmarkIntSetHasIntsUpTo10000000-8            59893634                19.5 ns/op
// BenchmarkIntSetHasIntsUpTo100000000-8           61334964                19.5 ns/op
// BenchmarkIntSetHasIntsUpTo1000000000-8          61584146                19.5 ns/op
// BenchmarkIntSetUnionWithIntsUpTo100-8           46418280                25.8 ns/op
// BenchmarkIntSetUnionWithIntsUpTo1000-8          45591897                26.3 ns/op
// BenchmarkIntSetUnionWithIntsUpTo10000-8         44225811                26.4 ns/op
// BenchmarkIntSetUnionWithIntsUpTo100000-8        46355760                25.8 ns/op
// BenchmarkIntSetUnionWithIntsUpTo1000000-8       45942476                25.6 ns/op
// BenchmarkIntSetUnionWithIntsUpTo10000000-8      43929394                26.2 ns/op
// BenchmarkIntSetUnionWithIntsUpTo100000000-8     25336431                47.8 ns/op
// BenchmarkIntSetUnionWithIntsUpTo1000000000-8    14300050                77.1 ns/op
// BenchmarkMapSetAddIntsUpTo100-8                 25005676                50.1 ns/op
// BenchmarkMapSetAddIntsUpTo1000-8                25618623                44.4 ns/op
// BenchmarkMapSetAddIntsUpTo10000-8               24975493                49.9 ns/op
// BenchmarkMapSetAddIntsUpTo100000-8              20945766                57.5 ns/op
// BenchmarkMapSetAddIntsUpTo1000000-8             12575019                93.5 ns/op
// BenchmarkMapSetAddIntsUpTo10000000-8             9427948               129 ns/op
// BenchmarkMapSetAddIntsUpTo100000000-8            9256411               154 ns/op
// BenchmarkMapSetAddIntsUpTo1000000000-8           9299414               155 ns/op
// BenchmarkMapSetHasIntsUpTo100-8                 52957160                22.7 ns/op
// BenchmarkMapSetHasIntsUpTo1000-8                51711829                22.6 ns/op
// BenchmarkMapSetHasIntsUpTo10000-8               52500178                22.6 ns/op
// BenchmarkMapSetHasIntsUpTo100000-8              52705354                22.6 ns/op
// BenchmarkMapSetHasIntsUpTo1000000-8             51617404                22.6 ns/op
// BenchmarkMapSetHasIntsUpTo10000000-8            52656842                22.6 ns/op
// BenchmarkMapSetHasIntsUpTo100000000-8           52720472                22.7 ns/op
// BenchmarkMapSetHasIntsUpTo1000000000-8          52609090                22.7 ns/op
// BenchmarkMapSetUnionWithIntsUpTo100-8           22740442                53.2 ns/op
// BenchmarkMapSetUnionWithIntsUpTo1000-8          22734658                52.6 ns/op
// BenchmarkMapSetUnionWithIntsUpTo10000-8         21614175                54.8 ns/op
// BenchmarkMapSetUnionWithIntsUpTo100000-8        18661153                62.1 ns/op
// BenchmarkMapSetUnionWithIntsUpTo1000000-8       11706801               104 ns/op
// BenchmarkMapSetUnionWithIntsUpTo10000000-8       8567504               139 ns/op
// BenchmarkMapSetUnionWithIntsUpTo100000000-8      8514171               165 ns/op
// BenchmarkMapSetUnionWithIntsUpTo1000000000-8     8633817               166 ns/op
// PASS
// ok      github.com/chibby0ne/go_book_exercises/chapter11/exercise11_7   64.244s

