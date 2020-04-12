package main

import (
	"fmt"
	"github.com/chibby0ne/go_book_exercises/chapter5/links"
	"log"
	"os"
)

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	worklist := make(chan []string)  // list of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	// Start with the command-line arguments
	go func() { worklist <- os.Args[1:] }()

	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				// The foundlinks are sent to worklist in a separate goroutine
				// to avoid possible deadlock, since worklist and unseenLinks
				// are unbuffered channels, and these goroutine will block
				// until a message is received from worklist in the main
				// goroutine and therefore would block all these 20 goroutines
				// from crawling, also the main goroutine would block if there
				// was already a message sent on unseenLinks and is waiting for
				// one of the 20 goroutines to receive it -> deadlock
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// The main goroutine de-duplicates worklist items and sends the unseen
	// ones to the crawlers
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}
