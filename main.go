package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
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

func getDailyUsage(acct, access string) {
	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
        log.Fatal(err)
	}

	client := &http.Client{
        Jar: jar,
	}

	url := "https://old.bwsc.org/ACCOUNTS/security_main.asp?AcctNum=%s&MtrNum=%s"
	url = fmt.Sprintf(url, acct, access)
	if _, err = client.Get(url); err != nil {
        log.Fatal(err)
	}

	var resp *http.Response
	daily_url := "https://old.bwsc.org/ACCOUNTS/readings_daily_30.asp"
	if resp, err = client.Get(daily_url); err != nil {
        log.Fatal(err)
	}
	// daily_html, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	findTable(resp.Body)
	resp.Body.Close()
}


func main() {
	debugPtr := flag.Bool("debug", false, "Enable debugging")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [flags] ACCOUNT_NUMBER ACCESS_NUMBER\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	debug := *debugPtr
	if debug {
		fmt.Println("debugging enabled")
	} else {
		fmt.Println("debugging disabled")
	}

	fmt.Println("args:", os.Args)
	fmt.Println("flag.Args:", flag.Args())
	if flag.NArg() < 2 {
		fmt.Println("not enough arguments")
		flag.Usage()
		os.Exit(2)
	}
	accountNum := flag.Arg(0)
	accessNum := flag.Arg(1)
	fmt.Println("accountNum:", accountNum)
	fmt.Println("accessNum:", accessNum)


	// Now retrieve data from bwsc.org
	getDailyUsage(accountNum, accessNum)


	return

	htmlReader, err := os.Open("bwsc-monthly.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlReader.Close()
	findTable(htmlReader)
}
