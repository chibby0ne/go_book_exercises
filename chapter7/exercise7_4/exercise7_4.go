// The strings.NewReader function returns a value that satisfies the io.Reader
// interface (and others) by reading from its argument, a string. Implement a
// simple version of NewReader yourself, and use it to make the HTML
// parser(5.2) take input from a string.

package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"os"
)

type StringReader struct {
	s string
}

func NewStringReader(s string) io.Reader {
	return &StringReader{s: s}
}

func (sr *StringReader) Read(p []byte) (n int, err error) {
	n = copy(p, sr.s)
	sr.s = sr.s[n:]
	if len(sr.s) == 0 {
		return n, io.EOF
	}
	return n, nil
}

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading from stdin: %v", err)
		os.Exit(1)
	}
	reader := NewStringReader(string(b))
	doc, err := html.Parse(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
