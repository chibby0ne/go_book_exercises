// The LimitReader function in the io package accepts an io.Reader r and a
// number of bytes n, and returns another Reader that reads from r but reports
// and end-of-file condition after n bytes. Implement it.
//
// func LimitReader(r io.Reader, n int64) io.Reader
//
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	BytesToRead = 20
)

type LimitedReader struct {
	r        io.Reader
	n, limit int64
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &LimitedReader{r: r, limit: n}
}

func (lr *LimitedReader) Read(p []byte) (int, error) {
	n, err := lr.r.Read(p[:lr.limit])
	lr.n += int64(n)
	if lr.n >= lr.limit {
		return n, io.EOF
	}
	return n, err
}

func main() {
	r := strings.NewReader("some io.Reader stream to be read\n")
	reader := LimitReader(r, BytesToRead)
	var b [BytesToRead * 2]byte
	n, err := reader.Read(b[:])
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	if n == BytesToRead && err == io.EOF {
		fmt.Printf("Read the limit set of %v bytes: \"%v\"\n", BytesToRead, string(b[:]))
		os.Exit(0)
	} else if err == io.EOF {
		fmt.Printf("Read only %v bytes: \"%v\"\n", n, string(b[:]))
		os.Exit(0)
	}
}
