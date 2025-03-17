package src

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const cacheDir = "./data"

func checkCache(host, path string) bool {
	sanitizedHost := sanitizeFileName(host)
	sanitizedPath := sanitizeFileName(path)
	filePath := fmt.Sprintf("%s/%s_%s.html", cacheDir, sanitizedHost, sanitizedPath)

	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func addToCache(host, path, response string) error {
	sanitizedHost := sanitizeFileName(host)
	sanitizedPath := sanitizeFileName(path)
	filePath := fmt.Sprintf("%s/%s_%s.html", cacheDir, sanitizedHost, sanitizedPath)

	err := writeToFile(filePath, response)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	fmt.Println("Added to cache successfully")

	return nil
}

func getFromCache(host, path string) (string, error) {
	sanitizedHost := sanitizeFileName(host)
	sanitizedPath := sanitizeFileName(path)
	filePath := fmt.Sprintf("%s/%s_%s.html", cacheDir, sanitizedHost, sanitizedPath)

	// Read the cached file content
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading cache file: %v", err)
	}

	return string(data), nil
}

func sanitizeFileName(input string) string {
	// Replace slashes with underscores, remove any characters not suitable for file names
	return strings.NewReplacer("/", "_").Replace(input)
}
