package src

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
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

	if !(checkCache(host, path)) {
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

		response, err := io.ReadAll(conn)
		if err != nil {
			return "", fmt.Errorf("error reading response: %v", err)
		}

		redirectURL := checkRedirect(string(response))
		if redirectURL != "" {
			fmt.Println("Redirecting to:", redirectURL)
			return Fetch(redirectURL) // Recursively fetch the new location
		}

		err = addToCache(host, path, string(response))
		if err != nil {
			return "", err
		}

		return string(response), nil
	}
	fmt.Println("\nThis website is cached!\n")
	response, err := getFromCache(host, path)
	if err != nil {
		return "", err
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

func checkRedirect(response string) string {
	scanner := bufio.NewScanner(strings.NewReader(response))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "HTTP/1.1 3") || strings.HasPrefix(line, "HTTP/1.0 3") {
			fmt.Println("Redirect detected:", line)
		}

		if strings.HasPrefix(line, "Location: ") {
			redirectURL := strings.TrimSpace(strings.TrimPrefix(line, "Location: "))
			return redirectURL
		}

		if line == "" {
			break
		}
	}
	return ""
}
