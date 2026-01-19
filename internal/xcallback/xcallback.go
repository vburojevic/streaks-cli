package xcallback

import "strings"

// Supported reports whether any discovered URL scheme implies x-callback support.
// Currently Streaks does not advertise x-callback URLs, so this stays false.
func Supported(urlSchemes []string) bool {
	for _, scheme := range urlSchemes {
		if strings.Contains(strings.ToLower(scheme), "x-callback") {
			return true
		}
	}
	return false
}
