// Failure of any client program to read data in a timely manner ultimately
// causes all clients to get struck. Modify the broadcaster to skip a message
// rather than wait if a client writer is not ready to accept it.
// Alternatively, add buffering to each client's outgoing message channel so
// that most messages are not dropped; the broadcaster should use a
// non-blocking send to this channel
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

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
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		// time.Sleep(5 * time.Second) // to simulate a client program which can't read data in a timely manner, write from two clients quickly so that it buffers
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	bufferedMessages := make(map[client][]string)
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				select {
				case cli <- msg:
					for _, oldmsg := range bufferedMessages[cli] {
						cli <- "buffered message:" + oldmsg
					}
				default:
					buffer := bufferedMessages[cli]
					buffer = append(buffer, msg)
					bufferedMessages[cli] = buffer
				}
			}
		case cli := <-entering:
			clients[cli] = true
			bufferedMessages[cli] = []string{}
		case cli := <-leaving:
			delete(clients, cli)
			delete(bufferedMessages, cli)
			close(cli)
		}
	}
}
