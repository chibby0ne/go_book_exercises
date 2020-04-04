// Modify the reverb2 server to use a sync.WaitGroup per connection to count
// the number of active echo goroutines. When it falls to zero, close the write
// half of the TCP connection as described in Exercise 8.3. Verify that yoiur
// modified netcat3 client from that exercise waits for the final echoes of
// multiple concurrent shouts, even after the standard input has been closed.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration, wg *sync.WaitGroup) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
	wg.Done()
}

func handleConn(c net.Conn) {
	log.Printf("got connection from: %v", c.RemoteAddr())
	var wg sync.WaitGroup
	input := bufio.NewScanner(c)
	for input.Scan() {
		log.Printf("input.Text(): %s", input.Text())
		wg.Add(1)
		go echo(c, input.Text(), 1*time.Second, &wg)
	}
	wg.Wait()
	if c, ok := c.(*net.TCPConn); ok {
		// Note: ignoring potential errors from input.Err()
		c.CloseWrite()
	}
	log.Print("exiting handleConn")
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go handleConn(conn)
	}
}
