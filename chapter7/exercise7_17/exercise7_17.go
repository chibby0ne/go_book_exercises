// Extend xmlselect so that elements may be selected not just by name, but by
// their attributes too, in the manner of CSS, so that, for instance, an
// element like <div id="page" class="wide"> could be selected by a matching id
// or class as well as its name.

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	elements *string = flag.String("elements", "", "elements to look for")
	classes  *string = flag.String("classes", "", "classes to look for")
	ids      *string = flag.String("ids", "", "attributes to look for")
	verbose  *bool   = flag.Bool("v", false, "verbose output")
)

func getListArguments(s *string) []string {
	return strings.Split(*s, ",")
}

func logfVerbose(fmt string, v ...interface{}) {
	if *verbose {
		log.Printf(fmt, v...)
	}
}

func logVerbose(v ...interface{}) {
	if *verbose {
		log.Print(v...)
	}
}

func main() {
	flag.Parse()

	elemsList := getListArguments(elements)
	classesList := getListArguments(classes)
	idsList := getListArguments(ids)

	dec := xml.NewDecoder(os.Stdin)
	var elementsFoundStack []string
	classesMap := make(map[string][]string)
	idsMap := make(map[string][]string)
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			elementsFoundStack = append(elementsFoundStack, tok.Name.Local)

			classesStack := containsAttribute(tok.Attr, classesList, "class")
			logfVerbose("classesStack: %+v\n", classesStack)
			if len(classesStack) > 0 {
				classesMap[createKeyFromElementStack(elementsFoundStack)] = classesStack
			}

			idsStack := containsAttribute(tok.Attr, idsList, "id")
			logfVerbose("idStack: %+v\n", idsStack)
			if len(idsStack) > 0 {
				idsMap[createKeyFromElementStack(elementsFoundStack)] = idsStack
			}

		case xml.EndElement:
			// Need to remove the key/pair from the classes map and Ids map
			// before we pop the stack of the elementsFoundStack
			// since we use the elementsFoundStack to create the key for the
			// aforementioned maps, and in this step we pop the
			// elementsFoundStack
			delete(classesMap, createKeyFromElementStack(elementsFoundStack))
			delete(idsMap, createKeyFromElementStack(elementsFoundStack))
			elementsFoundStack = elementsFoundStack[:len(elementsFoundStack)-1]

		case xml.CharData:
			if containsAll(elementsFoundStack, elemsList) {
				fmt.Printf("%s: %s\n", strings.Join(elementsFoundStack, " "), tok)
			}
			if classes, ok := classesMap[createKeyFromElementStack(elementsFoundStack)]; ok {
				fmt.Printf("%s: %s\n", strings.Join(classes, " "), tok)
			}
			if ids, ok := idsMap[createKeyFromElementStack(elementsFoundStack)]; ok {
				fmt.Printf("%s: %s\n", strings.Join(ids, " "), tok)
			}
		}
	}
}

func createKeyFromElementStack(elements []string) string {
	return strings.Join(elements, ",")
}

func containsAttribute(attrs []xml.Attr, attrsToSearch []string, attrName string) []string {
	var attributes []string
	if foundAttrs := getAttributes(attrs, attrName); len(foundAttrs) > 0 {
		attributes = getMatchingValues(foundAttrs, attrsToSearch)
	}
	return attributes
}

func getAttributes(attrs []xml.Attr, attributeName string) []xml.Attr {
	var foundAttrs []xml.Attr
	for _, attr := range attrs {
		logfVerbose("attr.Name.Local: %v, attributeName: %v", attr.Name.Local, attributeName)
		if attr.Name.Local == attributeName {
			foundAttrs = append(foundAttrs, attr)
		}
	}
	return foundAttrs
}

func getMatchingValues(attrs []xml.Attr, attrsToSearch []string) []string {
	var matchingAttrValues []string
	for _, attributeToSearch := range attrsToSearch {
		for _, attr := range attrs {
			logfVerbose("atttributeToSearch: %v, attr.Name.Local: %v", attributeToSearch, attr.Value)
			if attributeToSearch == attr.Value {
				matchingAttrValues = append(matchingAttrValues, attr.Value)
			}
		}
	}
	return matchingAttrValues
}

func containsAll(x, y []string) bool {
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if x[0] == y[0] {
			y = y[1:]
		}
		x = x[1:]
	}
	return false
}
