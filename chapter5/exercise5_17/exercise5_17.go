// Write a variadic function ElementsByTagName that, given an HTML node tree and zero or more names,
// returns all the elements that match one of those names. Here are two example calls:
//
// func Elements ByTagName(doc *html.Node, name ...string) []*html.Node
//
// images := ElementsByTagName(doc, "img")
// headings := ElementsByTagName(doc, "h1", "h2", "h3", "h4")

package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
)

func ElementsByTagName(doc *html.Node, name ...string) []*html.Node {
	var nodes []*html.Node
	visited := make(map[*html.Node]bool)
	stack := []*html.Node{doc}
	for len(stack) != 0 {
		c := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !visited[c] {
			visited[c] = true
			for n := c.FirstChild; n != nil; n = n.NextSibling {
				stack = append(stack, n)
			}
		}
		if c.Type == html.ElementNode {
			for _, tag := range name {
				if c.Data == tag {
					nodes = append(nodes, c)
					break
				}
			}
		}
	}
	return nodes
}

func getNodesStringRepr(nodes []*html.Node) []string {
	var results []string
	for _, node := range nodes {
		if node.FirstChild != nil {
			results = append(results, fmt.Sprintf("%v: %+v\n", node.Data, node.FirstChild.Data))
		} else {
			results = append(results, fmt.Sprintf("%v: %+v\n", node.Data, node.Attr))
		}
	}
	return results
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	images := ElementsByTagName(doc, "img")
	fmt.Printf("%v\n", getNodesStringRepr(images))
	headings := ElementsByTagName(doc, "h1", "h2", "h3", "h4")
	fmt.Printf("%v\n", getNodesStringRepr(headings))
}
