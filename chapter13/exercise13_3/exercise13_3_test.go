package bzip

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"fmt"
	"io"
	"strconv"
	"sync"
	"testing"
)

func initializeChan(numMessages int) chan int {
	ch := make(chan int, numMessages)
	for i := 0; i < numMessages; i++ {
		ch <- i
	}
	close(ch)
	return ch
}

func writeBzip(wg *sync.WaitGroup, ch chan int, writer io.WriteCloser, t *testing.T) {
	defer wg.Done()
	for i := range ch {
		// Need the space for parsing correctly using scanner.Text() due to the
		// split func of the scanner being bufio.ScanWords
		_, err := writer.Write([]byte(fmt.Sprintf("%d ", i)))
		if err != nil {
			t.Errorf("Error writing using bzip writer: %v", err)
		}
	}
}

const (
	MaxSizeChannel int = 50
	NumRoutines    int = 3
)

func TestConcurrentUse(t *testing.T) {
	// Initialize channel which will be used to get messages from different
	// goroutines competing for writing with the writer
	ch := initializeChan(MaxSizeChannel)

	// Writer that writes to buffer
	var buffer bytes.Buffer
	writer := NewWriter(&buffer)

	// For keeping track and waiting for the goroutines to finish
	var wg sync.WaitGroup

	// Create the goroutines
	for i := 0; i < NumRoutines; i++ {
		wg.Add(1)
		go writeBzip(&wg, ch, writer, t)
	}

	// Wait for the goroutines to end
	wg.Wait()
	writer.Close()

	// uncompress and add them to a map to check which are missing
	numbers := make(map[int]bool)
	var uncompressed bytes.Buffer

	io.Copy(&uncompressed, bzip2.NewReader(&buffer))

	scanner := bufio.NewScanner(&uncompressed)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			t.Errorf("Number was not written correctly: %v", err)
		}
		numbers[num] = true
	}

	// report which numbers were not written
	for i := 0; i < MaxSizeChannel; i++ {
		if !numbers[i] {
			t.Errorf("Missing number: %v", i)
		}
	}

}
