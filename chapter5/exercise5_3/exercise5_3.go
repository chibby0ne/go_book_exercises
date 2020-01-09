// Write a function to print the contents of all text nodes in an HTML document
// tree. Do not descend into <script> or <style> elements, since their contents
// are not visible in a web server.
package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
	"strings"
)

func printContentTextNodes(n *html.Node, contents []string) []string {
	if n == nil || n.DataAtom.String() == "script" || n.DataAtom.String() == "style" {
		return contents
	}
	if n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" {
		contents = append(contents, n.Data)
	}
	contents = printContentTextNodes(n.FirstChild, contents)
	contents = printContentTextNodes(n.NextSibling, contents)
	return contents
}

func main() {
	node, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var contents []string
	printContentTextNodes(node, contents)
	for _, v := range printContentTextNodes(node, contents) {
		fmt.Println(v)
	}
}
