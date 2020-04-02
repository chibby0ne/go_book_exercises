package main

import (
	"github.com/adonovan/gopl.io/ch8/thumbnail"
)

func makeThumbnails(filenames []string) {
	ch := make(chan struct{})
	for _, f := range filenames {
		go func(f string) {
			thumbnail.ImageFile(f)
			ch <- struct{}{}
		}(f)
	}
	// Wait for gorutines to complete
	for range filenames {
		<-ch
	}
}

func main() {

}
