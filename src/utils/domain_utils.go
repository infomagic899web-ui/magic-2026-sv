package utils

import (
	"net/url"
	"strings"
)

// DomainOrigin normalizes a URL or hostname into a consistent origin form (scheme + domain).
// e.g. "https://sub.example.com" → "https://example.com"
func DomainOrigin(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		// Handle input without scheme (like just "localhost")
		parsed = &url.URL{Scheme: "http", Host: rawURL}
	}

	hostParts := strings.Split(parsed.Hostname(), ".")
	length := len(hostParts)
	domain := parsed.Hostname()

	// Convert subdomain -> main domain
	if length > 2 {
		domain = strings.Join(hostParts[length-2:], ".")
	}

	return parsed.Scheme + "://" + domain, nil
}

// DomainName extracts the main domain portion from a host string.
// e.g. "sub.api.example.com" → "example.com"
func DomainName(host string) string {
	hostParts := strings.Split(host, ".")
	length := len(hostParts)
	if length > 2 {
		return strings.Join(hostParts[length-2:], ".")
	}
	return host
}

// IsSameOrigin safely compares two URLs/origins to verify if they share the same scheme + main domain.
// Handles cases like ports, subdomains, and missing schemes.
func IsSameOrigin(a, b string) bool {
	if a == "" || b == "" {
		return false
	}

	aParsed, errA := url.Parse(a)
	bParsed, errB := url.Parse(b)
	if errA != nil || errB != nil {
		return false
	}

	// Normalize missing schemes
	if aParsed.Scheme == "" {
		aParsed.Scheme = "http"
	}
	if bParsed.Scheme == "" {
		bParsed.Scheme = "http"
	}

	// Extract main domains
	aDomain := DomainName(aParsed.Hostname())
	bDomain := DomainName(bParsed.Hostname())

	// Match both scheme and normalized domain
	return strings.EqualFold(aParsed.Scheme, bParsed.Scheme) &&
		strings.EqualFold(aDomain, bDomain)
}
