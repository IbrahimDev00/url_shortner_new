package helpers

import (
	"os"
	"strings"
)

func EnsureHTTPPrefix(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
}

func IsDifferentDomain(url string) bool {
	domain := os.Getenv("DOMAIN")

	// Check if the URL exactly matches the domain
	if url == domain {
		return false
	}

	// Remove protocol and www prefix, then extract the domain part of the URL
	cleanURL := strings.TrimPrefix(url, "http://")
	cleanURL = strings.TrimPrefix(cleanURL, "https://")
	cleanURL = strings.TrimPrefix(cleanURL, "www.")
	cleanURL = strings.Split(cleanURL, "/")[0]

	// Compare the cleaned URL with the domain
	return cleanURL != domain
}
