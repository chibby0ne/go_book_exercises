package main

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outline: %v\n", err)
		os.Exit(1)
	}
	outline(nil, doc)
}

func outline(stack *[]string, n *html.Node) {
	if n.Type == html.ElementNode {
		if stack != nil {
			*stack = append(*stack, n.Data)
		} else {
			s := []string{}
			s = append(s, n.Data)
			stack = &s
		}
		fmt.Println(*stack)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c)
	}
}
