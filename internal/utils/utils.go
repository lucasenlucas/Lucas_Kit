package utils

import (
	"strings"
)

// NormalizeDomain cleans up user input to extract just the base domain.
func NormalizeDomain(d string) string {
	d = strings.TrimSpace(d)
	d = strings.TrimPrefix(d, "http://")
	d = strings.TrimPrefix(d, "https://")
	d = strings.Split(d, "/")[0]
	return strings.TrimSuffix(d, ".")
}
