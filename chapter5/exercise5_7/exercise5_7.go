// Develop startElement and endElement into a general HTML pretty printer.
// Print comment node,s text nodes, and the attributes of each element (<a href='...'). Use
// short forms like <img/> instead of <img></img> when an element has no children. Write a
// test to ensure that the output can be parsed successfully. (See Chapter 11)
package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
	"strings"
)

var depth int

func isSelfClosingTag(n *html.Node) bool {
	if n.FirstChild != nil {
		return false
	}
	if n.Type != html.ElementNode {
		return false
	}
	switch n.Data {
	case "area", "base", "br", "col", "embed", "hr", "img", "input", "meta", "param", "source", "track", "wbr":
		return true
	}
	return false
}

func createAttrString(attrs []html.Attribute) string {
	var str strings.Builder
	for _, attr := range attrs {
		str.WriteString(fmt.Sprintf("%s=%q ", attr.Key, attr.Val))
	}
	return strings.TrimSpace(str.String())
}

func startElement(n *html.Node) {
	if n.Type == html.ElementNode {
		if len(n.Attr) != 0 {
			fmt.Printf("%*s<%s %s>", depth*2, "", n.Data, createAttrString(n.Attr))
		} else {
			fmt.Printf("%*s<%s>", depth*2, "", n.Data)
		}
		depth++

		if !(n.FirstChild != nil && n.FirstChild.Type == html.TextNode && strings.TrimSpace(n.FirstChild.Data) != "") {
			fmt.Println()
		}
	}
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		depth--
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode && strings.TrimSpace(n.FirstChild.Data) != "" {
			fmt.Printf("</%s>\n", n.Data)
		} else {
			fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
		}
	}
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if isSelfClosingTag(n) {
		fmt.Printf("%*s<%s %s/>\n", depth*2, "", n.Data, createAttrString(n.Attr))
		return
	}
	if pre != nil {
		pre(n)
	}
	if n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" {
		fmt.Printf("%s", strings.TrimSpace(n.Data))
	} else if n.Type == html.CommentNode {
		fmt.Printf("%*s<!--%s-->\n", depth*2, "", n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func main() {
	node, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	forEachNode(node, startElement, endElement)
}
