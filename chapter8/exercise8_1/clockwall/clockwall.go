// Modify clock2 to accept a port number, and write a program, clockwall, that
// acts as a client of several clock servers at once, reading the times from
// each one and displays the result in a table, akin to the wall of clocks seen
// in some business offices. If you have access to geographically distributed
// computers, run instances remotely; otherwise run local instances on
// different ports with fake time zones
//
// $ TZ=US/Eastern ./clock2 -port 8010 &
// $ TZ=Asia/Tokyo ./clock2 -port 8020 &
// $ TZ=Europe/London ./clock2 -port 8010 &
// $ clockwall NewYork=localhost:8010 Tokyo=localhost:8020 London=localhost:8030
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	locationsAddresses := make(map[string]string)
	for _, input := range os.Args[1:] {
		location := strings.Split(input, "=")
		locationsAddresses[location[0]] = location[1]
	}
	showLocationsTime(locationsAddresses)
}

func showLocationsTime(locationsAddresses map[string]string) {
	for location, address := range locationsAddresses {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		go showLocationTime(os.Stdout, conn, location)
	}
	for {
	}
}

func showLocationTime(dst io.Writer, src io.Reader, location string) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		fmt.Fprintf(dst, "%s: %s\n", location, scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Fprintf(dst, "error while reading from location:%s: %s\n", location, scanner.Err())
	}
}
