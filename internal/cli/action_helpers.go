package cli

import (
	"fmt"

	"streaks-cli/internal/discovery"
)

func findActionDef(id string) (discovery.ActionDef, error) {
	for _, def := range availableActionDefs() {
		if def.ID == id {
			return def, nil
		}
	}
	return discovery.ActionDef{}, fmt.Errorf("unknown action: %s", id)
}

func samplePayload(def discovery.ActionDef) map[string]any {
	payload := map[string]any{}
	if def.RequiresTask {
		payload["task"] = "<task>"
	}
	if len(def.ParamOptions) > 0 {
		for key, values := range def.ParamOptions {
			if len(values) > 0 {
				payload[key] = values[0]
			} else {
				payload[key] = ""
			}
		}
	}
	return payload
}
