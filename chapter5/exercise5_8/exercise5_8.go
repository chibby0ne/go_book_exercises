// Modify forEachNode so that the pre and post functions return a boolean result
// indicating whether to continue the traversal. Use it to write a function ElementByID with the
// folowing signature that finds the first HTML element with the specified id attribute. The
// function should stop the traversal as soon as a a match is found.
// func ElementByID(doc *html.Node, id string) *html.Node
package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
)

func forEachNode(n *html.Node, s string, pre, post func(n *html.Node, s string) bool) (*html.Node, bool) {
	if n == nil {
		return nil, false
	}
	if pre != nil {
		if pre(n, s) {
			return n, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		n, found := forEachNode(c, s, pre, post)
		if found {
			return n, found
		}
	}
	return nil, false
}

func startElement(n *html.Node, s string) bool {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == s {
				return true
			}
		}
	}
	return false
}

func ElementByID(doc *html.Node, id string) *html.Node {
	node, _ := forEachNode(doc, id, startElement, nil)
	return node
}

func main() {
	node, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	nodeWithId := ElementByID(node, "content")
	fmt.Printf("node is: %+v\n", nodeWithId)
}
