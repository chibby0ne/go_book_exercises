package main

import (
	"log"
	"os"
	"sync"

	"github.com/adonovan/gopl.io/ch8/thumbnail"
)

// This function requires the passing of a closed channel otherwise it will result in runtime time error: "all goroutines are asleep - deadlock!"
func makeThumbnails6(filenames <-chan string) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup

	for f := range filenames {
		wg.Add(1)

		// Worker
		go func(f string) {
			defer wg.Done()
			thumb, err := thumbnail.ImageFile(f)
			if err != nil {
				log.Println(err)
				return
			}
			info, _ := os.Stat(thumb)
			sizes <- info.Size()
		}(f)
	}

	// Closer
	go func() {
		wg.Wait()
		close(sizes)
	}()

	var total int64
	for size := range sizes {
		total += size
	}
	return total
}

func main() {

}
