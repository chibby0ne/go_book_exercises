// Write a version of du that comptues and periodically displays separate totla
// for each of the root directories.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RootSize struct {
	id   int
	size int64
}

var (
	verbose = flag.Bool("v", false, "show verbose progress messages")
	sema    = make(chan struct{}, 20)
)

func walkDir(id int, dir string, n *sync.WaitGroup, fileSizes chan<- RootSize) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(id, subdir, n, fileSizes)
		} else {
			fileSizes <- RootSize{id: id, size: entry.Size()}
		}
	}
}

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

func main() {
	// Determine the initial directories
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse each root of the file tree in parallel
	fileSizes := make(chan RootSize)
	var n sync.WaitGroup
	for i, root := range roots {
		n.Add(1)
		go walkDir(i, root, &n, fileSizes)
	}

	// Closer of the fileSizes channel
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	// Print the results periodically
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}

	nfiles, nbytes := make([]int64, len(roots)), make([]int64, len(roots))

loop:
	for {
		select {
		case rootSize, ok := <-fileSizes:
			if !ok {
				break loop
			}
			nfiles[rootSize.id]++
			nbytes[rootSize.id] += rootSize.size
		case <-tick:
			printDiskUsage(roots, nfiles, nbytes)
		}
	}
	printDiskUsage(roots, nfiles, nbytes)
}

func printDiskUsage(roots []string, nfiles, nbytes []int64) {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s: %d files, %.1f GB", roots[0], nfiles[0], float64(nbytes[0])/1e9))
	for i := 1; i < len(roots); i++ {
		builder.WriteString(fmt.Sprintf("\t%s: %d files, %.1f GB", roots[i], nfiles[i], float64(nbytes[i])/1e9))
	}
	fmt.Println(builder.String())
}
