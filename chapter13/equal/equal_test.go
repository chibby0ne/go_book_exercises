package equal

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




