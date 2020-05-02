package main

import (
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
