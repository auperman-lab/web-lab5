package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

func main() {
	url := "https://www.scrapingcourse.com/ecommerce/"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}
	traverse(doc)
}

func traverse(n *html.Node) {
	if n.Type == html.ElementNode {
		fmt.Println("Tag:", n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(c)
	}
}
