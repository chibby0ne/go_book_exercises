// Make the chat server disconnect idle clients, such as those that have sent
// no messages in the last five minutes. Hint: calling conn.Close() in another
// goroutine unblocks active Read calls such as the one done by input.Scan()

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

const IdleTime = 5 * time.Minute

type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	messageSent := make(chan struct{})
	go clientWriter(conn, ch)
	go disconnectIdleConn(conn, ch, messageSent)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
		messageSent <- struct{}{}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func disconnectIdleConn(conn net.Conn, ch chan string, messageSent <-chan struct{}) {
	var lastSent time.Time
	timer := time.NewTimer(IdleTime)
	for {
		select {
		case <-timer.C:
			if time.Since(lastSent) >= IdleTime {
				ch <- "You are being disconnected for being idle"
				conn.Close()
			}
		case <-messageSent:
			t := time.Now()
			lastSent = t
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(IdleTime)
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}
