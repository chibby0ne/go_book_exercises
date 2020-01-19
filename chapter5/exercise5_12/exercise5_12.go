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

func forEachNode(n *html.Node, pre, post func(depth *int, n *html.Node)) {
	forEachNodeRec(0, n, pre, post)
}

func forEachNodeRec(depth int, n *html.Node, pre, post func(depth *int, n *html.Node)) {
	if pre != nil {
		pre(&depth, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNodeRec(depth, c, pre, post)
	}
	if post != nil {
		post(&depth, n)
	}
}

func startElement(depth *int, n *html.Node) {
	f := func(n *html.Node) {
		if n.Type == html.ElementNode {
			fmt.Printf("%*s<%s>\n", *depth*2, "", n.Data)
			*depth++
		}
	}
	f(n)
}

func endElement(depth *int, n *html.Node) {
	f := func(n *html.Node) {
		if n.Type == html.ElementNode {
			*depth--
			fmt.Printf("%*s</%s>\n", *depth*2, "", n.Data)
		}
	}
	f(n)
}

func main() {
	node, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	forEachNode(node, startElement, endElement)
}
