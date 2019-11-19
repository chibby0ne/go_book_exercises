package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch, body := make(chan string), make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch, body)
	}
	for range os.Args[1:] {
		fmt.Fprintln(os.Stderr, <-ch)
		fmt.Fprintln(os.Stdout, <-body)
	}
	fmt.Fprintf(os.Stderr, "%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string, body chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, len(b), url)
	body <- fmt.Sprintf("%s", b)
}
