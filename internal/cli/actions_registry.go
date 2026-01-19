package cli

import (
	"context"
	"os"

	"streaks-cli/internal/discovery"
)

func availableActionDefs() []discovery.ActionDef {
	defs := discovery.DefaultActionDefinitions()
	if os.Getenv(envDisableDiscovery) != "" {
		return defs
	}
	ctx := context.Background()
	disc, err := discovery.Discover(ctx)
	if err != nil || len(disc.Actions) == 0 {
		return defs
	}
	present := make(map[string]discovery.Action, len(disc.Actions))
	for _, action := range disc.Actions {
		present[action.ID] = action
	}
	return filterDefs(defs, present)
}
