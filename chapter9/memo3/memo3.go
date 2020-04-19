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

// A Memo caches the results of calling a Func.
type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]result
}

// Func is the type of the function to memoize
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// Get is concurrency safe and performance improves compared to memo2 where the
// whole function is the critical section. Unfortunately now sometimes there
// are some URLS fetched twice. We have redundant work. We need duplicate
// suppression.
func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	memo.mu.Unlock()
	if !ok {
		res.value, res.err = memo.f(key)
		// Between the two cirtical sections, several goroutines may race to
		// compute f(key) and update the map
		memo.mu.Lock()
		memo.cache[key] = res
		memo.mu.Unlock()
	}
	return res.value, res.err
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
