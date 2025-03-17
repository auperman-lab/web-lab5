package src

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/url"
)

func Fetch(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	host := parsedURL.Host
	path := parsedURL.Path
	if path == "" {
		path = "/"
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", host), &tls.Config{})
	if err != nil {
		return "", err
	}
	defer conn.Close()

	request := fmt.Sprintf(
		"GET %s HTTP/1.1\r\n"+
			"Host: %s\r\n"+
			"Connection: close\r\n"+
			"User-Agent: : Mozilla/5.0\r\n\r\n", path, host)

	_, err = conn.Write([]byte(request))
	if err != nil {
		return "", err
	}

	var response string
	scanner := bufio.NewScanner(conn)
	headersDone := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			headersDone = true
			continue
		}
		if headersDone {
			response += line + "\n"
		}
	}

	return response, nil

}
