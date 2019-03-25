package main

import (
	"fmt"
	"golang.org/x/net/html"
)

func childrenOfType(node *html.Node, typ string) []*html.Node {
	var result []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == typ {
			result = append(result, c)
		}
	}
	return result
}

func findById(node *html.Node, id string) *html.Node {
	if node.Type == html.ElementNode {
		for _, a := range node.Attr {
			if a.Key == "id" && a.Val == id {
				fmt.Printf("Found %s with id %s\n", node.Data, id)
				return node
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if item := findById(c, id); item != nil {
			return item
		}
	}
	return nil
}
