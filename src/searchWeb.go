package src

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func ScrapeDuckDuckGo(query string) ([]string, error) {
	searchURL := "https://html.duckduckgo.com/html/"
	limit := 10

	data := url.Values{}
	data.Set("q", query)

	resp, err := http.PostForm(searchURL, data) // Use POST to simulate a real search
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search results: %v", err)
	}
	defer resp.Body.Close()

	links, err := extractLinks(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to extract links: %v", err)
	}

	if len(links) > limit {
		links = links[:limit]
	}

	return links, nil
}

func extractLinks(body io.Reader) ([]string, error) {
	tokenizer := html.NewTokenizer(body)
	var links []string

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return links, nil
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" && strings.HasPrefix(attr.Val, "http") {
						if !contains(links, attr.Val) {
							links = append(links, attr.Val)

						}
					}
				}
			}
		}
	}
}

func contains(links []string, link string) bool {
	for _, l := range links {
		if l == link {
			return true
		}
	}
	return false
}
