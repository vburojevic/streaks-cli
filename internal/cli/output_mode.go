package cli

import (
	"fmt"
	"strings"
)

type outputMode string

const (
	outputHuman outputMode = "human"
	outputJSON  outputMode = "json"
	outputPlain outputMode = "plain"
)

func parseOutputMode(value string) (outputMode, error) {
	if value == "" {
		return outputHuman, nil
	}
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case string(outputHuman):
		return outputHuman, nil
	case string(outputJSON):
		return outputJSON, nil
	case string(outputPlain):
		return outputPlain, nil
	default:
		return outputHuman, fmt.Errorf("invalid output mode: %s", value)
	}
}
