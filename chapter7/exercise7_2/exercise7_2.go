// Write a function CountingWriter with the signature below that, given an
// io.Writer returns a new Writer that wraps the original, and a a pointer to
// an int64 variable that at any moment contains the number of bytes writen to
// the new Writer.
//
// func CountingWriter(w io.Writer) (io.Writer, *int64)
//
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
)

type countWriter struct {
	writer io.Writer
	count  *int64
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	c := int64(0)
	writer := countWriter{
		writer: w,
		count:  &c,
	}
	return writer, writer.count
}

func (c countWriter) Write(p []byte) (int, error) {
	n, err := c.writer.Write(p)
	*c.count += int64(n)
	if err != nil {
		return n, err
	}
	return int(*c.count), nil

}

func main() {
	var buffer bytes.Buffer
	cw, numBytes := CountingWriter(&buffer)
	n, err := cw.Write([]byte("hello"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("num of bytes written is: %d, which should also be in this case n: %d\n", *numBytes, n)
}
