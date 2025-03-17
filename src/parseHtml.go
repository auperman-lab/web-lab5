package src

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

func ParseHtml(body string) (string, error) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return "", err
	}

	tagsToRemove := []string{"script", "style", "meta", "link"}
	cleanHtml := ""

	cleanHtml += removeTags(doc.FirstChild, tagsToRemove)

	fmt.Println(cleanHtml)

	return cleanHtml, nil

}

func removeTags(n *html.Node, tags []string) string {
	var sb strings.Builder

	if n.Type == html.ElementNode {
		for _, tag := range tags {
			if n.Data == tag {
				return ""
			}
		}

		sb.WriteString("<" + n.Data)

		for _, attr := range n.Attr {
			sb.WriteString(fmt.Sprintf(" %s=%q", attr.Key, attr.Val))
		}
		sb.WriteString(">")

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			sb.WriteString(removeTags(c, tags))
		}

		sb.WriteString(fmt.Sprintf("</%s>", n.Data))
	}

	if n.Type == html.TextNode {
		sb.WriteString(n.Data)
	}

	if n.Type == html.CommentNode {
		sb.WriteString("<!--")
		sb.WriteString(n.Data)
		sb.WriteString("-->")
	}

	if n.Type == html.DoctypeNode {

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			sb.WriteString(removeTags(c, tags))
		}
	}

	return sb.String()
}
