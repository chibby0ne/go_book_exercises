// Modify crawl to make local copies of the pages it finds, creating
// directories as necessary. Don't make copies of pages that come from a
// different domain. For example, if the original page comes from golang.org,
// save all files from there, but exclude ones from vimeo.com

package main

import (
	"fmt"
	"github.com/chibby0ne/go_book_exercises/chapter5/links"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

var REGEX = regexp.MustCompile(`(?:\/|\.)((?:\w+\.)(?:org|com|net|io))/?(.*)`)

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

func downloadPage(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("downloadPage: could not create file: %v", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = file.Write(b)
	if err != nil {
		return fmt.Errorf("downloadPage: could not write to file: %v", err)
	}
	file.Close()
	return nil
}

func createDirectories(url string) error {
	for i, c := range url {
		if c == '/' {
			if _, err := os.Stat(url[:i]); os.IsNotExist(err) {
				if err := os.Mkdir(url[:i], 0755); err != nil {
					return fmt.Errorf("createDirectories: %v", err)
				}
			}
		}
	}
	if _, err := os.Stat(url); os.IsNotExist(err) {
		if err := os.Mkdir(url, 0755); err != nil {
			return fmt.Errorf("createDirectories: %v", err)
		}
	}
	return nil
}

func getRootUrlAndPath(url string) (string, string) {
	for _, val := range REGEX.FindAllStringSubmatch(url, -1) {
		return val[1], val[2]
	}
	return "", ""
}

func crawl(url string) []string {
	fmt.Println(url)
	rootUrl, urlPath := getRootUrlAndPath(url)
	if rootUrl == "" && urlPath == "" {
		log.Fatalf("not a valid url: %v", url)
	}
	if _, err := os.Stat(rootUrl); os.IsNotExist(err) {
		if err := os.Mkdir(rootUrl, 0755); err != nil {
			log.Fatalf("could no create dir: %v", err)
		}
	}
	urlDirs := urlPath
	filename := path.Join(rootUrl, urlPath)
	if strings.HasSuffix(urlPath, ".html") {
		index := strings.LastIndex(urlPath, "/")
		urlDirs = urlPath[:index]
	} else {
		filename += "/index.html"
	}
	if err := createDirectories(path.Join(rootUrl, urlDirs)); err != nil {
		log.Fatalf("could no create nested dirs: %v", err)
	}
	fmt.Println("filename: ", filename)
	if err := downloadPage(url, filename); err != nil {
		log.Fatalf("could not download page: %v", err)
	}
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	breadthFirst(crawl, os.Args[1:])
}
