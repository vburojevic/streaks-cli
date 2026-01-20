package discovery

import (
	"os"
	"strings"
)

func PreferredLocales() []string {
	envs := []string{"LC_ALL", "LC_MESSAGES", "LANG"}
	locales := make([]string, 0, len(envs))
	for _, env := range envs {
		if value := os.Getenv(env); value != "" {
			locale := strings.Split(value, ".")[0]
			locale = strings.ReplaceAll(locale, "_", "-")
			if locale != "" {
				locales = append(locales, locale)
			}
		}
	}
	return locales
}
