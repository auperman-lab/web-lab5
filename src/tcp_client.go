package src

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
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

	var conn net.Conn
	conn, err = fetchHttps(host)

	if err != nil {
		conn, err = fetchHttp(host)
		if err != nil {
			return "", fmt.Errorf("failed to connect to %s (HTTPS/HTTP): %v", urlStr, err)
		}
	}
	defer conn.Close()

	request := fmt.Sprintf(
		"GET %s HTTP/1.1\r\n"+
			"Host: %s\r\n"+
			"Connection: close\r\n"+
			"User-Agent: Mozilla/5.0\r\n\r\n", path, host)

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

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	return response, nil
}

func fetchHttps(host string) (*tls.Conn, error) {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", host), &tls.Config{})
	if err != nil {
		fmt.Printf("failed to connect via HTTPS: %v", err)
		return nil, fmt.Errorf("failed to connect via HTTPS: %v", err)
	}
	fmt.Printf("\n\nConnection via HTTPS succesful\n\n")
	return conn, nil
}

func fetchHttp(host string) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:80", host))
	if err != nil {
		fmt.Printf("failed to connect via HTTP: %v", err)
		return nil, fmt.Errorf("failed to connect via HTTP: %v", err)
	}
	fmt.Printf("\n\nConnection via HTTP succesful\n\n")

	return conn, nil
}
