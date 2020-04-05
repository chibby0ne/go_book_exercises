package main

import (
	"fmt"
	"github.com/chibby0ne/go_book_exercises/chapter5/links"
	"log"
	"os"
)

// tokens is a counting semaphore used to enforce a limit of 20 concurrent
// requests.
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	worklist := make(chan []string)
	var n int // number of pending sends to worklist. keeps track of the number of sends to worklist that are yet to occur

	// Start with the command-line arguments
	n++
	go func() { worklist <- os.Args[1:] }()

	seen := make(map[string]bool)
	// Crawl the web concurrently
	// Program terminates when it has discovered all the links reachable from the initial URLs
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}
