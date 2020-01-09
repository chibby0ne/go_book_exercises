// Implement countWordsAndImages (see exercise 4.9 for word-splitting)
package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

func countWordsAndImages(n *html.Node) (words, images int) {
	var nodes []*html.Node
	nodes = append(nodes, n)
	for len(nodes) > 0 {
		current := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]

		switch current.Type {

		case html.TextNode:
			reader := bufio.NewScanner(strings.NewReader(current.Data))
			reader.Split(bufio.ScanWords)
			for reader.Scan() {
				words++
			}
		case html.ElementNode:
			if current.Data == "img" {
				images++
			}
		}
		for node := current.FirstChild; node != nil; node = node.NextSibling {
			nodes = append(nodes, node)
		}
	}
	return words, images
}

func CountWordsAndImages(url string) (words, images int, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(res.Body)
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}

func main() {
	for _, url := range os.Args[1:] {
		words, images, err := CountWordsAndImages(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		}
		fmt.Println("words:", words, "images:", images)
	}
}
