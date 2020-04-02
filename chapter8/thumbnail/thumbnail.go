package main

import (
	"github.com/adonovan/gopl.io/ch8/thumbnail"
	"log"
)

func makeThumbnails(filenames []string) {
	for _, f := range filenames {
		if _, err := thumbnail.ImageFile(f); err != nil {
			log.Println(err)
		}
	}
}

func main() {

}
