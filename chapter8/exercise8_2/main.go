// Implement a concurrent FTP sever. The server should interpret commands from
// each client such as cd to change directory, ls to list a directory, get to
// send the contents of a file, and close to close the connection. you can use
// the standard ftp command as the client or write your own.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	f "github.com/chibby0ne/go_book_exercises/chapter8/exercise8_2/flag"
	"github.com/chibby0ne/go_book_exercises/chapter8/exercise8_2/ftp"
	l "github.com/chibby0ne/go_book_exercises/chapter8/exercise8_2/log"
)

func main() {
	flag.Parse()
	address := fmt.Sprintf("localhost:%d", *f.Port)
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		log.Fatal(err)
	}
	l.LogfVerbose("listening to %s...\n", address)
	ftp.ServeFTP(listener)
}
