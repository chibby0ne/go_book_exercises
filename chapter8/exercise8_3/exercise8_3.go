// In netcat3, the interface value conn has the concrete type *net.TCPConn,
// which represents a TCP connection. A TCP connection consists of two halves
// that may be closed independently using its CloseRead and CloseWrite methods.
// Modify the main goroutine of netcat3 to close only the write half of the
// connection so that the program will continue to print the final echoes from
// the reverb1 server even after the standard input has been closed.

package main

import (
	"io"
	"log"
	"net"
	"os"
)

func mustCopy(dst io.Writer, src io.Reader) {
	var n int64
	var err error
	if n, err = io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
	log.Printf("bytes copied: %v", n)
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn)
		log.Println("done")
		done <- struct{}{}
	}()
	mustCopy(conn, os.Stdin)
	if conn, ok := conn.(*net.TCPConn); ok {
		conn.CloseWrite()
	}
	log.Print("closed conn writer closer. waiting for echoes to finish...")
	<-done
}
