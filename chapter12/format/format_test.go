package format_test

import (
	"github.com/chibby0ne/go_book_exercises/chapter12/format"
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	var x int64 = 1
	var d time.Duration = 1 * time.Nanosecond
	t.Log(format.Any(x))
	t.Log(format.Any(d))
	t.Log(format.Any([]int64{x}))
	t.Log(format.Any([]time.Duration{d}))
}
