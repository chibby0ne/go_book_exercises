// Use the breadthFirst function to explore a different structure. For example you could use the
// course dependencies from the topoSort example (a directed graph), the file system hierarchy on
// you computer (a tree) or a list of bus or subway routes downloaded from your city government's
// web site (an undirected graph).

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crawl(p string) []string {
	fmt.Println(p)
	fileInfo, err := os.Stat(p)
	if err != nil {
		log.Fatal(err)
	}
	if !fileInfo.IsDir() {
		return nil
	}
	files, err := ioutil.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}
	var fileNames []string
	for _, f := range files {
		fileNames = append(fileNames, path.Join(p, f.Name()))
	}
	return fileNames
}

func getAllAbsPaths(paths []string) []string {
	var absPaths []string
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			log.Fatal(err)
		}
		absPaths = append(absPaths, absPath)
	}
	return absPaths
}

func main() {
	breadthFirst(crawl, getAllAbsPaths(os.Args[1:]))
}
