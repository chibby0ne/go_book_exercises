// Add depth-limiting to the concurrent crawler. That is, if the user sets
// -depth=3 then only URLs reachable by at most threee links will be fetched.
package main

import (
	"flag"
	"fmt"
	"github.com/chibby0ne/go_book_exercises/chapter5/links"
	"log"
	"sync"
)

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

type DepthMap struct {
	m map[string]int
	l sync.Mutex
}

func (dm *DepthMap) Store(link string, depth int) {
	defer dm.l.Unlock()
	dm.l.Lock()
	dm.m[link] = depth
}

func (dm *DepthMap) Check(link string) (int, bool) {
	defer dm.l.Unlock()
	dm.l.Lock()
	v, ok := dm.m[link]
	return v, ok
}

func NewDepthMap() DepthMap {
	return DepthMap{
		m: make(map[string]int),
		l: sync.Mutex{},
	}
}

func main() {
	maxDepth := flag.Int("depth", 3, "depth of the crawler, it will search at most this many links deep")

	flag.Parse()
	worklist := make(chan []string)  // list of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs
	n := 0                           // number of pending sends to worklist and also within the specified depth
	depthMap := NewDepthMap()

	// Initialize the depth map with level 0
	for _, link := range flag.Args() {
		depthMap.Store(link, 0)
	}

	// Start with the command-line arguments
	n++
	go func() { worklist <- flag.Args() }()

	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				depth, _ := depthMap.Check(link)
				if depth >= *maxDepth {
					continue
				}
				foundLinks := crawl(link)
				go func(link string, depth int) {
					for _, newLink := range foundLinks {
						if _, ok := depthMap.Check(newLink); !ok {
							depthMap.Store(newLink, depth+1)
						}
					}
					worklist <- foundLinks
				}(link, depth)
			}
		}()
	}

	// The main goroutine de-duplicates worklist items and sends the unseen
	// ones to the crawlers
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			depth, _ := depthMap.Check(link)
			if !seen[link] && depth < *maxDepth {
				n++
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
	close(unseenLinks)
	fmt.Println("That's all")
}
