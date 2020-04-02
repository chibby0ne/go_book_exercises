package main

import (
	"github.com/adonovan/gopl.io/ch8/thumbnail"
)

// Incorrect. The main goroutine doesn't wait for all the goroutines to finish,
// and therefore no work is done
func makeThumbnails2(filenames []string) {
	for _, f := range filenames {
		go thumbnail.ImageFile(f)
	}
}

func main() {

}
