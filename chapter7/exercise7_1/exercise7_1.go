// Using the ideas from ByteCounter, implement counters for words and for
// lines. You will find bufio.ScanWords useful.

package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type WordCounter int
type LineCounter int

func (c *WordCounter) Write(p []byte) (n int, err error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(p))
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}
	*c += WordCounter(count)
	if scanner.Err() != nil {
		return count, scanner.Err()
	}
	if len(p) != count {
		return count, fmt.Errorf("Could not write all bytes")
	}
	return count, nil
}

func (c *LineCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(p))
	scanner.Split(bufio.ScanLines)
	count := 0
	for scanner.Scan() {
		count++
	}
	*c += LineCounter(count)
	if scanner.Err() != nil {
		return count, scanner.Err()
	}
	if len(p) != count {
		return count, fmt.Errorf("Could not write all bytes")
	}
	return count, nil
}

func main() {
	fmt.Println("Word Counter")
	var c WordCounter
	c.Write([]byte("hello"))
	fmt.Println(c)
	c = 0

	var name = "Dolly"
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(c)

	fmt.Println("Now with line counter")

	var cl WordCounter
	cl.Write([]byte("hello\n"))
	fmt.Println(cl)

	cl = 0
	fmt.Fprintf(&cl, "hello,\n %s\n", name)
	fmt.Println(cl)

}
