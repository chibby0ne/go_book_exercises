// Following the approach of mirroredQuery in Section 8.4.4, implement a variant of fetch that requests several URLs concurrently. As soon as the first response arrives, cancel the other requests.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	cancel    = make(chan struct{})
	responses = make(chan *http.Response, 3)
	wg        sync.WaitGroup
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Need at least 2 or more urls to fetch concurrently as command line arguments\n")
	}
	for _, url := range os.Args[1:] {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Print(err)
			}
			req.Cancel = cancel
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Print(err)
			}
			responses <- resp
		}(url)
	}
	resp := <-responses
	close(cancel)
	fmt.Printf("%v was the fastest\n", resp.Request.Host)
	wg.Wait()
}
