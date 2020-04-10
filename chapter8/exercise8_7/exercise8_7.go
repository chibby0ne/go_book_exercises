// Write a concurrent program that creates a local mirror of a web site,
// fetching each reachable page and writing it to a directory on the local
// disk. Only pages withing the original domain should be fetched. URls withing
// mirrored pages should be altered as needed so that they refer to the
// mirrored page, not the original.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var (
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func Extract(uri string) ([]string, error) {
	uriParsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("error getting %s: %s", uri, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", uri, err)
	}
	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "link") {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue
				}
				// fmt.Println("found link: ", link.String(), " link.Host:", link.Host, "link.Path:", link.Path)
				if link.Host == uriParsed.Host {
					links = append(links, link.String())
				}
			}
		} else if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key != "src" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue
				}
				// fmt.Println("found link: ", link.String(), " link.Host:", link.Host, "link.Path:", link.Path)
				if link.Host == uriParsed.Host {
					links = append(links, link.String())
				}
			}
		} else if n.Type == html.ElementNode && n.Data == "use" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(strings.Split(a.Val, "#")[0])
				if err != nil {
					continue
				}
				// fmt.Println("found link: ", link.String(), " link.Host:", link.Host, "link.Path:", link.Path)
				if link.Host == uriParsed.Host {
					links = append(links, link.String())
				}
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func fetchAllDomainPages(uri string) []string {
	var validlinks, foundlinks []string
	foundlinks = append(foundlinks, uri)
	visited := make(map[string]bool)
	for i := 0; i < len(foundlinks); i++ {
		link := foundlinks[i]
		if visited[link] {
			continue
		}
		visited[link] = true
		links, err := Extract(link)
		if err != nil {
			log.Fatal(err)
		}
		foundlinks = append(foundlinks, links...)
	}
	for k, _ := range visited {
		validlinks = append(validlinks, k)
	}
	return validlinks
}

func createFile(rootdir string, link *url.URL) (*os.File, error) {
	var path string
	// If path is the root path then use dir as the directory to create, otherwise prependir dir to the path's dir
	if link.Path != "/" {
		path = filepath.Join(rootdir, filepath.Dir(link.Path))
	} else {
		path = rootdir
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	// create directory
	err = os.MkdirAll(absPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot create directory: %v", err)
	}
	base := filepath.Base(link.Path)
	// if the file ends with a / then this is the index of a section, or it
	// could be the home page with a /, in that case just open the file (if it
	// has already been dowloaded) i.e: www.google.com == www.google.com/ but
	// both strings could be in the links list, we need to open it the second
	// time, as creating it would throw an error
	if base == "." || base == "/" && link.Path == "/" {
		base = "index.html"
	}
	// create file by concatenating dir path and base path
	absPathFile := filepath.Join(absPath, base)

	// if the file already exists we return nil, to skip file creation and page download
	if _, err := os.Stat(absPathFile); os.IsExist(err) {
		log.Printf("file: %v already exists skipping download", absPathFile)
		return nil, nil
	}
	file, err := os.Create(absPathFile)
	if err != nil {
		return nil, fmt.Errorf("cannot create file: %v", err)
	}
	return file, nil
}

func downloadPage(file *os.File, urlString string) error {
	resp, err := client.Get(urlString)
	if err != nil {
		log.Print(err)
		return err
	}
	defer resp.Body.Close()
	if err != nil {
		log.Print(err)
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func stringToURL(s string) (*url.URL, error) {
	parsedURL, err := url.Parse(s)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return parsedURL, nil
}

func main() {
	domainRoot := os.Args[1]
	links := fetchAllDomainPages(domainRoot)
	rootURL, err := stringToURL(domainRoot)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	tokens := make(chan struct{}, 8)

	for _, link := range links {
		// fmt.Println(link)
		tokens <- struct{}{}
		wg.Add(1)
		go downloadPageAndCreateFile(&wg, rootURL, link)
		<-tokens
	}
	wg.Wait()
}

func downloadPageAndCreateFile(wg *sync.WaitGroup, rootURL *url.URL, link string) {
	defer wg.Done()
	linkURL, err := stringToURL(link)
	if err != nil {
		log.Fatal("error creating link:", err)
	}
	file, err := createFile(rootURL.Host, linkURL)
	if err != nil {
		log.Fatal("error creating file: ", err)
	}
	if file == nil {
		return
	}
	downloadPage(file, linkURL.String())
}
