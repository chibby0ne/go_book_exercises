// Extend the Func type and the (*Memo).Get method so that callers may provide an optional done channel through whicht hey can cancel the operation. The results of a cancelled Func call should not be cached.
package memo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

// Func is the type of the function to memoize
type Func func(key string, done <-chan struct{}) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

// A request is a message requesting that the Func be applied to Key
type request struct {
	key      string
	response chan<- result // the client wants a single result
	done     <-chan struct{}
}

// A Memo caches the results of calling a Func.
type Memo struct {
	requests, cancels chan request
}

func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request), cancels: make(chan request)}
	go memo.server(f)
	return memo
}

// Get is concurrency safe
func (memo *Memo) Get(key string, done <-chan struct{}) (interface{}, error) {
	response := make(chan result)
	req := request{key, response, done}
	memo.requests <- req
	res := <-response
	select {
	case <-done:
		memo.cancels <- req
	default:
		// Not canceled
	}
	return res.value, res.err
}

func (memo *Memo) Close() {
	close(memo.requests)
	close(memo.cancels)
}

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for {
	Cancel:
		// Need to process cancels first in case there are early cancels, since
		// if we allowed the processing of requests and cancels at the same
		// time select would randomly choose one or the other case (channel
		// event), and we can't guarantee that the request url will not be for
		// a cancelled request url. Whenever that channel is empty we process
		// to the main loop of the server which will queue between.
		for {
			select {
			case req := <-memo.cancels:
				delete(cache, req.key)
			default:
				break Cancel
			}
		}
		// Main loop of server. Handles requests and cancels concurrently as they come
		for {
			select {
			case req := <-memo.cancels:
				delete(cache, req.key)
				continue
			case req := <-memo.requests:
				e := cache[req.key]
				if e == nil {
					// This is the first reuqest for this key.
					e = &entry{ready: make(chan struct{})}
					cache[req.key] = e
					go e.call(f, req.key, req.done)
				}
				go e.deliver(req.response)
			}
		}

	}
}

func (e *entry) call(f Func, key string, done <-chan struct{}) {
	// Evaluate the function
	e.res.value, e.res.err = f(key, done)
	// Broadcast the ready condition
	close(e.ready)
}

func (e *entry) deliver(res chan<- result) {
	// Wait for ready condition
	<-e.ready
	// Send the result to the client
	res <- e.res
}

func httpGetBody(url string, done <-chan struct{}) (interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Cancel = done
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)

}

var doneChan chan struct{}

func Concurrent() {
	m := New(httpGetBody)
	var wg sync.WaitGroup
	for url := range incomingURLs() {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			start := time.Now()
			value, err := m.Get(url, doneChan)
			if err != nil {
				log.Print(err)
				return
			}
			fmt.Printf("%s, %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))
		}(url)
	}
	wg.Wait()
}

var (
	ch = make(chan string, 10)
)

func incomingURLs() <-chan string {
	return ch
}
