// Using a aslect statement, add a timeout to the echo server from Section 8.3
// so that it disconnects any client that shouts nothing withing 10 seconds.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	scanChan := make(chan bool)
	breakChan := make(chan struct{})

	// Note: ignoring potential errors from input.Err()
	defer c.Close()

	// Send input.Scan() events on channel so that events can be multiplexed
	// between timeout due to inactivity and text received, since we don't want
	// to block waiting for something to be sent, we send it in another
	// goroutine, additionally when input.Scan() returns false it means it
	// reached the end of an input or an error therefore we need to send to the
	// channel used for returning from this function and closing the connection
	// in the process.
	go func() {
		for input.Scan() {
			scanChan <- true
		}
		breakChan <- struct{}{}
	}()

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Printf("Disconnecting from client: %v due to inactivity\n", c.RemoteAddr())
			ticker.Stop()
			return
		case <-scanChan:
			ticker.Stop()
			fmt.Printf("Received: %v\n", input.Text())
			go echo(c, input.Text(), 1*time.Second)
			ticker = time.NewTicker(10 * time.Second)
		case <-breakChan:
			return
		}
	}
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
