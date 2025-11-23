package utils

import "strings"

func Slugify(str string) string {
	// Lowercase
	s := strings.ToLower(str)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "")

	// Remove invalid characters
	allowed := func(r rune) rune {
		if (r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') ||
			r == '-' {
			return r
		}
		return -1
	}

	s = strings.Map(allowed, s)
	return s
}
