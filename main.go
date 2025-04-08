package main

import (
	"flag"
	"fmt"
	"github.com/auperman-lab/web-lab5/src"
)

func main() {
	u := flag.String("u", "", "Insert a http address")
	h := flag.Bool("h", false, "Show this help")
	s := flag.String("s", "", "Search Something")

	flag.Parse()

	if *h {
		fmt.Printf("Available commands:\n")
		fmt.Printf("\t -u <URL>:  make an HTTP request to the specified URL \n")
		fmt.Printf("\t -s <search-term>:  make an HTTP request to search and print top 10 results\n")
		fmt.Printf("\t -h: Show this help\n")
	}
	if *u != "" {
		body, err := src.Fetch(*u)
		if err != nil {
			fmt.Printf("Invalid URL page \n Error: %s\n", err.Error())
		}
		//fmt.Println(body)
		resp, err := src.ParseHtml(body)
		if err != nil {
			fmt.Printf("Unable to parse URL page \n Error: %s\n", err.Error())
		}
		fmt.Println(resp)
	}
	if *s != "" {
		result, err := src.ScrapeDuckDuckGo(*s)
		if err != nil {
			fmt.Println("Error fetching results:", err)
			return
		}
		for _, link := range result {
			fmt.Println(link)
		}

	}
}
