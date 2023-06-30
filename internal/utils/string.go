package utils

import "strings"

func NormalizeString(in string) string {
	return strings.ToLower(strings.TrimSpace(in))
}
