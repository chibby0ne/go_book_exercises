// Extend the `visit` function so that it extacts other kinds of links from the
// document such as images, scripts and style sheets.

package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
)

func appendAttributeValue(attrs []html.Attribute, links []string, attr string) []string {
	for _, a := range attrs {
		if a.Key == attr {
			links = append(links, a.Val)
		}
	}
	return links
}

func visit(links []string, n *html.Node) []string {
	if n == nil {
		return links
	}
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a", "link":
			links = appendAttributeValue(n.Attr, links, "href")
		case "img", "script":
			links = appendAttributeValue(n.Attr, links, "src")
		default:
		}
	}
	links = visit(links, n.FirstChild)
	links = visit(links, n.NextSibling)
	return links
}

func main() {
	node, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var links []string
	for _, link := range visit(links, node) {
		fmt.Println(link)
	}
}
