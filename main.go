package main

import (
	"encoding/csv"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"os"
	"strings"
)

func nodeValue(node *html.Node) string {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			return strings.TrimSpace(c.Data)
		}
	}
	return ""
}

func tdValues(tr_node *html.Node) []string {
	var result []string
	for c := tr_node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "td" {
			result = append(result, nodeValue(c))
		}
	}
	return result
}

func findTable(r io.Reader) {
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}
	dataTable := findById(doc, "ConsumptionData")
	if dataTable == nil {
		// Return an error
		return
	}
	// Now dig out the rows of data

	// Find the tbody child,
	tbodies := childrenOfType(dataTable, "tbody")
	fmt.Printf("Found %d tbody children\n", len(tbodies))
	if len(tbodies) == 0 {
		fmt.Println("no tbody found")
		return
	}
	tbody := tbodies[0]
	//then for each tr in tbody, dig out the
	// text inside each td
	rows := childrenOfType(tbody, "tr")
	fmt.Printf("Found %d tr children\n", len(rows))
	if len(rows) == 0 {
		fmt.Println("no rows found")
		return
	}

	csvWriter := csv.NewWriter(os.Stdout)
	for _, tr := range rows {
		dataRow := tdValues(tr)
		// fmt.Println("Row: ", dataRow)
		if err := csvWriter.Write(dataRow); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	// Write any buffered data to the underlying writer (standard output).
	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// fmt.Printf("hello, world\n")
	// s := `<p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul>`
	// htmlReader := strings.NewReader(s)
	htmlReader, err := os.Open("bwsc-monthly.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlReader.Close()
	findTable(htmlReader)
}
