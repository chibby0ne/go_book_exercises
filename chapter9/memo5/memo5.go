// Package memo provides a concurrency-unsafe memorization of a function of type Func.
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
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

// A request is a message requesting that the Func be applied to Key
type request struct {
	key      string
	response chan<- result // the client wants a single result
}

// A Memo caches the results of calling a Func.
type Memo struct {
	requests chan request
}

func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

// Get is concurrency safe
func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() {
	close(memo.requests)
}

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			// This is the first reuqest for this key.
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string) {
	// Evaluate the function
	e.res.value, e.res.err = f(key)
	// Broadcast the ready condition
	close(e.ready)
}

func (e *entry) deliver(res chan<- result) {
	// Wait for ready condition
	<-e.ready
	// Send the result to the client
	res <- e.res
}

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)

}

func Concurrent() {
	m := New(httpGetBody)
	var wg sync.WaitGroup
	for url := range incomingURLs() {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			start := time.Now()
			value, err := m.Get(url)
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
