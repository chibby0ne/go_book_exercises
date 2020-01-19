// The startElement and endElement functions in gopl.io/ch5/outline2 share a global
// variable depth. Turn them into anonymous functions that share a variable local
// to the outline function

package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
)

func forEachNode(n *html.Node) {
	forEachNodeRec(0, n)
}

func forEachNodeRec(depth int, n *html.Node) {
	func(n *html.Node) {
		if n.Type == html.ElementNode {
			fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
			depth++
		}
	}(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNodeRec(depth, c)
	}
	func(n *html.Node) {
		if n.Type == html.ElementNode {
			depth--
			fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
		}
	}(n)
}

func main() {
	node, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	forEachNode(node)
}
