// Using the token-based decoder API, write a promgram that will read an
// arbitrary XML document and construct a tree of generic nodes that represents
// it. Nodes are of two kinds: CharData nodes represent text strings, and
// Element nodes represent named elements and their attributes. Each element
// node has a slice of child nodes.
//
// You many find the following delcarations helpful:
//
// import "encoding/xml"
//
// type Node interface{} // CharData or *Element
//
// type CharData string
//
// type Element struct {
// 	Type     xml.Name
// 	Attr     []xml.Attr
// 	Children []Node
// }
package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Node interface{}

type CharData string

type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

func (e *Element) String() string {
	var builder strings.Builder
	traverse(e, &builder, 0)
	return builder.String()
}

// traverse is using the principle: be flexible in what you receive
func traverse(n Node, writer io.Writer, depth int) {
	switch n := n.(type) {

	case *Element:
		fmt.Fprintf(writer, "%*s%s %s\n", depth*2, "", n.Type, n.Attr)
		for _, child := range n.Children {
			traverse(child, writer, depth+1)
		}

	case CharData:
		// CharData is detected even for elements that don't hold text between
		// their tags, therefore we remove it otherwise it will keep adding
		// blank lines since the "\n" would still be printed
		if strings.TrimSpace(string(n)) == "" {
			return
		}
		fmt.Fprintf(writer, "%*s%s\n", depth*2, "", string(n))

	default:
		panic(fmt.Sprintf("unexpected type: %T", n))
	}
}

// parseXML is using the principle: be flexible (interfaces) in what you receive and strict in what you return (actual dynamic types)
func parseXML(reader io.Reader) (*Element, error) {
	dec := xml.NewDecoder(reader)

	var root *Element
	var currentElem *Element
	var prevElem *Element

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		switch tok := tok.(type) {

		case xml.StartElement:
			elem := &Element{
				Type: tok.Name,
				Attr: tok.Attr,
			}
			if root == nil {
				root = elem
			} else {
				if currentElem.Children == nil {
					currentElem.Children = []Node{elem}
				} else {
					currentElem.Children = append(currentElem.Children, elem)
				}
				prevElem = currentElem
			}
			currentElem = elem

		case xml.EndElement:
			currentElem = prevElem

		case xml.CharData:
			currentElem.Children = append(currentElem.Children, CharData(tok))
		}
	}
	return root, nil
}

func main() {
	root, err := parseXML(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(root)
}
