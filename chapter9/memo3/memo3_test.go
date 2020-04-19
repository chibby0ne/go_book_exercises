package memo

import (
	"testing"
)

// This test will report a race condition in line 36
func TestConcurrent(t *testing.T) {
	input := []string{
		"https://golang.org",
		"https://godoc.org",
		"https://play.golang.org",
		"http://gopl.io",
		"https://golang.org",
		"https://godoc.org",
		"https://play.golang.org",
		"http://gopl.io",
	}
	for _, v := range input {
		ch <- v
	}
	close(ch)
	Concurrent()
}
