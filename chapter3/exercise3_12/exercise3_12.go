package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println(are_anagrams("alo", "ola"))
	fmt.Println(are_anagrams("kalo", "pola"))
	fmt.Println(are_anagrams("lamina", "animal"))
	fmt.Println(are_anagrams("asd", "dpas"))
	fmt.Println(are_anagrams("123", "320"))
}

func are_anagrams(s, t string) (string, string, bool) {
	s_slice := []byte(s)
	t_slice := []byte(t)
	sort.Slice(s_slice, func(i, j int) bool { return s_slice[i] < s_slice[j] })
	sort.Slice(t_slice, func(i, j int) bool { return t_slice[i] < t_slice[j] })
	for i, l := range s_slice {
		if l != t_slice[i] {
			return s, t, false
		}
	}
	return s, t, true
}
