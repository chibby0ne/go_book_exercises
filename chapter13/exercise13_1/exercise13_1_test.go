package exercise13_1

import (
    "testing"
)

type link struct {
    value string
    tail *link
}

func TestEqual(t *testing.T) {
    if got := Equal([]int{1, 2, 3}, []int{1, 2, 3}); !got {
        t.Errorf("Equal(%v, %v), got: %v", []int{1, 2, 3}, []int{1, 2, 3}, got)
    }
    if got := Equal([]string{"foo"}, []string{"bar"}); got {
        t.Errorf("Equal(%v, %v), got: %v", []string{"foo"}, []string{"bar"}, got)
    }
    if got := Equal([]string(nil), []string{}); !got {
        t.Errorf("Equal(%v, %v), got: %v", []string(nil), []string{}, got)
    }
    if got := Equal(map[string]int(nil), map[string]int{}); !got {
        t.Errorf("Equal(%v, %v), got: %v", map[string]int(nil), map[string]int{}, got)
    }
    a, b, c := &link{value: "a"}, &link{value: "b"}, &link{value: "c"}
    a.tail, b.tail, c.tail = b, a, c
    t.Logf("Equal(a, a): %v\n", Equal(a, a))
    t.Logf("Equal(b, b): %v\n", Equal(b, b))
    t.Logf("Equal(c, c): %v\n", Equal(c, c))
    t.Logf("Equal(a, b): %v\n", Equal(a, b))
    t.Logf("Equal(a, c): %v\n", Equal(a, c))

}


func TestEqualNumbers(t *testing.T) {
    
    // expect true
    if got := Equal(int(1e12), int(1e12 -1)); !got {
        t.Errorf("Equal(%v, %v), got: %v", 1e12, 1e12 -1, got)
    }

    // expect false
    if got := Equal(int(1e2), int(1e2 -1)); got {
        t.Errorf("Equal(%v, %v), got: %v", 1e2, 1e2 -1, got)
    }

    // expect true
    if got := Equal(int64(1e12), int64(1e12-1)); !got {
        t.Errorf("Equal(%v, %v), got: %v", 1e12, 1e12 -1, got)
    }

    // expect false
    if got := Equal(int32(3000), int32(3001)); got {
        t.Errorf("Equal(%v, %v), got: %v", 3000, 3001, got)
    }

    // expect true
    if got := Equal(float64(3e12), float64(3e12 - 1)); !got {
        t.Errorf("Equal(%v, %v), got: %v", float64(3e12), float64(3e12-1), got)
    }
    
    // expect false
    if got := Equal(uint32(3000), uint32(3001)); got {
        t.Errorf("Equal(%v, %v), got: %v", 3000, 3001, got)
    }
}


